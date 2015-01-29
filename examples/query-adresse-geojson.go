// +build ignore

package main

import (
	"fmt"
	"github.com/klauspost/dawa"
)

func main() {
	// Ask for an address with Vejnavn = Rødkildevej and Husnr = 46
	query := dawa.NewAdresseQuery().Vejnavn("Rødkildevej").Husnr("46", "44", "42")
	fmt.Println("Url:" + query.URL())

	geo, err := query.GeoJSON()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result Type:%s\n", geo.Type)
	for _, feat := range geo.Features {
		fmt.Printf("\tFeature Type:%s\n", feat.Type)
		if feat.Crs != nil {
			fmt.Printf("\tFeature CRS:%s\n", *feat.Crs)
		}
		geometry, err := feat.GetGeometry()
		if err != nil {
			fmt.Printf("\tGeometry decoding error %s\n", err)
		} else {
			fmt.Printf("\tGeometry type: %s\n", geometry.GetType())
			fmt.Printf("\tGeometry internal type: %T\n", geometry)
			fmt.Printf("\tGeometry:\n\t\t %+v\n", geometry)
		}
		fmt.Printf("\tProperties:\n")
		for key, prop := range feat.Properties {
			fmt.Printf("\t\t%s:%v\n", key, prop)
		}
	}
}
