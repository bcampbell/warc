package warc

import (
	//	"fmt"
	"io"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	f, err := os.Open("foo.warc")
	if err != nil {
		panic(err)
	}

	_, err = Read(f)
	if err != nil {
		panic(err)
	}
	//	fmt.Printf("version: %s\n", rec.Version)
	//	fmt.Printf("%d headers\n", len(rec.Header))
	//	fmt.Printf("block is %d bytes\n", len(rec.Block))

	_, err = Read(f)
	if err != io.EOF {
		t.Fatalf("didn't see the expected EOF (got %s)", err)
	}
}
