package dawa

import (
	"io"
	"strconv"
)

// AdgangsAdresseQuery is a new query for 'adgangsadresser' objects for searching DAWA.
// Use NewAdgangsAdresseQuery() or NewAdgangsAdresseComplete() to get an initialized object.
// Example:
//			// Search for "Rødkildevej 46"
//			item, err := dawa.NewAdgangsAdresseQuery().Vejnavn("Rødkildevej").Husnr("46").First()
//
//			// If err is nil, we go a result
//			if err == nil {
//				fmt.Printf("Got item:%+v\n", item)
//			}
type AdgangsAdresseQuery struct {
	queryGeoJSON
}

// NewAdgangsAdresseQuery returns a new query for 'adgangsadresser objects for searching DAWA.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func NewAdgangsAdresseQuery() *AdgangsAdresseQuery {
	return &AdgangsAdresseQuery{queryGeoJSON: queryGeoJSON{query: query{host: DefaultHost, path: "/adgangsadresser"}}}
}

// NewAdgangsAdresseQuery returns a new query for 'adgangsadresser' objects for searching DAWA with autocomplete.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adresseautocomplete
func NewAdgangsAdresseComplete() *AdgangsAdresseQuery {
	return &AdgangsAdresseQuery{queryGeoJSON: queryGeoJSON{query: query{host: DefaultHost, path: "/autocomplete"}}}
}

// GetAAID will return a single AdgangsAdresse with the specified ID.
// Will return (nil, io.EOF) if there is no results.
func GetAAID(id string) (*AdgangsAdresse, error) {
	return NewAdgangsAdresseQuery().ID(id).First()
}

// Iter will return an iterator that allows you to read the results
// one by one.
//
// An example:
//			iter, err := dawa.NewAdgangsAdresseQuery().Vejnavn("Rødkildevej").Husnr("46").Iter()
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
func (q AdgangsAdresseQuery) Iter() (*AdgangsAdresseIter, error) {
	resp, err := q.NoFormat().Request()
	if err != nil {
		return nil, err
	}

	iter, err := ImportAdgangsAdresserJSON(resp)
	if err != nil {
		return nil, err
	}
	iter.AddCloser(resp)
	return iter, nil
}

// All returns all results as an array.
func (q AdgangsAdresseQuery) All() ([]AdgangsAdresse, error) {
	resp, err := q.NoFormat().Request()
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	ret := make([]AdgangsAdresse, 0)
	iter, err := ImportAdgangsAdresserJSON(resp)
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
func (q AdgangsAdresseQuery) First() (*AdgangsAdresse, error) {
	resp, err := q.NoFormat().Request()
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	iter, err := ImportAdgangsAdresserJSON(resp)
	if err != nil {
		return nil, err
	}

	a, err := iter.Next()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Q will add a parameter for 'q' to the AdgangsAdresseQuery.
//
// Søgetekst. Der søges i vejnavn, husnr, etage, dør, supplerende bynavn, postnr og postnummerets navn.
// Alle ord i søgeteksten skal matche adressebetegnelsen. Wildcard * er tilladt i slutningen af hvert ord.
// Der skelnes ikke mellem store og små bogstaver.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Q(s string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "q", Values: []string{s}})
	return q
}

// ID will add a parameter for 'id' to the AdgangsAdresseQuery.
//
// AdgangsAdressens unikke id, f.eks. 0a3f5095-45ec-32b8-e044-0003ba298018. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) ID(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "id", Values: s, Multi: true})
	return q
}

// Kvh will add a parameter for 'kvh' to the AdgangsAdresseQuery.
//
// KVH-nøgle. 12 tegn bestående af 4 cifre der repræsenterer kommunekode,
/// 4 cifre der repræsenterer vejkode efterfulgt af 4 tegn der repræsenter husnr
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Kvh(s string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "kvh", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Status will add a parameter for 'status' to the AdgangsAdresseQuery.
//
// AdgangsAdressens status, som modtaget fra BBR. "1" angiver en endelig adresse og "3" angiver en foreløbig adresse".
// AdgangsAdresser med status "2" eller "4" er ikke med i DAWA.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Status(i int) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "status", Values: []string{strconv.Itoa(i)}, Multi: false, Null: false})
	return q
}

// Vejkode will add a parameter for 'vejkode' to the AdgangsAdresseQuery.
//
// Vejkoden. 4 cifre. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Vejkode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "vejkode", Values: s, Multi: true, Null: false})
	return q
}

// Vejnavn will add a parameter for 'vejnavn' to the AdgangsAdresseQuery.
//
// Vejnavn. Der skelnes mellem store og små bogstaver.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Vejnavn(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "vejnavn", Values: s, Multi: true, Null: true})
	return q
}

// Husnr will add a parameter for 'husnr' to the AdgangsAdresseQuery.
//
// Husnummer. Max 4 cifre eventuelt med et efterfølgende bogstav. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Husnr(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "husnr", Values: s, Multi: true, Null: false})
	return q
}

// SupplerendeBynavn will add a parameter for 'supplerendebynavn' to the AdgangsAdresseQuery.
//
// Det supplerende bynavn. (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) SupplerendeBynavn(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "supplerendebynavn", Values: s, Multi: true, Null: true})
	return q
}

// Postnr will add a parameter for 'postnr' to the AdgangsAdresseQuery.
//
// Postnummer. 4 cifre. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Postnr(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "postnr", Values: s, Multi: true, Null: false})
	return q
}

