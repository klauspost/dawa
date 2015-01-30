package dawa

import (
	"bufio"

	"github.com/ugorji/go/codec"

	"io"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
)

// ListQuery returns query item for searching DAWA for specific list types.
// Use dawa.NewListQuery(type string, autocomplete bool) to create a new query.
//
// See 'examples/query-list.go' for a usage example.
//
// See documentation at http://dawa.aws.dk/listerdok
type ListQuery struct {
	queryGeoJSON
	listType string
}

// ListQuery returns query item for searching DAWA for specific list types.
//
// Supported list types are "regioner","sogne","retskredse","politikredse","opstillingskredse","valglandsdele","ejerlav".
// Use the corresponding iterator function, for instance i.NextRegion() to get typed results.
//
// See 'examples/query-list.go' for a usage example.
//
// See documentation at http://dawa.aws.dk/listerdok
func NewListQuery(listType string, autoComplete bool) *ListQuery {
	path := "/" + listType
	if autoComplete {
		path += "/autocomplete"
	}
	q := &ListQuery{listType: listType, queryGeoJSON: queryGeoJSON{query: query{host: DefaultHost, path: path}}}
	return q
}

// Q will add a parameter for 'q' to the ListQuery.
//
// Søgetekst. Der søges i kode og navn. Alle ord i søgeteksten skal matche. Wildcard * er tilladt i slutningen af hvert ord.
//
// See http://dawa.aws.dk/listerdok
func (q *ListQuery) Q(s string) *ListQuery {
	q.add(textQuery{Name: "q", Values: []string{s}, Multi: true, Null: false})
	return q
}

// Kode will add a parameter for 'kode' to the ListQuery.
//
// Kode for det der søges.
func (q *ListQuery) Kode(s ...string) *ListQuery {
	q.add(textQuery{Name: "kode", Values: s, Multi: true, Null: false})
	return q
}

// Navn will add a parameter for 'navn' to the ListQuery.
//
// Navn for det der søges.
func (q *ListQuery) Navn(s string) *ListQuery {
	q.add(textQuery{Name: "navn", Values: []string{s}, Multi: true, Null: false})
	return q
}

// NoFormat will disable extra whitespace. Always enabled when querying
func (q *ListQuery) NoFormat() *ListQuery {
	q.add(textQuery{Name: "noformat", Multi: false, Null: true})
	return q
}

// ListIter is an Iterator that enable you to get individual entries.
type ListIter struct {
	closer
	a     reflect.Value // Channel
	eType reflect.Type  // Type of the element
	err   error
}

func makeChannel(t reflect.Type, chanDir reflect.ChanDir, buffer int) reflect.Value {
	ctype := reflect.ChanOf(chanDir, t)
	return reflect.MakeChan(ctype, buffer)
}

// Iter creates a list iterator that will allow you to get the items one by one.
//
func (q ListQuery) Iter() (*ListIter, error) {
	resp, err := q.NoFormat().Request()
	if err != nil {
		return nil, err
	}

	typ := q.Type()
	if typ == nil {
		return nil, fmt.Errorf("Unknown list type: %s", q.listType)
	}
	var h codec.JsonHandle
	h.DecodeOptions.ErrorIfNoField = JSONStrictFieldCheck

	// use a buffered reader for efficiency
	in := bufio.NewReader(resp)
	ret := &ListIter{}
	ret.eType = reflect.TypeOf(typ)
	// We create a channel with the expected type
	ret.a = makeChannel(ret.eType, reflect.BothDir, 100)
	go func() {
		defer ret.a.Close()
		var dec *codec.Decoder = codec.NewDecoder(in, &h)
		channel := ret.a.Interface()
		ret.err = dec.Decode(&channel)
		if ret.err == nil {
			ret.err = io.EOF
		}
	}()

	if err != nil {
		return nil, err
	}
	ret.AddCloser(resp)
	return ret, nil
}

// Returns a writeable
func (q ListQuery) Type() interface{} {
	switch q.listType {
	case "regioner":
		return &Region{}
	case "kommuner":
		return &Kommune{}
	case "sogne":
		return &Sogn{}
	case "retskredse":
		return &Retskreds{}
	case "politikredse":
		return &Politikreds{}
	case "opstillingskredse":
		return &Opstillingskreds{}
	case "valglandsdele":
		return &Valglandsdel{}
	case "ejerlav":
		return &Ejerlav{}
	}
	return nil
}

// Next will return the next item untyped.
// It will return an error if that has been encountered.
// When there are not more entries nil, io.EOF will be returned.
func (a *ListIter) Next() (interface{}, error) {
	v, ok := a.a.Recv()
	if ok {
		return v.Interface(), nil
	}
	return nil, a.err
}

// NextKommune will return the next item.
// The query must be built using the corresponding type. See NewListQuery() function.
func (a *ListIter) NextKommune() (*Kommune, error) {
	if !a.eType.ConvertibleTo(reflect.TypeOf(&Kommune{})) {
		return nil, fmt.Errorf("Wrong type requested from iterator. Expected %s", a.eType.String())
	}
	item, err := a.Next()
	if err != nil {
		return nil, a.err
	}
	return item.(*Kommune), nil
}

