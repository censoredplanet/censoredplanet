//Copyright 2021 Censored Planet

//Package tarballReader tar.gz raw data files that are published on the Censored Planet website.
package tarballReader

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

//ReadTarball reads the tar.gz file and sends back the file bytes
//Input - tar.gz file , Filename
//Output - File byte stream, error
func ReadTarball(reader io.Reader, filename string) ([]byte, error) {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, errors.New("File not found")
		}
		if err != nil {
			return nil, err
		}
		if strings.Contains(hdr.Name, filename) {
			bar := pb.New(int(hdr.Size))
			bar.Start()
			barReader := bar.NewProxyReader(tr)
			bs, err := ioutil.ReadAll(barReader)
			if err != nil {
				return nil, err
			}
			bar.Finish()
			return bs, nil
		}
	}
}
