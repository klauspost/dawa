package dawa

import (
	"bufio"
	"github.com/ugorji/go/codec"
	"io"
)

type Postnummer struct {
	Href     string    `json:"href"`     // Postnummerets unikke URL.
	Kommuner []Kommune `json:"kommuner"` // De kommuner hvis areal overlapper postnumeret areal.
	Navn     string    `json:"navn"`     // Det navn der er knyttet til postnummeret, typisk byens eller bydelens navn. Repræsenteret ved indtil 20 tegn. Eksempel: ”København NV”.
	Nr       string    `json:"nr"`       // Unik identifikation af det postnummeret. Postnumre fastsættes af Post Danmark. Repræsenteret ved fire cifre. Eksempel: ”2400” for ”København NV”.
	// Never set to anything but null
	Stormodtageradresser []AdgangsAdresseRef `json:"stormodtageradresser"` // Hvis postnummeret er et stormodtagerpostnummer rummer feltet adresserne på stormodtageren.
}

// PostnummerIter is an Iterator that enable you to get individual entries.
type PostnummerIter struct {
	a   chan Postnummer
	err error
}

// Next will return addresses.
// It will return an error if that has been encountered.
// When there are not more entries nil, io.EOF will be returned.
func (a *PostnummerIter) Next() (*Postnummer, error) {
	v, ok := <-a.a
	if ok {
		return &v, nil
	}
	return nil, a.err
}

// ImportPostnumreJSON will import "postnumre" from a JSON input, supplied to the reader.
// An iterator will be returned that return all items.
func ImportPostnumreJSON(in io.Reader) (*PostnummerIter, error) {
	var h codec.JsonHandle
	h.DecodeOptions.ErrorIfNoField = JSONStrictFieldCheck
	// use a buffered reader for efficiency
	if _, ok := in.(io.ByteScanner); !ok {
		in = bufio.NewReader(in)
	}
	ret := &PostnummerIter{a: make(chan Postnummer, 100)}
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
