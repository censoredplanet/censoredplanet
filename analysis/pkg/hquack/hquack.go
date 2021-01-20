//Copyright 2020 Censored Planet

// Package hquack contains analysis scripts for quack and hyperquack protocols
package hquack

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/censoredplanet/censoredplanet/analysis/pkg/geolocate"
	"github.com/censoredplanet/censoredplanet/analysis/pkg/tarballReader"
	"github.com/cheggaaa/pb/v3"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/trustmaster/goflow"
)

//analysisType holds the type of analysis input from the user
var analysisType string

//Different error types for interference in application layer.
//TODO: Need to update error codes that appear less frequently
var resetError = regexp.MustCompile("connection reset by peer|EOF")
var timeoutError = regexp.MustCompile("Client.Timeout")
var blockpageError = regexp.MustCompile("Incorrect web response")

//For now, we are not confirmng any other error responses, as they need to be checked in a more case-by-case basis. They will be counted as anomalies.
//var otherError = regexp.MustCompile("response missing Location header|malformed HTTP status code")

//These error types indicate network errors.
var filterError = regexp.MustCompile("too many open files|address already in use|no route to host|connection refused|connect: network is unreachable|connect: connection timed out|getsockopt: network is unreachable|remote error: internal error|trailer header without chunked transfer encoding|remote error: bad record MAC|remote error: handshake failure|remote error: alert(112)|local error:")

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

//Filter is the component for filtering failed measurements
type Filter struct {
	InFilterData  <-chan map[string]interface{}
	OutFilterData chan<- map[string]interface{}
}

//Process filters out failed measurements
//Either FailSanity is set to true, or the measurement error matches one of the filter error types
func (f *Filter) Process() {
	for data := range f.InFilterData {
		filterLine := false
		if data["FailSanity"].(bool) == true {
			filterLine = true
		}
		resultsSlice := data["Results"].([]interface{})
		for _, result := range resultsSlice {
			if result.(map[string]interface{})["Success"].(bool) == false && filterError.MatchString(result.(map[string]interface{})["Error"].(string)) {
				filterLine = true
			}
		}
		if filterLine == false {
			f.OutFilterData <- data
		}
	}
}

//MetaData is the component that assigns measurement metadata to each row
type MetaData struct {
	InMetaData  <-chan map[string]interface{}
	OutMetaData chan<- map[string]interface{}
}

//Process assigns measurement metadata to each row
func (m *MetaData) Process() {
	for data := range m.InMetaData {
		vantagePoint := data["Server"].(string)
		geolocation, err := geolocate.Geolocate(vantagePoint)
		if err != nil {
			log.Warn("Cannot geolocate vantage point: ", err.Error())
		} else {
			data["Geolocation"] = geolocation.CountryCode
			m.OutMetaData <- data
		}
	}
}

//Anomaly is the component that holds the type of anomaly, and whether the anomaly is confirmed or not
type Anomaly struct {
	InAnomalyData       <-chan map[string]interface{}
	OutAnomalyData      chan<- map[string]interface{}
	InBlockpageData     <-chan map[string]*regexp.Regexp
	InFalsePositiveData <-chan map[string]*regexp.Regexp
	InTechnique         <-chan string
}

//Process analyses the type of blocking and determines whether the anomaly is confirmed or not
func (a *Anomaly) Process() {
	blockpages := <-a.InBlockpageData
	falsePositives := <-a.InFalsePositiveData
	technique := <-a.InTechnique

	blockpageOrder := make([]string, len(blockpages))
	i := 0
	for k := range blockpages {
		blockpageOrder[i] = k
		i++
	}

	//Sort to map the blockpage strings in the right order
	sort.Strings(blockpageOrder)

	for data := range a.InAnomalyData {
		var errorType string
		var confirmed string
		var fingerprint string
		if data["Blocked"] == true {
			resultsSlice := data["Results"].([]interface{})
			for _, result := range resultsSlice {
				if result.(map[string]interface{})["Success"].(bool) == false {
					if resetError.MatchString(result.(map[string]interface{})["Error"].(string)) {
						if confirmed != "false" {
							if errorType != "" && errorType != "Reset" {
								errorType = "Mixed"
							} else {
								errorType = "Reset"
							}
							confirmed = "true"
						}
					}
					if timeoutError.MatchString(result.(map[string]interface{})["Error"].(string)) {
						if confirmed != "false" {
							if errorType != "" && errorType != "Timeout" {
								errorType = "Mixed"
							} else {
								errorType = "Timeout"
							}
							confirmed = "true"
						}
					}
					//TODO: Status line errors will only be confirmed if the body has a blockpage. Could there be a set of fingerpritns for status lines too??
					if blockpageError.MatchString(result.(map[string]interface{})["Error"].(string)) {
						var body string
						if technique == "echo" || technique == "discard" {
							body = result.(map[string]interface{})["Received"].(string)
						} else {
							body = result.(map[string]interface{})["Received"].(map[string]interface{})["body"].(string)
						}

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
						if confirmed == "true" {
							if errorType != "" && errorType != "Blockpage" {
								errorType = "Mixed"
							} else {
								errorType = "Blockpage"
							}
						}
					}
					//TODO: Confirm other errors by default or no? Right now, take the more conservative approach and do not confirm.
					if confirmed == "" {
						errorType = "Other"
					}
				}
			}
		}
		data["ErrorType"] = errorType
		data["Confirmed"] = confirmed
		data["Fingerprint"] = fingerprint
		a.OutAnomalyData <- data
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
			a.OutAnalysisData <- map[string]interface{}{"Keyword": data["Keyword"], "Anomaly": data["Blocked"], "Confirmed": data["Confirmed"], "Country": data["Geolocation"]}
		} else if analysisType == "Vantage Point" {
			a.OutAnalysisData <- map[string]interface{}{"Server": data["Server"], "Anomaly": data["Blocked"], "Confirmed": data["Confirmed"], "Country": data["Geolocation"]}

		} else if analysisType == "Error Type" {
			if data["Blocked"] == true {
				a.OutAnalysisData <- map[string]interface{}{"ErrorType": data["ErrorType"], "Confirmed": data["Confirmed"], "Country": data["Geolocation"]}
			}
		}

	}

}

