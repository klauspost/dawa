# dawa
Golang implementation of DAWA AWS Suite 4 (Danish Address Info)

The dawa package can be used to de-serialize structures received from "Danmarks Adressers Web API (DAWA)" (Addresses of Denmark Web API).

* This package allows to de-serialize JSON responses from the web api into typed structs.
* The package also allows importing JSON or CSV downloads from the official web page.
* The package helps you build queries against the official web api.

See the /examples folder for more information.

Package home: https://github.com/klauspost/dawa

Information about the format and download/API options, see http://dawa.aws.dk/

# Installation 

```go get github.com/klauspost/dawa/...```

This will also install the only dependecy "go-codec": https://github.com/ugorji/go which is used to decode JSON streams more efficiently than the standard golang libraries.

To use the library in your own code, simply add ```import "github.com/klauspost/dawa"``` to your imports.

# Documentation
[![GoDoc][1]][2] [![Build Status][3]][4]
[1]: https://godoc.org/github.com/klauspost/dawa?status.svg
[2]: https://godoc.org/github.com/klauspost/dawa
[3]: https://travis-ci.org/klauspost/dawa.svg
[4]: https://travis-ci.org/klauspost/dawa

# Usage

Note that these examples ignore error values. You should of course never do that. For complete examples with error checking, see the /examples folder.

The library supplies all known data structures, which means that it can decode responses from the Web API:

```Go
	// Get a single address from Web API
	resp, _ := http.Get("http://dawa.aws.dk/adresser/0a3f509d-07be-32b8-e044-0003ba298018")

	// Read the response into r
	r, _ := ioutil.ReadAll(resp.Body)

	// Create destination struct and decode the response
	addresse := dawa.Adresse{}
	_ = json.Unmarshal(r, &add)

```

The library also supplies streaming decoders for an array of responses:

```Go
	// Get an array of addresses
	resp, _ := http.Get("http://dawa.aws.dk/adgangsadresser?husnr=14&postnr=9000")

	// Send the response stream to a decoder
	iter, _ := dawa.ImportAdgangsAdresserJSON(resp.Body)

	// Read responses one by one
	for {
		address, err := iter.Next()
		if err == io.EOF {
			break
		}
		// 'address' now contains an 'Adgangsadresse'
	}

```
If you would like to have help in building queries for the web API, there are query builders for each data type. For instance to execute the query above, you can do it like this:
```Go
	// Get the matching adgangsadresser
	iter, _ := dawa.NewAdgangsAdresseQuery().Husnr("14").Postnr("9000").Iter()

	// Read responses one by one
	for {
		address, err := iter.Next()
		if err == io.EOF {
			break
		}
		// 'address' now contains an 'Adgangsadresse'
	}

```
To get all results in a single array is a simple oneliner:
```Go
	// Get the matching adgangsadresser
	results, _ := dawa.NewAdgangsAdresseQuery().Husnr("14").Postnr("9000").All()
```


The same API is used to decode content from files downloaded from http://download.aws.dk/

```Go
	// Open file
	file, _ := os.Open("adgangsadresser.json")

	// Send the file stream to the decoder
	iter, _ := dawa.ImportAdgangsAdresserJSON(file)

	// Read responses one by one
	for {
		address, err := iter.Next()
		if err == io.EOF {
			break
		}
		// 'address' now contains an 'Adgangsadresse'
	}

```

For 'Adgangsadresser' and 'Adresser' also gives the option to decode from CSV files instead of JSON. Note however, that not all information is present in the CSV files, so not all fields will be filled.

The API is similar to the JSON API:

```Go
	// Open file
	file, _ := os.Open("adgangsadresser.csv")

	// Send the file stream to the decoder
	iter, _ := dawa.ImportAdgangsAdresserCSV(file)

	// Read responses one by one
	for {
		address, err := iter.Next()
		if err == io.EOF {
			break
		}
		// 'address' now contains an 'Adgangsadresse'
	}

```

# Queries

There is a search API to assist you in building queries for the DAWA Web API.

Currently it supports queries for "adresser" and "adgangsadresser".

You can use a ```dawa.NewAdresseQuery()``` to start a new query. Parameters can be appended to the query, by simply calling the matching functions. For example to get Danmarksgade in Aalborg, use a query like this ```query := dawa.NewAdresseQuery().Vejnavn("Danmarksgade").Postnr("9000")```.

If you want the URL for a query, you can call the .URL() function, but you can also request all results by calling .All(), get an iterator for the results with .Iter(), or just get the first result with .First()

To send multiple query values of the same type, you should specify them in the same function call, so if you are looking for "postnr" with values 6400 and 6500 you can use the query ```q := dawa.NewAdresseQuery().Postnr("6400", "6500")```. For values that support this, you can signify a query for an empty value, by simply not sending any parameters, for example ```q := dawa.NewAdresseQuery().Etage()``` will search for values where 'etage' is unset.

# Query Examples
Get a single item:
```Go
// Search for "Rødkildevej 46"
item, err := dawa.NewAdgangsAdresseQuery().Vejnavn("Rødkildevej").Husnr("46").First()

// If io.EOF, there were no results.
if err == io.EOF {
	fmt.Printf("No results")
} else if err != nil {
	// There was another error
	fmt.Printf("Error:%s", err.Error())
} else {
	// We got a result
	fmt.Printf("Got item:%+v\n", item)
}
```

Query where a parameter can have multiple values.
```Go
// Search for "Rødkildevej 44,45 and 46"
item, err := dawa.NewAdgangsAdresseQuery().Vejnavn("Rødkildevej").Husnr("44", "45", "46").All()

fmt.Printf("Got item:%+v\n", item)
```


Get all results from a query:
```Go
	iter, err := dawa.NewAdresseQuery().Vejnavn("Rødkildevej").Husnr("46").Iter()
	if err != nil {
		panic(err)
	}

	for {
		a, err := iter.Next()
		if err == io.EOF {
			iter.Close()
			break  // we are finished
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", a)
	}
```

You can get the results as GeoJSON by using the GeoJSON function on any query:
```Go
geoj, err := dawa.NewAdgangsAdresseQuery().Vejnavn("Rødkildevej").Husnr("44").GeoJSON()
fmt.Printf("Got location:%+v\n", geoj)
```
See ```examples/query-adresse-geojson.go``` on how to parse the result.


# License

This code is published under an MIT license. See LICENSE file for more information.
