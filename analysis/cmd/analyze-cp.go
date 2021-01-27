//Copyright 2020 Censored Planet

package main

import (
	"flag"
	"os"
	"regexp"
	"strings"

	"github.com/censoredplanet/censoredplanet/analysis/pkg/geolocate"
	"github.com/censoredplanet/censoredplanet/analysis/pkg/hquack"
	log "github.com/sirupsen/logrus"
)

//Flags stores values from user-entered command line flags
type Flags struct {
	inputFile         string
	outputFile        string
	blockpageFile     string
	falsePositiveFile string
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
	r := regexp.MustCompile("CP_[a-zA-Z]+-[a-zA-Z]+-20[1-3][0-9]-[0-1][0-9]-[0-3][0-9]-[0-2][0-9]-[0-5][0-9]-[0-5][0-9].tar.gz")
	if !r.MatchString(filename) {
		log.Fatal("Input file does not match expected file name pattern. Please use same file name pattern as in the censoredplanet.org website")
	}

	//Extract the scan technique and scan date
	technique := strings.ToLower(strings.Split(filename, "-")[1])
	scandate := strings.Split(filename, "-")[2] + strings.Split(filename, "-")[3] + strings.Split(filename, "-")[4]

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
	if technique == "echo" || technique == "discard" || technique == "http" || technique == "https" {
		hquack.AnalyzeHquack(f.inputFile, f.outputFile, technique)
	} else if technique == "satellite" {
		log.Fatal("Support for Satellite analysis is coming soon")
	} else {
		log.Fatal("Unsupported technique for analysis")
	}
}
