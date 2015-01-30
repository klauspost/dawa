// +build ignore

package main

import (
	"fmt"
	"github.com/klauspost/dawa"
	"io"
)

func main() {
	// Ask for kommuner that start with "aa"
	query := dawa.NewListQuery("kommuner", false).Q("aa*")
	fmt.Println("Url: " + query.URL())

	iter, err := query.Iter()
	if err != nil {
		panic(err)
	}

	// Close the iterator when done.
	defer iter.Close()

	// Iterate the results.
	for {
		item, err := iter.NextKommune()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Item: %s\n", item.Navn)
	}

	fmt.Printf("Finished.\n")

}
