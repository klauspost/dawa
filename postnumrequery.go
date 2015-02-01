package dawa

import (
	"fmt"
	"io"
)

// PostnrQuery is a new query for 'postnummer' objects for searching DAWA.
// Use NewPostnrQuery() or NewPostnrComplete() to get an initialized object.
// Example:
//			// Search for "Rødkildevej 46"
//			item, err := dawa.NewPostnrQuery().Navn("Rødby").First()
//
//			// If err is nil, we go a result
//			if err == nil {
//				fmt.Printf("Got item:%+v\n", item)
//			}
type PostnrQuery struct {
	queryGeoJSON
}

// NewPostnummerQuery returns a new query for 'postnummer' objects for searching DAWA.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func NewPostnrQuery() *PostnrQuery {
	return &PostnrQuery{queryGeoJSON: queryGeoJSON{query: query{host: DefaultHost, path: "/postnumre"}}}
}

// NewPostnrCompleteQuery returns a new autocomplete query for 'postnummer' objects for searching DAWA.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func NewPostnrComplete() *PostnrQuery {
	return &PostnrQuery{queryGeoJSON: queryGeoJSON{query: query{host: DefaultHost, path: "/postnumre/autocomplete"}}}
}

// GetPostnr will return a single Postnummer with the specified ID.
// Will return (nil, io.EOF) if there is no results.
func GetPostnr(id string) (*Postnummer, error) {
	return NewPostnrQuery().Nr(id).First()
}

// Iter will return an iterator that allows you to read the results
// one by one.
//
// An example:
//			iter, err := dawa.NewPostnrQuery().Vejnavn("Rødkildevej").Husnr("46").Iter()
//			if err != nil {
// 				panic(err)
// 			}
//
//			for {
//				a, err := iter.Next()
//				if err == io.EOF {
// 					iter.Close()
//					break  // we are finished
//				}
//				if err != nil {
//					panic(err)
//				}
// 				fmt.Printf("%+v\n", a)
//			}
//		}
func (q PostnrQuery) Iter() (*PostnummerIter, error) {
	resp, err := q.NoFormat().Request()
	if err != nil {
		return nil, err
	}

	iter, err := ImportPostnumreJSON(resp)
	if err != nil {
		return nil, err
	}
	iter.AddCloser(resp)
	return iter, nil
}

// All returns all results as an array.
func (q PostnrQuery) All() ([]Postnummer, error) {
	resp, err := q.NoFormat().Request()
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	ret := make([]Postnummer, 0)
	iter, err := ImportPostnumreJSON(resp)
	if err != nil {
		return nil, err
	}

	for {
		a, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if a != nil {
			ret = append(ret, *a)
		}
	}
	return ret, nil
}

// First will return the first result from a query.
// Note the entire query is executed, so only use this if you expect a few results.
//
// Will return (nil, io.EOF) if there is no results.
func (q PostnrQuery) First() (*Postnummer, error) {
	resp, err := q.NoFormat().Request()
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	iter, err := ImportPostnumreJSON(resp)
	if err != nil {
		return nil, err
	}

	a, err := iter.Next()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Nr will add a parameter for 'nr' to the PostnrQuery.
//
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Nr(s ...string) *PostnrQuery {
	q.add(&textQuery{Name: "nr", Values: s, Multi: true, Null: false})
	return q
}

// Navn will add a parameter for 'navn' to the PostnrQuery.
//
// Postnummernavn (Flerværdisøgning mulig).
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Navn(s ...string) *PostnrQuery {
	q.add(&textQuery{Name: "navn", Values: s, Multi: true, Null: false})
	return q
}

// Kommunekode will add a parameter for 'kommunekode' to the PostnrQuery.
//
// Kommunekode. 4 cifre. Eksempel: 0101 for Københavns kommune. (Flerværdisøgning mulig).
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Kommunekode(s ...string) *PostnrQuery {
	q.add(&textQuery{Name: "kommunekode", Values: s, Multi: true, Null: false})
	return q
}

// Q will add a parameter for 'q' to the PostnrQuery.
//
// Søgetekst. Der søges i postnummernavnet.
// Alle ord i søgeteksten skal matche postnummernavnet. Wildcard * er tilladt i slutningen af hvert ord.
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Q(s string) *PostnrQuery {
	q.add(&textQuery{Name: "q", Values: []string{s}, Multi: false, Null: true})
	return q
}

// Stormodtagere will add a parameter for 'stormodtagere' to the PostnrQuery.
//
// Hvis denne parameter er sat til 'true', vil stormodtager-postnumre medtages i resultatet. Default er false.
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Stormodtagere(b bool) *PostnrQuery {
	q.add(&textQuery{Name: "stormodtagere", Values: []string{fmt.Sprintf("%v", b)}, Multi: false, Null: false})
	return q
}

// NoFormat will disable extra whitespace. Always enabled when querying
func (q *PostnrQuery) NoFormat() *PostnrQuery {
	q.add(&textQuery{Name: "noformat", Multi: false, Null: true})
	return q
}
