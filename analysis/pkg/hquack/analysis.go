//Copyright 2020 Censored Planet

// Package hquack contains analysis scripts for quack and hyperquack protocols
package hquack

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
			if _, ok := dataMap[result["Country"].(string)][result["Keyword"].(string)]; !ok {
				dataMap[result["Country"].(string)][result["Keyword"].(string)] = map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["Keyword"].(string)]["Measurements"]; !ok {
				dataMap[result["Country"].(string)][result["Keyword"].(string)]["Measurements"] = 0
			}
			if _, ok := dataMap[result["Country"].(string)][result["Keyword"].(string)]["Anomalies"]; !ok {
				dataMap[result["Country"].(string)][result["Keyword"].(string)]["Anomalies"] = 0
			}
			if _, ok := dataMap[result["Country"].(string)][result["Keyword"].(string)]["Confirmed"]; !ok {
				dataMap[result["Country"].(string)][result["Keyword"].(string)]["Confirmations"] = 0
			}

			dataMap[result["Country"].(string)][result["Keyword"].(string)]["Measurements"]++

			if result["Anomaly"] == true {
				dataMap[result["Country"].(string)][result["Keyword"].(string)]["Anomalies"]++
			}

			if result["Confirmed"] == "true" {
				dataMap[result["Country"].(string)][result["Keyword"].(string)]["Confirmations"]++
			}
		}
	} else if analysisType == "Vantage Point" {
		for _, result := range data {
			if _, ok := dataMap[result["Country"].(string)]; !ok {
				dataMap[result["Country"].(string)] = map[string]map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["Server"].(string)]; !ok {
				dataMap[result["Country"].(string)][result["Server"].(string)] = map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["Server"].(string)]["Measurements"]; !ok {
				dataMap[result["Country"].(string)][result["Server"].(string)]["Measurements"] = 0
			}
			if _, ok := dataMap[result["Country"].(string)][result["Server"].(string)]["Anomalies"]; !ok {
				dataMap[result["Country"].(string)][result["Server"].(string)]["Anomalies"] = 0
			}
			if _, ok := dataMap[result["Country"].(string)][result["Server"].(string)]["Confirmed"]; !ok {
				dataMap[result["Country"].(string)][result["Server"].(string)]["Confirmations"] = 0
			}

			dataMap[result["Country"].(string)][result["Server"].(string)]["Measurements"]++

			if result["Anomaly"] == true {
				dataMap[result["Country"].(string)][result["Server"].(string)]["Anomalies"]++
			}

			if result["Confirmed"] == "true" {
				dataMap[result["Country"].(string)][result["Server"].(string)]["Confirmations"]++
			}
		}
	} else if analysisType == "Error Type" {
		for _, result := range data {
			if _, ok := dataMap[result["Country"].(string)]; !ok {
				dataMap[result["Country"].(string)] = map[string]map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["ErrorType"].(string)]; !ok {
				dataMap[result["Country"].(string)][result["ErrorType"].(string)] = map[string]int{}
			}
			if _, ok := dataMap[result["Country"].(string)][result["ErrorType"].(string)]["Measurements"]; !ok {
				dataMap[result["Country"].(string)][result["ErrorType"].(string)]["Measurements"] = 0
			}
			if _, ok := dataMap[result["Country"].(string)][result["ErrorType"].(string)]["Anomalies"]; !ok {
				dataMap[result["Country"].(string)][result["ErrorType"].(string)]["Anomalies"] = 0
			}
			if _, ok := dataMap[result["Country"].(string)][result["ErrorType"].(string)]["Confirmed"]; !ok {
				dataMap[result["Country"].(string)][result["ErrorType"].(string)]["Confirmations"] = 0
			}

			dataMap[result["Country"].(string)][result["ErrorType"].(string)]["Anomalies"]++

			if result["Confirmed"] == "true" {
				dataMap[result["Country"].(string)][result["ErrorType"].(string)]["Confirmations"]++
			}
		}
	}
	return dataMap

}
