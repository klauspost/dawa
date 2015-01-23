// +build ignore

package main

import (
	"fmt"
	"github.com/klauspost/dawa"
	"io"
	"log"
	"net/http"
	"time"
)

// Makes a request to the web api and prints the result.
func main() {
	// Get an array of addresses
	resp, err := http.Get("http://dawa.aws.dk/supplerendebynavne")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Send the response stream to a decoder
	iter, err := dawa.ImportSupplBynavnJSON(resp.Body)
	if err != nil {
		panic(err)
	}

	// Read responses one by one
	n := 0
	t := time.Now()
	for {
		a, err := iter.Next()
		if err == io.EOF {
			log.Printf("Finished reading %d entries in %v.\n", n, time.Now().Sub(t))
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		n++
		fmt.Printf("Entry:%#v\n", a)
	}
}
