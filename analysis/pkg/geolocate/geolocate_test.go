//Copyright 2020 Censored Planet

//Tests for tarballReader.go
package geolocate

import (
	"errors"
	"fmt"
	"testing"

	"github.com/censoredplanet/hyperquackv2/geolocate"
)

func TestGeolocate(t *testing.T) {
	var tests = []struct {
		testName    string
		ip          string
		geolocation string
		err         error
	}{
		{"Invalid IP", "300.300.300.300", "", errors.New("invalid IP address \"300.300.300.300\"")},
		{"MMDB uninitialized", "141.212.123.125", "", errors.New("GeoIp2 reader uninitialized")},
		{"Correct geolocation", "141.212.123.125", "US", nil},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.testName)
		if testname == "Correct geolocation" {
			geolocate.Initialize("/home/ram/repos/src/github.com/censoredplanet/censoredplanet-scheduler/maxmind/GeoLite2-City.mmdb")
		}
		t.Run(testname, func(t *testing.T) {
			geolocation, err := geolocate.Geolocate(tt.ip)
			if err != nil && tt.err == nil {
				t.Errorf("Received Error when none was wanted")
			}
			if err == nil && tt.err != nil {
				t.Errorf("Did not receive error when required")
			}
			if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Wrong error returned. Observed: %v, Expected: %v", err.Error(), tt.err.Error())
			}
			if err == nil && geolocation.CountryCode != tt.geolocation {
				t.Errorf("Wrong geolocation: Observed: %v, Expected: %v", geolocation.CountryCode, tt.geolocation)
			}
		})
	}

}
