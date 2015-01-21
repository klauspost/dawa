package dawa

import (
	"bufio"
	"github.com/ugorji/go/codec"
	"io"
)

// VejstykkeIter is an Iterator that enable you to get individual entries.
type VejstykkeIter struct {
	a   chan Vejstykke
	err error
}

// Next will return addresses.
// It will return an error if that has been encountered.
// When there are not more entries nil, io.EOF will be returned.
func (a *VejstykkeIter) Next() (*Vejstykke, error) {
	v, ok := <-a.a
	if ok {
		return &v, nil
	}
	return nil, a.err
}

// ImportVejstykkerJSON will import "vejstykker" from a JSON input, supplied to the reader.
// An iterator will be returned that return all items.
func ImportVejstykkerJSON(in io.Reader) (*VejstykkeIter, error) {
	var h codec.JsonHandle
	h.DecodeOptions.ErrorIfNoField = JSONStrictFieldCheck
	// use a buffered reader for efficiency
	if _, ok := in.(io.ByteScanner); !ok {
		in = bufio.NewReader(in)
	}
	ret := &VejstykkeIter{a: make(chan Vejstykke, 100)}
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
