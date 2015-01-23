package dawa

import (
	"bufio"
	"github.com/ugorji/go/codec"
	"io"
)

// Et supplerende bynavn – typisk landsbyens navn – eller andet lokalt stednavn
// der er fastsat af kommunen for at præcisere adressens beliggenhed indenfor postnummeret.
//
// Indgår som en del af den officielle adressebetegnelse.
type SupplBynavn struct {
	Navn      string          `json:"navn"`      // Det supplerende bynavn. Indtil 34 tegn. Eksempel: ”Sønderholm”.
	Href      string          `json:"href"`      // Det supplerende bynavns unikke URL
	Kommuner  []Kommune       `json:"kommuner"`  // Kommuner, som det supplerende bynavn er beliggende i.
	Postnumre []PostnummerRef `json:"postnumre"` // Postnumre, som det supplerende bynavn er beliggende i.
}

// SupplBynavnIter is an Iterator that enable you to get individual entries.
type SupplBynavnIter struct {
	a   chan SupplBynavn
	err error
}

// Next will return the next item in the array.
// It will return an error if that has been encountered.
// When there are not more entries nil, io.EOF will be returned.
func (a *SupplBynavnIter) Next() (*SupplBynavn, error) {
	v, ok := <-a.a
	if ok {
		return &v, nil
	}
	return nil, a.err
}

// ImportSupplBynavnJSON will import "supplerende bynavne" from a JSON input, supplied to the reader.
// An iterator will be returned that return all items.
func ImportSupplBynavnJSON(in io.Reader) (*SupplBynavnIter, error) {
	var h codec.JsonHandle
	h.DecodeOptions.ErrorIfNoField = JSONStrictFieldCheck
	// use a buffered reader for efficiency
	if _, ok := in.(io.ByteScanner); !ok {
		in = bufio.NewReader(in)
	}
	ret := &SupplBynavnIter{a: make(chan SupplBynavn, 100)}
	go func() {
		defer close(ret.a)
		var dec *codec.Decoder = codec.NewDecoder(in, &h)
		ret.err = dec.Decode(&ret.a)
		if ret.err == nil {
			ret.err = io.EOF
		}
	}()

	return ret, nil
}
