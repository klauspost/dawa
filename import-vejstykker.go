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

// ImportVejstykkerJSON will import "adresser" from a JSON input, supplied to the reader.
// An iterator will be returned that return all addresses.
func ImportVejstykkerJSON(in io.Reader) (*VejstykkeIter, error) {
	reader := bufio.NewReader(in)
	// Skip until after '['
	_, err := reader.ReadBytes('[')
	if err != nil {
		return nil, err
	}
	// Start decoder
	ret := &VejstykkeIter{a: make(chan Vejstykke, 100)}
	go func() {
		defer close(ret.a)
		var h codec.JsonHandle
		h.ErrorIfNoField = true
		for {
			var dec *codec.Decoder = codec.NewDecoder(reader, &h)
			a := Vejstykke{}
			if err := dec.Decode(&a); err != nil {
				ret.err = err
				return
			}
			ret.a <- a

			// Skip comma
			if b, err := readByteSkippingSpace(reader); err != nil {
				ret.err = err
				return
			} else {
				switch b {
				case ',':
					continue
				case ']':
					ret.err = io.EOF
					return
				default:
					panic("Invalid character in JSON data: " + string([]byte{b}))
				}
			}

		}
	}()
	return ret, nil
}
