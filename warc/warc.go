package warc

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
)

type WARCRecord struct {
	Version string

	// Header contains the WARC headers fields.
	// Note that the names are canonicalised, so
	// use "Warc-Target-Uri" instead of "WARC-Target-URI", for example.
	Header textproto.MIMEHeader

	// the payload data
	Block []byte
}

// TODO:
//  - writing
//  - helper for setting up gzip support

type WARCReader struct {
	rdr *textproto.Reader
}

func NewWARCReader(in io.Reader) *WARCReader {
	bufin := bufio.NewReader(in)
	rdr := textproto.NewReader(bufin)
	r := &WARCReader{rdr: rdr}
	return r
}

// Reads the next WARC record.
// nil,io.EOF is returned if no more records are available.
func (r *WARCReader) Read() (*WARCRecord, error) {
	rdr := r.rdr

	// read the version
	ver, err := rdr.ReadLine()
	if err == io.EOF {
		// graceful exit - no more records
		return nil, io.EOF
	}
	if err != nil {
		return nil, fmt.Errorf("couldn't read version: %s", err)
	}
	if ver != "WARC/1.0" {
		return nil, fmt.Errorf("unknown version: '%s'", ver)
	}
	out := &WARCRecord{Version: ver}

	// read the header pairs
	out.Header, err = rdr.ReadMIMEHeader()
	if err != nil {
		return nil, fmt.Errorf("couldn't read header: %s", err)
	}

	// read the payload
	var length int
	if foo := out.Header.Get("Content-Length"); foo != "" {
		length, err = strconv.Atoi(foo)
		if err != nil {
			return nil, fmt.Errorf("bad Content-Length: %s", err)
		}
	} else {
		return nil, fmt.Errorf("record is missing Content-Length header")
	}

	out.Block = make([]byte, length)
	_, err = io.ReadFull(rdr.R, out.Block)
	if err != nil {
		return nil, fmt.Errorf("error reading block: %s", err)
	}

	// two CRLF to finish off
	for i := 0; i < 2; i++ {
		blank, err := rdr.ReadLine()
		if err != nil || blank != "" {
			return nil, fmt.Errorf("Missing blank lines after block")
		}
	}
	return out, nil
}
