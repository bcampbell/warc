package warc

// helpers to write out raw HTTP requests/responses as noddy .warc files

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

// copy a Response, leaving a new Body reader which contains
// the slurped data.
func copyResponse(orig *http.Response) (*http.Response, error) {
	// read body first
	bod, err := ioutil.ReadAll(orig.Body)
	if err != nil {
		return nil, err
	}
	orig.Body.Close()
	orig.Body = nopCloser{bytes.NewReader(bod)}

	clone := *orig
	clone.Body = nopCloser{bytes.NewReader(bod)}

	return &clone, nil
}

// WriteWARC writes out an http response (including it's body).
// It tries to leave the response unaltered, although it works by
// reading in the entire Body, replacing it with a []byte-backed
// reader reset back to the beginning. This should be fine for most
// applications, just be aware that this means it's not 100% non-intrusive.
// TODO: pass in optional extra headers instead of srcURL
func Write(w io.Writer, resp *http.Response, srcURL string, timeStamp time.Time) error {

	// copy the response so we can peek at the body
	tmpResp, err := copyResponse(resp)
	if err != nil {
		return err
	}

	var payload bytes.Buffer
	err = tmpResp.Write(&payload)
	if err != nil {
		return err
	}

	warcHdr := http.Header{}
	// required fields
	warcHdr.Set("WARC-Record-ID", fmt.Sprintf("urn:X-scrapeomat:%d", time.Now().UnixNano()))
	warcHdr.Set("Content-Length", fmt.Sprintf("%d", payload.Len()))
	warcHdr.Set("WARC-Date", timeStamp.UTC().Format(time.RFC3339))
	warcHdr.Set("WARC-Type", "response")
	// some extras

	warcHdr.Set("WARC-Target-URI", tmpResp.Request.URL.String())
	// cheesy custom field for original url, in case we were redirected
	warcHdr.Set("X-Scrapeomat-Srcurl", srcURL)
	//	warcHdr.Set("WARC-IP-Address", "")
	warcHdr.Set("Content-Type", "application/http; msgtype=response")

	fmt.Fprintf(w, "WARC/1.0\r\n")
	err = warcHdr.Write(w)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\r\n")
	_, err = payload.WriteTo(w)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\r\n")
	fmt.Fprintf(w, "\r\n")

	return nil
}
