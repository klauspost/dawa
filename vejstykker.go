package dawa

import (
	"bufio"
	"github.com/ugorji/go/codec"
	"io"
)

// Et vejstykke er en vej, som er afgrænset af en kommune.
// Et vejstykke er identificeret ved en kommunekode og en vejkode og har desuden et navn.
// En vej som gennemløber mere end en kommune vil bestå af flere vejstykker.
// Det er p.t. ikke muligt at få information om hvilke vejstykker der er en del af den samme vej.
// Vejstykker er udstillet under /vejstykker
type Vejstykke struct {
	Adresseringsnavn string          `json:"adresseringsnavn"` //En evt. forkortet udgave af vejnavnet på højst 20 tegn, som bruges ved adressering på labels og rudekuverter og lign., hvor der ikke plads til det fulde vejnavn.
	Historik         Historik        `json:"historik"`         // Væsentlige tidspunkter for vejstykket
	Href             string          `json:"href"`             // Vejstykkets unikke URL.
	Kode             string          `json:"kode"`             // Identifikation af vejstykke. Er unikt indenfor den pågældende kommune. Repræsenteret ved fire cifre. Eksempel: I Københavns kommune er ”0004” lig ”Abel Cathrines Gade”.
	Kommune          KommuneRef      `json:"kommune"`          // Kommunen som vejstykket er beliggende i.
	Navn             string          `json:"navn"`             // Vejens navn som det er fastsat og registreret af kommunen. Repræsenteret ved indtil 40 tegn. Eksempel: ”Hvidkildevej”.
	Postnumre        []PostnummerRef `json:"postnumre"`        // Postnummrene som vejstykket er beliggende i.
}

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
