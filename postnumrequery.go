package dawa

import (
	"fmt"
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
	query
}

// NewAdgangsAdresseQuery returns a new query for 'postnummer' objects for searching DAWA.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func NewPostnrQuery() *PostnrQuery {
	return &PostnrQuery{query: query{host: DefaultHost, path: "/postnumre"}}
}

// NewPostnrCompleteQuery returns a new autocomplete query for 'postnummer' objects for searching DAWA.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func NewPostnrComplete() *PostnrQuery {
	return &PostnrQuery{query: query{host: DefaultHost, path: "/postnumre/autocomplete"}}
}

// Nr will add a parameter for 'nr' to the PostnrQuery.
//
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Nr(s ...string) *PostnrQuery {
	q.add(textQuery{Name: "nr", Values: s, Multi: true, Null: false})
	return q
}

// Navn will add a parameter for 'navn' to the PostnrQuery.
//
// Postnummernavn (Flerværdisøgning mulig).
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Navn(s ...string) *PostnrQuery {
	q.add(textQuery{Name: "navn", Values: s, Multi: true, Null: false})
	return q
}

// Kommunekode will add a parameter for 'kommunekode' to the PostnrQuery.
//
// Kommunekode. 4 cifre. Eksempel: 0101 for Københavns kommune. (Flerværdisøgning mulig).
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Kommunekode(s ...string) *PostnrQuery {
	q.add(textQuery{Name: "kommunekode", Values: s, Multi: true, Null: false})
	return q
}

// Q will add a parameter for 'q' to the PostnrQuery.
//
// Søgetekst. Der søges i postnummernavnet.
// Alle ord i søgeteksten skal matche postnummernavnet. Wildcard * er tilladt i slutningen af hvert ord.
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Q(s string) *PostnrQuery {
	q.add(textQuery{Name: "q", Values: []string{s}, Multi: false, Null: true})
	return q
}

// Stormodtagere will add a parameter for 'stormodtagere' to the PostnrQuery.
//
// Hvis denne parameter er sat til 'true', vil stormodtager-postnumre medtages i resultatet. Default er false.
//
// See http://dawa.aws.dk/postnummerdok#postnummersoegning
func (q *PostnrQuery) Stormodtagere(b bool) *PostnrQuery {
	q.add(textQuery{Name: "stormodtagere", Values: []string{fmt.Sprintf("%v", b)}, Multi: false, Null: false})
	return q
}
