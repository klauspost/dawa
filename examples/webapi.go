// +build ignore
package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/klauspost/dawa"
	"io/ioutil"
	"net/http"
)

// Uses spew for nicer formatting of output, to use, execute:
//
// go get github.com/davecgh/go-spew/spew
//
// to install.
func main() {
	// Get an address
	resp, err := http.Get("http://dawa.aws.dk/adresser/0a3f509d-07be-32b8-e044-0003ba298018")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response into r
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Create destination and decode the response
	add := dawa.Adresse{}
	err = json.Unmarshal(r, &add)
	if err != nil {
		panic(err)
	}

	// We got the response. Print it.
	fmt.Printf("JSON response:%s\nDecoded:\n", string(r))
	spew.Dump(add)
}
