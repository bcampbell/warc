package warc

import (
	//	"fmt"
	"io"
	"os"
	"testing"
)

// wget 1.14+ supports warc output, eg:
//   $ wget --warc-file=example --no-warc-compression http://www.example.com

func TestRead(t *testing.T) {
	f, err := os.Open("example.warc")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rdr := NewReader(f)
	for {
		_, err = rdr.ReadRecord()
		if err != nil {
			if err != io.EOF {
				t.Fatalf("didn't see the expected EOF (got %s)", err)
			}
			break
		}
	}
}
