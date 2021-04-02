//Copyright 2021 Censored Planet
//analyze-cp provides example analysis for Censored Planet public raw data. 
package main

import (
	"flag"
	"os"
	"regexp"
	"strings"

	"github.com/censoredplanet/censoredplanet/analysis/pkg/geolocate"
	"github.com/censoredplanet/censoredplanet/analysis/pkg/hquack"
	"github.com/censoredplanet/censoredplanet/analysis/pkg/satellite"
	log "github.com/sirupsen/logrus"
)

//Flags stores values from user-entered command line flags
type Flags struct {
	inputFile         string
	outputFile        string
	blockpageFile     string
	falsePositiveFile string
	satellitev1HtmlFile string
	skipDatesFile     string
	mmdbFile          string
	logLevel          uint
	logFileFlag       string
}

//ReadSkipScanDates checkes whether the current scandate of the file matches any of the scandates to skip analysis
func ReadSkipScanDates(url string, technique string, scandate string) bool {
	data := strings.Split(hquack.AssetsHTTPRequest(url), "\n")
	for _, line := range data {
		parts := strings.Split(line, ",")
		if strings.ToLower(parts[0]) == technique && parts[1] == scandate {
			return true
		}
	}
	return false
}

func main() {
	//create Flag state object
	var f Flags

	flag := flag.NewFlagSet("flags", flag.ExitOnError)
	flag.StringVar(&f.inputFile, "input-file", "", "REQUIRED - Input tar.gz file (downloaded from censoredplanet.org)")
	flag.StringVar(&f.outputFile, "output-file", "output.csv", "Output csv file (default - output.csv)")
	flag.StringVar(&f.satellitev1HtmlFile, "satellitev1-html-file", "", "(Optional) json file that contains HTML responses for detecting blockpages from satellitev1 resolved IP addresses. The JSON file should have the following fields: 1) ip (resolved ip), query (query performed by satellitev1), body (HTML body). If unspecified, the blockpage matching process will be skipped.")
	flag.StringVar(&f.mmdbFile, "mmdb-file", "", "REQUIRED - Maxmind Geolocation MMDB file (Download from maxmind.com)")
	flag.UintVar(&f.logLevel, "verbosity", 3, "level of log detail (0-5)")
	flag.StringVar(&f.logFileFlag, "log-file", "-", "file name for logging, (- is stderr)")

	//Parse Flags
	flag.Parse(os.Args[1:])

	//Check required inputs

	if f.inputFile == "" || f.mmdbFile == "" {
		log.Fatal("Please provide required input flags - Input tar.gz file and Maxming mmdb file")
	}

	//Set log file. By default, it is stderr
	logFile := os.Stderr
	if f.logFileFlag != "-" {
		var err error
		if logFile, err = os.Create(f.logFileFlag); err != nil {
			log.Fatal(err)
		}
	}
	log.SetOutput(logFile)

	switch f.logLevel {
	case 1:
		log.SetLevel(log.FatalLevel)
	case 2:
		log.SetLevel(log.ErrorLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	case 4:
		log.SetLevel(log.InfoLevel)
	case 5:
		log.SetLevel(log.DebugLevel)
	default:
		log.Fatal("MAIN: invalid log level")
	}

	//timeout := time.Duration(time.Millisecond * time.Duration(*timeoutFlag))

	log.Info("Logging set up.")

	parts := strings.Split(f.inputFile, "/")
	if len(parts) == 0 {
		log.Fatal("Input file not found")
	}
	filename := parts[len(parts)-1]
	//Compile regex for filename from Censored Planet website
	r := regexp.MustCompile("CP_[a-zA-Z]+[-]*[a-zA-Z]*-20[1-3][0-9]-[0-1][0-9]-[0-3][0-9]-[0-2][0-9]-[0-5][0-9]-[0-5][0-9].tar.gz")
	if !r.MatchString(filename) {
		log.Fatal("Input file does not match expected file name pattern. Please use same file name pattern as in the censoredplanet.org website")
	}

	//Extract the scan technique and scan date
	technique := strings.ToLower(strings.Split(strings.Split(filename, "-")[0], "_")[1])
	protocol := ""
	scandate := ""
	if technique == "quack"{
		protocol = strings.ToLower(strings.Split(filename, "-")[1])
		scandate = strings.Split(filename, "-")[2] + strings.Split(filename, "-")[3] + strings.Split(filename, "-")[4]
	} else if technique == "satellite" {
		scandate = strings.Split(filename, "-")[1] + strings.Split(filename, "-")[2] + strings.Split(filename, "-")[3]
	} else {
		log.Fatal("Unsupported technique for analysis")
	}

	//Should this scan be skipped?
	if skip := ReadSkipScanDates("https://assets.censoredplanet.org/avoid_scandates.txt", technique, scandate); skip == true {
		log.Fatal("This scan is in the do-not-include list.")
	} 

	//Initialize maxmind
	log.Info("Input File okay!")
	err := geolocate.Initialize(f.mmdbFile)
	if err != nil {
		log.Fatal("Could not initialize Maxmind DB: ", err.Error())
	}
	log.Info("Maxmind init success")

	//Start analysis
	if technique == "quack" {
		hquack.AnalyzeHquack(f.inputFile, f.outputFile, protocol)
	} else if technique == "satellite" {
		if scandate >= "20210301" {
			log.Fatal("Satellitev2 support is not provided yet")
		}
		satellite.AnalyzeSatellite(f.inputFile, f.outputFile, f.satellitev1HtmlFile)
	} 
}
