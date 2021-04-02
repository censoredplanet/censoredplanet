//Copyright 2021 Censored Planet

package satellite

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strings"

	//"github.com/censoredplanet/censoredplanet/analysis/pkg/geolocate"
	"github.com/censoredplanet/censoredplanet/analysis/pkg/tarballReader"
	"github.com/censoredplanet/censoredplanet/analysis/pkg/hquack"
	"github.com/cheggaaa/pb/v3"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/trustmaster/goflow"
	set "github.com/deckarep/golang-set"
)


//Big CDNs regex
var cdnRegex = regexp.MustCompile("AMAZON|Akamai|OPENDNS|CLOUDFLARENET|GOOGLE") 

//analyses stores the format for user prompt
type analyses struct {
	Name  string
	Index string
}

//Parse forms the network components for parsing json information
type Parse struct {
	InLine      <-chan string
	ResJSONData chan<- map[string]interface{}
}

//Process parses the json data
func (p *Parse) Process() {
	for line := range p.InLine {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
			log.Warn("Cannot unmarshal JSON data: ", line, ", Error: ", err.Error())
		} else {
			p.ResJSONData <- jsonData
		}
	}
}

//MetaData is the component that assigns measurement metadata to each row
type MetaData struct {
	InMetaData  <-chan map[string]interface{}
	Geolocation <-chan map[string]string
	OutMetaData chan<- map[string]interface{}
}

//Process assigns measurement metadata to each row
func (m *MetaData) Process() {
	geolocation := <-m.Geolocation
	for data := range m.InMetaData {
		vantagePoint := data["resolver"].(string)
		if geolocation, ok := geolocation[vantagePoint]; !ok {
			log.Warn("Did not find geolocation for vantage point: ", vantagePoint)
		} else {
			data["Geolocation"] = geolocation
			m.OutMetaData <- data
		}
	}
}

//Filter is the component for marking untagged measurements
type Filter struct {
	InFilterData  <-chan map[string]interface{}
	ControlAnswers <-chan map[string]*tagsSet
	OutFilterData chan<- map[string]interface{}
}

//Process marks untagged measurements
func (f *Filter) Process() {
	controlAnswers := <-f.ControlAnswers
	for data := range f.InFilterData {
		UntaggedAnswer := false
		query := data["query"].(string)
		numControlTags := 0
		if tags, ok := controlAnswers[query]; ok {
			numControlTags = tags.http.Cardinality() + tags.cert.Cardinality() + tags.asnum.Cardinality() + tags.asname.Cardinality()
		}
		if numControlTags == 0 {
			UntaggedAnswer = true
		}
		
		answersMap := data["answers"].(map[string]interface{})
		flag := true
		for _, answers := range answersMap {
			answerSlice := answers.([]interface{})
			if len(answerSlice) == 0 {
				flag = false
			}
			for _, answerTag := range answerSlice {
				if answerTag != "no_tags" {
					flag = false
				}
			}
		}
		if flag {
			UntaggedAnswer = true
		}
		data["UntaggedAnswer"] = UntaggedAnswer
		f.OutFilterData <- data
	}
}

//Fetch is the component for applying blockpage and unexpected responses regex matching
type Fetch struct {
	InFetchData <-chan map[string]interface{}
	InBlockpageData     <-chan map[string]*regexp.Regexp
	InFalsePositiveData <-chan map[string]*regexp.Regexp
	HTMLPages 				<-chan map[string]string
	OutFetchData 		chan<- map[string]interface{}

}

//Process applies blockpage and unexpected responses regex matching
func (f *Fetch) Process() {
	blockpages := <-f.InBlockpageData
	falsePositives := <-f.InFalsePositiveData
	html := <-f.HTMLPages
	blockpageOrder := make([]string, len(blockpages))
	i := 0
	for k := range blockpages {
		blockpageOrder[i] = k
		i++
	}

	//Sort to map the blockpage strings in the right order
	sort.Strings(blockpageOrder)

	for data := range f.InFetchData {
		fetched := false
		if len(html) != 0 {
			
			fetched = true
			var confirmed string
			var fingerprint string
			if data["passed"].(bool) == false && data["UntaggedAnswer"].(bool) == false {
				answersMap := data["answers"].(map[string]interface{})
				for answer, _ := range answersMap {
					if body, ok := html[answer+data["query"].(string)]; ok {
						for fp, pattern := range falsePositives {
							if pattern.MatchString(body) {
								confirmed = "false"
								fingerprint = fp
								break
							}
						}
						for _, fp := range blockpageOrder {
							pattern := blockpages[fp]
							if pattern.MatchString(body) {
								confirmed = "true"
								fingerprint = fp
								break
							}
						}
					}
				}
			}
			data["Confirmed"] = confirmed
			data["Fingerprint"] = fingerprint
		}
		data["Fetched"] = fetched
		f.OutFetchData <- data
	}
}

