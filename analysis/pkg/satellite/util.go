//Copyright 2021 Censored Planet

// Package satellite contains analysis scripts for satellite
package satellite

import (
	"strings"
	"os"
	"encoding/json"
	"regexp"
	"bufio"
	"github.com/censoredplanet/censoredplanet/analysis/pkg/tarballReader"
	set "github.com/deckarep/golang-set"	
	log "github.com/sirupsen/logrus"
	
)

type tags struct {
	http   string
	cert   string
	asnum  float64
	asname string
}

type tagsSet struct {
	ip     set.Set
	http   set.Set
	cert   set.Set
	asnum  set.Set
	asname set.Set
}

//newTagsSet creates a new set of tags for Satellite responses
//Returns new tag set
func newTagsSet() *tagsSet {
	t := new(tagsSet)
	t.ip = set.NewSet()
	t.http = set.NewSet()
	t.cert = set.NewSet()
	t.asnum = set.NewSet()
	t.asname = set.NewSet()
	return t
}

//loadAnsTags loads Satellite answers IPs and their tags from "tagged_answers.json"
//Input - input tar.gz file, and cdn regex for finding CDN IPs
//Output - Map of answers IPs and tags, and set of CDN IPs
func loadAnsTags(inputFile string, cdnRegex *regexp.Regexp) (map[string]*tags, set.Set) {
	fileReader, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Could not open input file", err.Error())

	}
	defer fileReader.Close()
	taggedAnswersFileBytes, err := tarballReader.ReadTarball(fileReader, "tagged_answers.json")
	if err != nil {
		log.Fatal("Could not read tarball: ", err.Error())
	}
	log.Info("Tagged Answers File read")
	taggedAnswersFileText := string(taggedAnswersFileBytes)
	taggedAnswersTextLines := strings.Split(taggedAnswersFileText, "\n")
	log.Info("Number of lines in Tagged Answers file: ", len(taggedAnswersTextLines))
	ansTags := make(map[string]*tags)
	cdnIPs := set.NewSet()
	for _, line := range taggedAnswersTextLines {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
			log.Warn("Cannot unmarshal Tagged Answers JSON data: ", line, ", Error: ", err.Error())
			continue
		}
		ip := jsonData["ip"].(string)
		if ansTags[ip] == nil {
			ansTags[ip] = new(tags)
		}
		if jsonData["http"] != nil {
			ansTags[ip].http = jsonData["http"].(string)
		}
		if jsonData["cert"] != nil {
			ansTags[ip].cert = jsonData["cert"].(string)
		}
		if jsonData["asnum"] != nil {
			ansTags[ip].asnum = jsonData["asnum"].(float64)
		}
		if jsonData["asname"] != nil {
			ansTags[ip].asname = jsonData["asname"].(string)
			if cdnRegex.MatchString(jsonData["asname"].(string)) {
				cdnIPs.Add(ip)
			}
		}
	}
	return ansTags, cdnIPs
}

//loadControls loads DNS resolution answer from Satellite's control resolvers
//Input - input tar.gz file, and tagged answers
//Output - Tagged control resolver answers
func loadControls(inputFile string, ansTags map[string]*tags) map[string]*tagsSet {
	fileReader, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Could not open input file", err.Error())

	}
	defer fileReader.Close()
	controlAnswersFileBytes, err := tarballReader.ReadTarball(fileReader, "answers_control.json")
	if err != nil {
		log.Fatal("Could not read tarball: ", err.Error())

	}
	log.Info("Control Answers File read")
	controlAnswersFileText := string(controlAnswersFileBytes)
	controlAnswersTextLines := strings.Split(controlAnswersFileText, "\n")
	log.Info("Number of lines in Control Answers file: ", len(controlAnswersTextLines))	
	controls := make(map[string]*tagsSet)
	for _, line := range controlAnswersTextLines {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
			log.Warn("Cannot unmarshal Tagged Answers JSON data: ", line, ", Error: ", err.Error())
			continue
		}
		query := jsonData["query"].(string)
		answers := jsonData["answers"].([]interface{})
		if controls[query] == nil {
			controls[query] = newTagsSet()
		}

		for _, answer := range answers {
			controls[query].ip.Add(answer.(string))
			// Add the tags corresponding to this IP answer to control set
			if t, ok := ansTags[answer.(string)]; ok {
				if t.http != "" {
					controls[query].http.Add(t.http)
				}
				if t.cert != "" {
					controls[query].cert.Add(t.cert)
				}
				if t.asnum != 0 {
					controls[query].asnum.Add(t.asnum)
				}
				if t.asname != "" {
					controls[query].asname.Add(t.asname)
				}
			}
		}
	}
	return controls
}

//loadGeolocation gets country information from the "tagged_resolvers.json" file
//Input - tar.gz file
//Output - Geolocation data
func loadGeolocation(inputFile string) map[string]string {
	fileReader, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Could not open input file", err.Error())

	}
	defer fileReader.Close()
	taggedResolversFileBytes, err := tarballReader.ReadTarball(fileReader, "tagged_resolvers.json")
	if err != nil {
		log.Fatal("Could not read tarball: ", err.Error())

	}
	log.Info("Tagged resolvers File read")
	taggedResolversFileText := string(taggedResolversFileBytes)
	taggedResolversTextLines := strings.Split(taggedResolversFileText, "\n")
	log.Info("Number of lines in Tagged resolvers file: ", len(taggedResolversTextLines))	
	geolocation := make(map[string]string)
	for _, line := range taggedResolversTextLines {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
			log.Warn("Cannot unmarshal Tagged resolvers JSON data: ", line, ", Error: ", err.Error())
			continue
		}
		geolocation[jsonData["resolver"].(string)] = jsonData["country"].(string)
	}
	return geolocation
}

//loadHTML loads HTML data from provided input file
//Input - satellitev1HTMLFile from user input, a JSON file containing ip, query, and HTML body
//Returns - Map of answer IP, query, and HTML body
func loadHTML(inputFile string) map[string]string {
	fileReader, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Could not open input file", err.Error())
	}
	defer fileReader.Close()
	s := bufio.NewScanner(fileReader)
	htmlData := make(map[string]string)
	for s.Scan() {
		data := make(map[string]interface{})
		if err := json.Unmarshal(s.Bytes(), &data); err != nil {
			log.Warn("Cannot unmarshal HTML JSON data: ", data, ", Error: ", err.Error())
			continue
		}
		htmlData[data["ip"].(string)+data["query"].(string)] = data["body"].(string)
	}
	return htmlData
}