//Copyright 2021 Censored Planet

// Package satellite contains analysis scripts for satellite data
package satellite

//Analyze processes data based on the type of analysis specified
//Input - The data ([]map[string]interface{}), the analysisType (Specified by the user prompt)
//Returns - The stats about the data (map[string]map[string]map[string]int)
func Analyze(data []map[string]interface{}, analysisType string) map[string]map[string]map[string]int {
	var dataMap = map[string]map[string]map[string]int{}
	if analysisType == "Domain" {
		for _, result := range data {
			if _, ok := dataMap[result["Country"].(string)]; !ok {
				dataMap[result["Country"].(string)] = map[string]map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["Query"].(string)]; !ok {
				dataMap[result["Country"].(string)][result["Query"].(string)] = map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["Query"].(string)]["Measurements"]; !ok {
				dataMap[result["Country"].(string)][result["Query"].(string)]["Measurements"] = 0
			}
			if _, ok := dataMap[result["Country"].(string)][result["Query"].(string)]["Anomalies"]; !ok {
				dataMap[result["Country"].(string)][result["Query"].(string)]["Anomalies"] = 0
			}
			if result["Fetched"] == true {
				if _, ok := dataMap[result["Country"].(string)][result["Query"].(string)]["Confirmations"]; !ok {
					dataMap[result["Country"].(string)][result["Query"].(string)]["Confirmations"] = 0
				}
			}

			dataMap[result["Country"].(string)][result["Query"].(string)]["Measurements"]++

			if result["Anomaly"] == true {
				dataMap[result["Country"].(string)][result["Query"].(string)]["Anomalies"]++
			}
			if result["Fetched"] == true {
				if result["Confirmed"] == "true" {
					dataMap[result["Country"].(string)][result["Query"].(string)]["Confirmations"]++
				}
			}
		}
	} else if analysisType == "Vantage Point" {
		for _, result := range data {
			if _, ok := dataMap[result["Country"].(string)]; !ok {
				dataMap[result["Country"].(string)] = map[string]map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["Resolver"].(string)]; !ok {
				dataMap[result["Country"].(string)][result["Resolver"].(string)] = map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["Resolver"].(string)]["Measurements"]; !ok {
				dataMap[result["Country"].(string)][result["Resolver"].(string)]["Measurements"] = 0
			}
			if _, ok := dataMap[result["Country"].(string)][result["Resolver"].(string)]["Anomalies"]; !ok {
				dataMap[result["Country"].(string)][result["Resolver"].(string)]["Anomalies"] = 0
			}
			if result["Fetched"] == true {
				if _, ok := dataMap[result["Country"].(string)][result["Resolver"].(string)]["Confirmations"]; !ok {
					dataMap[result["Country"].(string)][result["Resolver"].(string)]["Confirmations"] = 0
				}
			}

			dataMap[result["Country"].(string)][result["Resolver"].(string)]["Measurements"]++

			if result["Anomaly"] == true {
				dataMap[result["Country"].(string)][result["Resolver"].(string)]["Anomalies"]++
			}
			if result["Fetched"] == true {
				if result["Confirmed"] == "true" {
					dataMap[result["Country"].(string)][result["Resolver"].(string)]["Confirmations"]++
				}
			}
		}
	}
	return dataMap

}
