//Copyright 2021 Censored Planet

// Package satellite contains analysis scripts for satellite data
package satellite

import (
	"encoding/csv"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

//WriteToCSV writes a map of analyzed data to a CSV file.
//Input - The data to write (map[string]map[string]map[string]int), the analysis type ("Vantage Point", "Domain"), Output csv filename
//Returns - None
//WriteToCSV expects the data to be in a certain nested map format.
func WriteToCSV(dataMap map[string]map[string]map[string]int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create output file: ", err.Error())
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"Country", "Vantage Point/Website", "Measurements", "Anomalies", "Confirmations"})

	for country, innerMap := range dataMap {
		for targetValue, statsMap := range innerMap {
			if _, ok := statsMap["Confirmations"]; ok {
				err := w.Write([]string{fmt.Sprintf("%v", country), fmt.Sprintf("%v", targetValue), fmt.Sprintf("%v", statsMap["Measurements"]), fmt.Sprintf("%v", statsMap["Anomalies"]), fmt.Sprintf("%v", statsMap["Confirmations"])})
				if err != nil {
					log.Warn("Could not write row due to: ", err.Error())
				}
			} else {
				err := w.Write([]string{fmt.Sprintf("%v", country), fmt.Sprintf("%v", targetValue), fmt.Sprintf("%v", statsMap["Measurements"]), fmt.Sprintf("%v", statsMap["Anomalies"]), "N/A"})
				if err != nil {
					log.Warn("Could not write row due to: ", err.Error())
				}
			}
		}
	}
}