//ProcessLine constructs the directed cyclic graph that handles data flow between different components.
func ProcessLine() *goflow.Graph {
	network := goflow.NewGraph()

	//Add network processes
	network.Add("parse", new(Parse))
	network.Add("filter", new(Filter))
	network.Add("metadata", new(MetaData))
	network.Add("anomaly", new(Anomaly))
	network.Add("analysis", new(Analysis))

	// Connect them with a channel
	network.Connect("parse", "ResJSONData", "filter", "InFilterData")
	network.Connect("filter", "OutFilterData", "metadata", "InMetaData")
	network.Connect("metadata", "OutMetaData", "anomaly", "InAnomalyData")
	network.Connect("anomaly", "OutAnomalyData", "analysis", "InAnalysisData")

	//Map the input ports for the network
	network.MapInPort("BlockpageInput", "anomaly", "InBlockpageData")
	network.MapInPort("FalsePositiveInput", "anomaly", "InFalsePositiveData")
	network.MapInPort("TechniqueInput", "anomaly", "InTechnique")
	network.MapInPort("AnalysisType", "analysis", "InAnalysisType")
	network.MapInPort("Input", "parse", "InLine")

	//Map the output ports for the network
	network.MapOutPort("ProcessingOutput", "analysis", "OutAnalysisData")

	return network
}

//Prompt gets the user's choice of anaylisys type
func Prompt() string {

	types := []analyses{
		{Name: "Blocked websites per country (csv)", Index: "Domain"},
		{Name: "Blocking per vantage point per country (csv)", Index: "Vantage Point"},
		{Name: "Blocking type per country (csv)", Index: "Error Type"},
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

//AnalyzeHquack is the main function that handles io and set up of the network
func AnalyzeHquack(inputFile string, outputFile string, technique string) {

	analysisType := Prompt()
	blockpages, err := ReadFingerprints("https://assets.censoredplanet.org/blockpage_signatures.json")
	if err != nil {
		log.Fatal("Could not read blockpage data: ", err.Error())
	}
	log.Info("Blockpage read successful")
	falsePositives, err := ReadFingerprints("https://assets.censoredplanet.org/false_positive_signatures.json")
	if err != nil {
		log.Fatal("Could not read false positive data: ", err.Error())
	}
	log.Info("False positive read successful")
	log.Info("Going through file: ", inputFile)
	fileReader, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Could not open onput file", err.Error())

	}
	defer fileReader.Close()

	//Read the Tar file
	//TODO: Read more than one tar file for files in 2020
	fileBytes, err := tarballReader.ReadTarball(fileReader)
	if err != nil {
		log.Fatal("Could not read tarball", err.Error())

	}
	log.Info("File read: ", inputFile)
	fileText := string(fileBytes)
	fileTextLines := strings.Split(fileText, "\n")
	log.Info("Number of lines in file: ", len(fileTextLines))

	processedData := make([]map[string]interface{}, 0)

	network := ProcessLine()

	BlockpageInput := make(chan map[string]*regexp.Regexp)
	FalsePositiveInput := make(chan map[string]*regexp.Regexp)
	TechniqueInput := make(chan string)
	AnalysisTypeInput := make(chan string)
	ProcessingOutput := make(chan map[string]interface{})
	done := make(chan bool)

	//Set all of the network ports
	network.SetInPort("BlockpageInput", BlockpageInput)
	network.SetInPort("FalsePositiveInput", FalsePositiveInput)
	network.SetInPort("TechniqueInput", TechniqueInput)
	network.SetInPort("AnalysisType", AnalysisTypeInput)
	network.SetOutPort("ProcessingOutput", ProcessingOutput)

	In := make(chan string)
	network.SetInPort("Input", In)

	//Start the network
	wait := goflow.Run(network)

	log.Info("Network set up. Starting data flow.")

	//Send the input
	BlockpageInput <- blockpages
	close(BlockpageInput)
	FalsePositiveInput <- falsePositives
	close(FalsePositiveInput)
	TechniqueInput <- technique
	close(TechniqueInput)
	AnalysisTypeInput <- analysisType
	close(AnalysisTypeInput)

	// create and start new progress bar
	bar := pb.StartNew(len(fileTextLines))

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

	for _, line := range fileTextLines {
		In <- line
		bar.Increment()
	}

	close(In)

	<-wait
	bar.Finish()

	<-done

	output := Analyze(processedData, analysisType)
	if output == nil {
		log.Warn("Analysis output is empty")
	}

	//Write the output
	WriteToCSV(output, analysisType, outputFile)
}
