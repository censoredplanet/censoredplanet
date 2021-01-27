//Copyright 2020 Censored Planet

// Package geolocate performs IP geolocation
package geolocate

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	geoip2 "github.com/oschwald/geoip2-golang"
)

var geoIP2Reader *geoip2.Reader

// Initialize starts the geolocation service with the MaxMind MMDB file
func Initialize(databasePath string) error {
	reader, err := geoip2.Open(databasePath)
	if err == nil {
		geoIP2Reader = reader
	}
	return err
}

//IPGeolocation stores the country information for an IP
type IPGeolocation struct {
	CountryName string `json:"country_name,omitempty"`
	// ISO 3166 2-letter uppercase contry code
	CountryCode string `json:"country_code,omitempty"`
}

const language = "en"

//Geolocate uses maxmind to geolocate IP to country
func Geolocate(ip string) (*IPGeolocation, error) {
	address := net.ParseIP(ip)
	if address == nil {
		return nil, errors.New(fmt.Sprintf("invalid IP address \"%s\"", ip))
	}
	if geoIP2Reader == nil {
		return nil, errors.New("GeoIp2 reader uninitialized")
	}
	data, err := geoIP2Reader.City(address)
	if err != nil {
		return nil, err
	}
	countryCode := strings.ToUpper(data.Country.IsoCode)
	if countryCode != "" && !IsCountryCode(countryCode) {
		log.Printf("Error: unrecognized country code \"%s\"", countryCode)
	}
	return &IPGeolocation{
		CountryName: data.Country.Names[language],
		CountryCode: countryCode,
	}, nil
}
