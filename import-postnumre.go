package dawa

import (
	"bufio"
	"github.com/ugorji/go/codec"
	"io"
)

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
