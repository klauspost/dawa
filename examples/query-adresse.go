// +build ignore

package main

import (
	"fmt"
	"github.com/klauspost/dawa"
)

func main() {
	// Ask for an address with Vejnavn = Rødkildevej and Husnr = 46
	query := dawa.NewAdresseQuery().Vejnavn("Rødkildevej").Husnr("46")
	fmt.Println("Url:" + query.URL())

	item, err := query.First()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nFirst Result ID: %s\n", item.ID)

	// Note that this will re-run the query:
	all, err := query.All()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nAll Results: %+v\n", all)

}