// Kommunekode will add a parameter for 'kommunekode' to the AdgangsAdresseQuery.
//
// Kommunekoden for den kommune som adressen skal ligge på. 4 cifre. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Kommunekode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "kommunekode", Values: s, Multi: true, Null: false})
	return q
}

// Ejerlavkode will add a parameter for 'ejerlavkode' to the AdgangsAdresseQuery.
//
// Koden på det matrikulære ejerlav som adressen skal ligge på. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Ejerlavkode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "ejerlavkode", Values: s, Multi: true, Null: false})
	return q
}

// Zonekode will add a parameter for 'zonekode' to the AdgangsAdresseQuery.
//
// Heltalskoden for den zone som adressen skal ligge i.
// Mulige værdier er 1 for byzone, 2 for sommerhusområde og 3 for landzone.
// (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Zonekode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "zonekode", Values: s, Multi: true, Null: false})
	return q
}

// Matrikelnr will add a parameter for 'matrikelnr' to the AdgangsAdresseQuery.
//
// Matrikelnummer. Unikt indenfor et ejerlav. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Matrikelnr(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "matrikelnr", Values: s, Multi: true, Null: false})
	return q
}

// Esrejendomsnr will add a parameter for 'esrejendomsnr' to the AdgangsAdresseQuery.
//
// ESR Ejendomsnummer. Indtil 7 cifre. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Esrejendomsnr(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "esrejendomsnr", Values: s, Multi: true, Null: false})
	return q
}

// Srid will add a parameter for 'srid' to the AdgangsAdresseQuery.
//
// Angiver SRID for det koordinatsystem, som geospatiale parametre er angivet i. Default er 4326 (WGS84).
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Srid(s string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "srid", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Polygon will add a parameter for 'polygon' to the AdgangsAdresseQuery.
//
// Find de adresser, som ligger indenfor det angivne polygon.
// Polygonet specificeres som et array af koordinater på samme måde som koordinaterne
// specificeres i GeoJSON's polygon.
// Bemærk at polygoner skal være lukkede, dvs. at første og sidste koordinat skal være identisk.
// Som koordinatsystem kan anvendes (ETRS89/UTM32 eller) WGS84/geografisk.
// Dette angives vha. srid parameteren, se ovenover.
// Eksempel: Polygon("[[[10.3,55.3],[10.4,55.3],[10.4,55.31],[10.4,55.31],[10.3,55.3]]]")
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Polygon(s string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "polygon", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Cirkel will add a parameter for 'cirkel' to the AdgangsAdresseQuery.
//
// Find de adresser, som ligger indenfor den cirkel angivet af koordinatet (x,y) og radius r.
// Som koordinatsystem kan anvendes (ETRS89/UTM32 eller) WGS84/geografisk.
// Radius angives i meter. Cirkel("{x},{y},{r}")
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Cirkel(s string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "cirkel", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Regionskode will add a parameter for 'regionskode' to the AdgangsAdresseQuery.
//
// Find de adresser som ligger indenfor regionen angivet ved regionkoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Regionskode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "regionskode", Values: s, Multi: true, Null: true})
	return q
}

// Sognekode will add a parameter for 'sognekode' to the AdgangsAdresseQuery.
//
// Find de adresser som ligger indenfor sognet angivet ved sognkoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Sognekode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "sognekode", Values: s, Multi: true, Null: true})
	return q
}

// Opstillingskredskode will add a parameter for 'opstillingskredskode' to the AdgangsAdresseQuery.
//
// Find de adresser som ligger indenfor opstillingskredsen angivet ved opstillingskredskoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Opstillingskredskode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "opstillingskredskode", Values: s, Multi: true, Null: true})
	return q
}

// Retskredskode will add a parameter for 'retskredskode' to the AdgangsAdresseQuery.
//
// Find de adresser som ligger indenfor retskredsen angivet ved retskredskoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Retskredskode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "retskredskode", Values: s, Multi: true, Null: true})
	return q
}

// Politikredskode will add a parameter for 'politikredskode' to the AdgangsAdresseQuery.
//
// Find de adresser som ligger indenfor politikredsen angivet ved politikredskoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Politikredskode(s ...string) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "politikredskode", Values: s, Multi: true, Null: true})
	return q
}

// Side will add a parameter for 'side' to the AdgangsAdresseQuery.
//
// Angiver hvilken siden som skal leveres. Se Paginering.
// http://dawa.aws.dk/generelt#paginering
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) Side(i int) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "side", Values: []string{strconv.Itoa(i)}, Multi: false, Null: true})
	return q
}

// PerSide will add a parameter for 'per_side' to the AdgangsAdresseQuery.
//
// Antal resultater per side. Se Paginering.
// http://dawa.aws.dk/generelt#paginering
//
// See documentation at http://dawa.aws.dk/adgangsadressedok#adressesoegning
func (q *AdgangsAdresseQuery) PerSide(i int) *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "per_side", Values: []string{strconv.Itoa(i)}, Multi: false, Null: true})
	return q
}

// NoFormat will disable extra whitespace. Always enabled when querying
func (q *AdgangsAdresseQuery) NoFormat() *AdgangsAdresseQuery {
	q.add(&textQuery{Name: "noformat", Multi: false, Null: true})
	return q
}