//Verify is the component for applying post procesing hueristics to avoid false positives
type Verify struct {
	InVerifyData  <-chan map[string]interface{}
	CDNIPs <-chan set.Set
	OutVerifyData chan<- map[string]interface{}
}

//Process applies post processig hueristics to avoid false positives
func (v *Verify) Process() {
	cdnIPs := <-v.CDNIPs
	for data := range v.InVerifyData {
		answersMap := data["answers"].(map[string]interface{})
		belongsToCDN := false
		for answer, _ := range answersMap {
			if cdnIPs.Contains(answer) {
				belongsToCDN = true
			}
		}
		data["BelongsToCDN"] = belongsToCDN
		v.OutVerifyData <- data
	}
}

//Analysis is the component that stores the input and output for different types of analysis
type Analysis struct {
	InAnalysisData  <-chan map[string]interface{}
	InAnalysisType  <-chan string
	OutAnalysisData chan<- map[string]interface{}
}

//Process performs analysis on the filtered data to calcuate different types of aggregates
//TODO: Add more analysis types
func (a *Analysis) Process() {
	analysisType := <-a.InAnalysisType
	for data := range a.InAnalysisData {
		if analysisType == "Domain" {
			a.OutAnalysisData <- map[string]interface{}{"Query": data["query"], "Anomaly": (!(data["passed"].(bool)) && !(data["BelongsToCDN"].(bool)) && !(data["UntaggedAnswer"].(bool))), "Fetched": data["Fetched"], "Confirmed": data["Confirmed"], "Country": data["Geolocation"]}
		} else if analysisType == "Vantage Point" {
			a.OutAnalysisData <- map[string]interface{}{"Resolver": data["resolver"], "Anomaly": (!(data["passed"].(bool)) && !(data["BelongsToCDN"].(bool)) && !(data["UntaggedAnswer"].(bool))), "Fetched": data["Fetched"], "Confirmed": data["Confirmed"], "Country": data["Geolocation"]}
		} 
	}

}


//ProcessLine constructs the directed cyclic graph that handles data flow between different components.
func ProcessLine() *goflow.Graph {
	network := goflow.NewGraph()

	//Add network processes
	network.Add("parse", new(Parse))
	network.Add("metadata", new(MetaData))
	network.Add("filter", new(Filter))
	network.Add("verify", new(Verify))
	network.Add("fetch", new(Fetch))
	network.Add("analysis", new(Analysis))


	// Connect them with a channel
	network.Connect("parse", "ResJSONData", "metadata", "InMetaData")
	network.Connect("metadata", "OutMetaData", "filter", "InFilterData")
	network.Connect("filter", "OutFilterData", "fetch", "InFetchData")
	network.Connect("fetch", "OutFetchData", "verify", "InVerifyData")
	network.Connect("verify", "OutVerifyData", "analysis", "InAnalysisData")

	network.MapInPort("Input", "parse", "InLine")
	network.MapInPort("ControlAnswersInput", "filter", "ControlAnswers")
	network.MapInPort("GeolocationInput", "metadata", "Geolocation")
	network.MapInPort("BlockpageInput", "fetch", "InBlockpageData")
	network.MapInPort("FalsePositiveInput", "fetch", "InFalsePositiveData")
	network.MapInPort("HTMLInput", "fetch", "HTMLPages")
	network.MapInPort("CDNIPInput", "verify", "CDNIPs")
	network.MapInPort("AnalysisType","analysis","InAnalysisType")

	//Map the output ports for the network
	network.MapOutPort("ProcessingOutput", "analysis", "OutAnalysisData")

	return network
}

