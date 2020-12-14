//Copyright 2020 Censored Planet
//Tests for tarballReader.go
package tarballReader

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestReadTarball(t *testing.T) {
	testResultsFile, err := os.Open("test_results.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer testResultsFile.Close()
	testNoResultsFile, err := os.Open("test_no_results.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer testNoResultsFile.Close()
	testGzGile, err := os.Open("test_file.txt.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer testGzGile.Close()
	testTxtFile, err := os.Open("test_file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer testTxtFile.Close()
	expectedResults := []byte{82, 101, 115, 117, 108, 116, 115, 32, 101, 120, 97, 109, 112, 108, 101, 32, 104, 101, 114, 101, 10}

	var tests = []struct {
		testName string
		ioReader io.Reader
		bytes    []byte
		err      error
	}{
		{"results", testResultsFile, expectedResults, nil},
		{"no results", testNoResultsFile, nil, errors.New("Results file not found")},
		{"Only gzip", testGzGile, nil, errors.New("unexpected EOF")},
		{"Text File", testTxtFile, nil, errors.New("gzip: invalid header")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.testName)
		t.Run(testname, func(t *testing.T) {
			fileBytes, err := ReadTarball(tt.ioReader)
			if err != nil && tt.err == nil {
				t.Errorf("Received Error when none was wanted")
			}
			if err == nil && tt.err != nil {
				t.Errorf("Did not receive error when required")
			}
			if err == nil && bytes.Compare(fileBytes, tt.bytes) != 0 {
				t.Errorf("Wrong file byte sequence read: Observed: %v, Expected: %v", fileBytes, tt.bytes)
			}
			if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Wrong error returned. Observed: %v, Expected: %v", err.Error(), tt.err.Error())
			}
		})
	}
}