// NextRegion will return the next item.
// The query must be built using the corresponding type. See NewListQuery() function.
func (a *ListIter) NextRegion() (*Region, error) {
	if !a.eType.ConvertibleTo(reflect.TypeOf(&Region{})) {
		return nil, fmt.Errorf("Wrong type requested from iterator. Expected %s", a.eType.String())
	}
	item, err := a.Next()
	if err != nil {
		return nil, a.err
	}
	return item.(*Region), nil
}

// NextSogn will return the next item.
// The query must be built using the corresponding type. See NewListQuery() function.
func (a *ListIter) NextSogn() (*Sogn, error) {
	if !a.eType.ConvertibleTo(reflect.TypeOf(&Sogn{})) {
		return nil, fmt.Errorf("Wrong type requested from iterator. Expected %s", a.eType.String())
	}
	item, err := a.Next()
	if err != nil {
		return nil, a.err
	}
	return item.(*Sogn), nil
}

// NextRetskreds will return the next item.
// The query must be built using the corresponding type. See NewListQuery() function.
func (a *ListIter) NextRetskreds() (*Retskreds, error) {
	if !a.eType.ConvertibleTo(reflect.TypeOf(&Retskreds{})) {
		return nil, fmt.Errorf("Wrong type requested from iterator. Expected %s", a.eType.String())
	}
	item, err := a.Next()
	if err != nil {
		return nil, a.err
	}
	return item.(*Retskreds), nil
}

// NextPolitikreds will return the next item.
// The query must be built using the corresponding type. See NewListQuery() function.
func (a *ListIter) NextPolitikreds() (*Politikreds, error) {
	if !a.eType.ConvertibleTo(reflect.TypeOf(&Politikreds{})) {
		return nil, fmt.Errorf("Wrong type requested from iterator. Expected %s", a.eType.String())
	}
	item, err := a.Next()
	if err != nil {
		return nil, a.err
	}
	return item.(*Politikreds), nil
}

// NextOpstillingskreds will return the next item.
// The query must be built using the corresponding type. See NewListQuery() function.
func (a *ListIter) NextOpstillingskreds() (*Opstillingskreds, error) {
	if !a.eType.ConvertibleTo(reflect.TypeOf(&Opstillingskreds{})) {
		return nil, fmt.Errorf("Wrong type requested from iterator. Expected %s", a.eType.String())
	}
	item, err := a.Next()
	if err != nil {
		return nil, a.err
	}
	return item.(*Opstillingskreds), nil
}

// NextValglandsdel will return the next item.
// The query must be built using the corresponding type. See NewListQuery() function.
func (a *ListIter) NextValglandsdel() (*Valglandsdel, error) {
	if !a.eType.ConvertibleTo(reflect.TypeOf(&Valglandsdel{})) {
		return nil, fmt.Errorf("Wrong type requested from iterator. Expected %s", a.eType.String())
	}
	item, err := a.Next()
	if err != nil {
		return nil, a.err
	}
	return item.(*Valglandsdel), nil
}

// NextEjerlav will return the next item.
// The query must be built using the corresponding type. See NewListQuery() function.
func (a *ListIter) NextEjerlav() (*Ejerlav, error) {
	if !a.eType.ConvertibleTo(reflect.TypeOf(&Ejerlav{})) {
		return nil, fmt.Errorf("Wrong type requested from iterator. Expected %s", a.eType.String())
	}
	item, err := a.Next()
	if err != nil {
		return nil, a.err
	}
	return item.(*Ejerlav), nil
}

// NewReverseQuery will create a reverse location to item lookup. Parameters are:
//
//	* listType: See NewListQuery() for valid options.
// 	* x: X koordinat. (Hvis ETRS89/UTM32 anvendes angives øst-værdien.) Hvis WGS84/geografisk anvendex angives bredde-værdien.
// 	* y: Y koordinat. (Hvis ETRS89/UTM32 anvendes angives nord-værdien.) Hvis WGS84/geografisk anvendex angives længde-værdien.
//	* srid: Angiver SRID for det koordinatsystem, som geospatiale parametre er angivet i. Default er 4326 (WGS84). Leave this empty for default value
//
//  See examples/query-list-reverse.go for usage example
//
// An iterator will be returned, but it will only contain zero or one values.
func NewReverseQuery(listType string, x, y float64, srid string) (*ListIter, error) {
	path := "/" + listType + "/reverse"
	q := &ListQuery{listType: listType, queryGeoJSON: queryGeoJSON{query: query{host: DefaultHost, path: path}}}
	typ := q.Type()
	if typ == nil {
		return nil, fmt.Errorf("unknown list type '%s'", listType)
	}
	q.add(textQuery{Name: "x", Values: []string{strconv.FormatFloat(x, 'f', -1, 64)}, Multi: false, Null: false})
	q.add(textQuery{Name: "y", Values: []string{strconv.FormatFloat(y, 'f', -1, 64)}, Multi: false, Null: false})
	if srid != "" {
		q.add(textQuery{Name: "srid", Values: []string{srid}, Multi: false, Null: false})
	}
	// Execute request
	resp, err := q.NoFormat().Request()
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	// We create the iterator and fill it with data.
	ret := &ListIter{}
	ret.eType = reflect.TypeOf(typ)
	ret.a = makeChannel(ret.eType, reflect.BothDir, 1)

	all, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	// Decode to typ
	err = json.Unmarshal(all, typ)
	if err != nil {
		return nil, err
	}
	ret.a.Send(reflect.ValueOf(typ))
	ret.a.Close()
	ret.err = io.EOF
	return ret, nil
}
