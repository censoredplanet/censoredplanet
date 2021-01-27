//Copyright 2020 Censored Planet

// Package hquack contains analysis scripts for quack and hyperquack protocols
package hquack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//AssetsHTTPRequest performs a HTTP request for a Censored Planet asset and returns the body of the response
func AssetsHTTPRequest(url string) string {
	client := http.Client{
		Timeout: time.Second * 5, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Could not create asset HTTP request: ", url, " due to: ", err)
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal("Could not create get asset: ", url, " due to: ", err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal("Could not read asset body: ", url, "due to: ", readErr)
	}
	return string(body)

}

//ReadFingerprints reads blockpage and false positive fingerprints into memory from a URL
func ReadFingerprints(url string) (map[string]*regexp.Regexp, error) {
	var signatures map[string]string
	fingerprintData := make(map[string]*regexp.Regexp)
	data := strings.Split(AssetsHTTPRequest(url), "\n")
	for _, line := range data {
		if line == "\n" || line == "" {
			continue
		}
		if err := json.Unmarshal([]byte(line), &signatures); err != nil {
			return nil, err
		}
		pattern := regexp.QuoteMeta(signatures["pattern"])
		pattern = strings.ReplaceAll(pattern, "%", ".*")
		fingerprintData[signatures["fingerprint"]] = regexp.MustCompile(pattern)
	}
	return fingerprintData, nil
}