//Prompt gets the user's choice of analysis type
func Prompt() string {

	types := []analyses{
		{Name: "Websites marked as anomaly per country (csv)", Index: "Domain"},
		{Name: "Anomalies per vantage point per country (csv)", Index: "Vantage Point"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U00002192 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U00002192 {{ .Name | red | cyan }}",
		Details: `
--------- Analysis Types ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Index:" | faint }}	{{ .Index }}`,
	}

	prompt := promptui.Select{
		Label:        "Select Type of Analysis",
		Items:        types,
		Templates:    templates,
		CursorPos:    0,
		HideSelected: false,
	}

	i, _, err := prompt.Run()

	if err != nil {
		log.Fatal("Prompt failed: ", err.Error())
	}

	return types[i].Index
}

//AnalyzeSatellite is the main function that handles io and set up of the network
func AnalyzeSatellite(inputFile string, outputFile string, satellitev1HtmlFile string) {
	analysisType := Prompt()
	blockpages, err := hquack.ReadFingerprints("https://assets.censoredplanet.org/blockpage_signatures.json")
	if err != nil {
		log.Fatal("Could not read blockpage data: ", err.Error())
	}
	log.Info("Blockpage read successful")
	falsePositives, err := hquack.ReadFingerprints("https://assets.censoredplanet.org/false_positive_signatures.json")
	if err != nil {
		log.Fatal("Could not read false positive data: ", err.Error())
	}

	//Load answer Tags
	log.Info("Loading answer tags")
	ansTags := make(map[string]*tags)
	cdnIPs := set.NewSet()
	ansTags, cdnIPs = loadAnsTags(inputFile, cdnRegex)

	//Load control answers
	log.Info("Loading control answers with tags")
	controlAnswers := make(map[string]*tagsSet)
	controlAnswers = loadControls(inputFile, ansTags)

	log.Info("Loading geolocation info")
	geolocation := make(map[string]string)
	geolocation = loadGeolocation(inputFile)

	htmlPages := make(map[string]string)
	if satellitev1HtmlFile != "" {
		log.Info("Loading Satellite-v1 HTML page")
		htmlPages = loadHTML(satellitev1HtmlFile)
	} else {
		log.Warn("HTML matching will be skipped")
	}

	log.Info("Going through file: ", inputFile)
	fileReader, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Could not open input file", err.Error())

	}
	defer fileReader.Close()

	//Read the Tar file and get the required files 
	//TODO: Read more than one tar file for files in 2020
	interferenceFileBytes, err := tarballReader.ReadTarball(fileReader, "interference.json")
	if err != nil {
		log.Fatal("Could not read tarball", err.Error())

	}
	log.Info("Interference File read: ", inputFile)
	interferenceFileText := string(interferenceFileBytes)
	interferenceFileTextLines := strings.Split(interferenceFileText, "\n")
	log.Info("Number of lines in Interference file: ", len(interferenceFileTextLines))
	processedData := make([]map[string]interface{}, 0)

	network := ProcessLine()

	ControlAnswersInput := make(chan map[string]*tagsSet)
	GeolocationInput := make(chan map[string]string)
	CDNIPInput := make(chan set.Set)
	HTMLInput := make(chan map[string]string)
	BlockpageInput := make(chan map[string]*regexp.Regexp)
	FalsePositiveInput := make(chan map[string]*regexp.Regexp)
	AnalysisTypeInput := make(chan string)
	ProcessingOutput := make(chan map[string]interface{})
	done := make(chan bool)

	network.SetInPort("ControlAnswersInput", ControlAnswersInput)
	network.SetInPort("GeolocationInput", GeolocationInput)
	network.SetInPort("CDNIPInput", CDNIPInput)
	network.SetInPort("BlockpageInput", BlockpageInput)
	network.SetInPort("FalsePositiveInput", FalsePositiveInput)
	network.SetInPort("HTMLInput", HTMLInput)
	network.SetInPort("AnalysisType", AnalysisTypeInput)
	network.SetOutPort("ProcessingOutput", ProcessingOutput)

	In := make(chan string)
	network.SetInPort("Input", In)

	//Start the network
	wait := goflow.Run(network)

	log.Info("Network set up. Starting data flow.")

	AnalysisTypeInput <- analysisType
	close(AnalysisTypeInput)

	ControlAnswersInput <- controlAnswers
	close(ControlAnswersInput)

	GeolocationInput <- geolocation
	close(GeolocationInput)

	BlockpageInput <- blockpages
	close(BlockpageInput)

	FalsePositiveInput <- falsePositives
	close(FalsePositiveInput)

	HTMLInput <- htmlPages
	close(HTMLInput)

	CDNIPInput <- cdnIPs
	close(CDNIPInput)

	// create and start new progress bar
	bar := pb.StartNew(len(interferenceFileTextLines))

	//Set the receiving channel
	go func() {
		for {
			processingOutputData, more := <-ProcessingOutput
			if more {
				processedData = append(processedData, processingOutputData)
			} else {
				log.Info("Received all the dataflow output")
				done <- true
				return
			}
		}
	}()

	//Send the input "interference.json" lines
	for _, line := range interferenceFileTextLines {
		In <- line
		bar.Increment()
	}

	close(In)

	<-wait
	bar.Finish()

	<-done

	//Analyze the dataflow output to create simplified CSV
	output := Analyze(processedData, analysisType)
	if output == nil {
		log.Warn("Analysis output is empty")
	}

	//Write the CSV to output file
	WriteToCSV(output, outputFile)
}
