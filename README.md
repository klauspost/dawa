# dawa
Golang implementation of DAWA AWS Suite 4 (Danish Address Info)

The dawa package can be used to de-serialize structures received from "Danmarks Adressers Web API (DAWA)" (Addresses of Denmark Web API).

* This package allows to de-serialize JSON responses from the web api into typed structs.
* The package also allows importing JSON or CSV downloads from the official web page.

See the /examples folder for more information.

Package home: https://github.com/klauspost/dawa

Information abou the format and download/API options, see http://dawa.aws.dk/

# Installation 

```go get github.com/klauspost/dawa/...```

This will also install the only dependecy "go-codec": https://github.com/ugorji/go which is used to decode JSON streams more efficiently than the standard golang libraries.

To use the library in your own code, simply add ```"github.com/klauspost/dawa"``` to your imports.

# Documentation
[![GoDoc][1]][2] [![Build Status][3]][4]
[1]: https://godoc.org/github.com/klauspost/dawa?status.svg
[2]: https://godoc.org/github.com/klauspost/dawa
[3]: https://travis-ci.org/klauspost/dawa.svg
[4]: https://travis-ci.org/klauspost/dawa

# Usage

Note that these examples ignore error values. You should of course never do that. For complete examples with error checking, see the /examples fodler.

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

# License

This code is published under an MIT license. See LICENSE file for more information.