// The dawa package can be used to de-serialize structures received from "Danmarks Adressers Web API (DAWA)" (Addresses of Denmark Web API).
//
// This package allows to de-serialize JSON responses from the web api into typed structs.
// The package also allows importing JSON or CSV downloads from the official web page.
// See the /examples folder for more information.
//
// Package home: https://github.com/klauspost/dawa
//
// Information abou the format and download/API options, see http://dawa.aws.dk/
//
// Description text in Danish:
//
// Danmarks Adressers Web API (DAWA) udstiller data og funktionalitet vedrørende Danmarks adresser, adgangsadresser, vejnavne samt postnumre.
// DAWA anvendes til etablering af adressefunktionalitet i it-systemer. Målgruppen for nærværende website er udviklere, som ønsker at indbygge adressefunktionalitet i deres it-systemer.
package dawa

import (
	"io"
)

// modify JSONStrictFieldCheck to return an error on unknown fields on JSON import.
// If true, return an error if a map in the stream has a key which does not map to any field; else read and discard the key and value in the stream and proceed to the next.
var JSONStrictFieldCheck = false

type closer struct {
	c []io.Closer
}

// Call this when you are finished using the object
func (c *closer) Close() error {
	for _, cl := range c.c {
		err := cl.Close()
		if err != nil {
			return err
		}
	}
	c.c = nil
	return nil
}

func (c *closer) AddCloser(a io.Closer) {
	c.c = append(c.c, a)
}
