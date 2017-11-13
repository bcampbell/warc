package warc

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// read an http response from a WARC file
// if filename has .gz suffix, gzip is assumed
func ReadFile(filename string) (*http.Response, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var in io.Reader
	if filepath.Ext(filename) == ".gz" {
		gin, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gin.Close()
		in = gin
	} else {
		in = f
	}

	return Read(in)
}

// read an http response from an io.Reader
func Read(in io.Reader) (*http.Response, error) {
	warcReader := NewReader(in)
	for {
		//	fmt.Printf("WARC\n")
		rec, err := warcReader.ReadRecord()
		if err != nil {
			return nil, err
		}
		if rec.Header.Get("Warc-Type") != "response" {
			continue
		}
		//reqURL := rec.Header.Get("Warc-Target-Uri")
		// parse response, grab raw html
		rdr := bufio.NewReader(bytes.NewReader(rec.Block))
		response, err := http.ReadResponse(rdr, nil)
		return response, err
	}
}
