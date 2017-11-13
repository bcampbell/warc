package main

// simple tool to list the records in a WARC file

import (
	"flag"
	"fmt"
	"github.com/bcampbell/warc"
	"io"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [warcfile]...\nList records in WARC files\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		return
	}

	for _, filename := range flag.Args() {
		err := doFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func doFile(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()
	r := warc.NewReader(in)
	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		fmt.Printf("%s %s %d\n", rec.Header.Get("Warc-Type"), rec.Header.Get("Content-Type"), len(rec.Block))

		//		fmt.Printf("--------------------------------------\n%s\n--------------------------", rec.Block)
	}
}
