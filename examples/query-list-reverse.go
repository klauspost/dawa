// +build ignore

package main

import (
	"fmt"
	"github.com/klauspost/dawa"
	"io"
)

func main() {
	// Ask for region at 12.5851471984198 y=55.6832383751223
	iter, err := dawa.NewReverseQuery("kommuner", 12.5851471984198, 55.6832383751223, "")
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
		fmt.Printf("Result:\n%#v\n", item)
	}

	fmt.Printf("Finished.\n")

}
