package dawa

import (
	"io"
	"net/http"
	"strconv"
)

// AdresseQuery is a new query for 'adresse' objects for searching DAWA.
// Use NewAdresseQuery() or NewAdresseComplete() to get an initialized object.
// Example:
//			// Search for "Rødkildevej 46"
//			item, err := dawa.NewAdresseQuery().Vejnavn("Rødkildevej").Husnr("46").First()
//
//			// If err is nil, we go a result
//			if err == nil {
//				fmt.Printf("Got item:%+v\n", item)
//			}
type AdresseQuery struct {
	query
}

// NewAdresseQuery returns a new query for 'adresse objects for searching DAWA.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func NewAdresseQuery() *AdresseQuery {
	return &AdresseQuery{query: query{host: DefaultHost, path: "/adresser"}}
}

// NewAdresseQuery returns a new query for 'adresse' objects for searching DAWA with autocomplete.
//
// See documentation at http://dawa.aws.dk/adressedok#adresseautocomplete
func NewAdresseComplete() *AdresseQuery {
	return &AdresseQuery{query: query{host: DefaultHost, path: "/adresser/autocomplete"}}
}

// Iter will return an iterator that allows you to read the results
// one by one.
//
// An example:
//			iter, err := dawa.NewAdresseQuery().Vejnavn("Rødkildevej").Husnr("46").Iter()
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
func (q AdresseQuery) Iter() (*AdresseIter, error) {
	resp, err := http.Get(q.NoFormat().URL())
	if err != nil {
		return nil, err
	}

	iter, err := ImportAdresserJSON(resp.Body)
	if err != nil {
		return nil, err
	}
	iter.AddCloser(resp.Body)
	return iter, nil
}

// All returns all results as an array.
func (q AdresseQuery) All() ([]Adresse, error) {
	resp, err := http.Get(q.NoFormat().URL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := make([]Adresse, 0)
	iter, err := ImportAdresserJSON(resp.Body)
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
func (q AdresseQuery) First() (*Adresse, error) {
	resp, err := http.Get(q.NoFormat().URL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	iter, err := ImportAdresserJSON(resp.Body)
	if err != nil {
		return nil, err
	}

	a, err := iter.Next()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Q will add a parameter for 'q' to the AdresseQuery.
//
// Søgetekst. Der søges i vejnavn, husnr, etage, dør, supplerende bynavn, postnr og postnummerets navn.
// Alle ord i søgeteksten skal matche adressebetegnelsen. Wildcard * er tilladt i slutningen af hvert ord.
// Der skelnes ikke mellem store og små bogstaver.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Q(s string) *AdresseQuery {
	q.add(textQuery{Name: "q", Values: []string{s}})
	return q
}

// ID will add a parameter for 'id' to the AdresseQuery.
//
// Adressens unikke id, f.eks. 0a3f5095-45ec-32b8-e044-0003ba298018. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) ID(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "id", Values: s, Multi: true})
	return q
}

// Adgangsadresseid will add a parameter for 'adgangsadresseid' to the AdresseQuery.
//
// Id på den til adressen tilknyttede adgangsadresse. UUID. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) AdgangsadresseID(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "adgangsadresseid", Values: s, Multi: true, Null: false})
	return q
}

// Etage will add a parameter for 'etage' to the AdresseQuery.
//
// Etagebetegnelse. Hvis værdi angivet kan den antage følgende værdier: tal fra 1 til 99, st, kl, kl2 op til kl9.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Etage(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "etage", Values: s, Multi: true, Null: true})
	return q
}

// Dør will add a parameter for 'dør' to the AdresseQuery.
//
// Dørbetegnelse. Tal fra 1 til 9999, små og store bogstaver samt tegnene / og -.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Dør(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "dør", Values: s, Multi: true, Null: true})
	return q
}

// Kvhx will add a parameter for 'kvhx' to the AdresseQuery.
//
// KVHX-nøgle. 19 tegn bestående af 4 cifre der repræsenterer kommunekode,
// 4 cifre der repræsenterer vejkode, 4 tegn der repræsenter husnr,
// 3 tegn der repræsenterer etage og 4 tegn der repræsenter dør
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Kvhx(s string) *AdresseQuery {
	q.add(textQuery{Name: "kvhx", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Status will add a parameter for 'status' to the AdresseQuery.
//
// Adressens status, som modtaget fra BBR. "1" angiver en endelig adresse og "3" angiver en foreløbig adresse".
// Adresser med status "2" eller "4" er ikke med i DAWA.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Status(s string) *AdresseQuery {
	q.add(textQuery{Name: "status", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Vejkode will add a parameter for 'vejkode' to the AdresseQuery.
//
// Vejkoden. 4 cifre. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Vejkode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "vejkode", Values: s, Multi: true, Null: false})
	return q
}

// Vejnavn will add a parameter for 'vejnavn' to the AdresseQuery.
//
// Vejnavn. Der skelnes mellem store og små bogstaver.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Vejnavn(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "vejnavn", Values: s, Multi: true, Null: true})
	return q
}

// Husnr will add a parameter for 'husnr' to the AdresseQuery.
//
// Husnummer. Max 4 cifre eventuelt med et efterfølgende bogstav. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Husnr(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "husnr", Values: s, Multi: true, Null: false})
	return q
}

// SupplerendeBynavn will add a parameter for 'supplerendebynavn' to the AdresseQuery.
//
// Det supplerende bynavn. (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) SupplerendeBynavn(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "supplerendebynavn", Values: s, Multi: true, Null: true})
	return q
}

// Postnr will add a parameter for 'postnr' to the AdresseQuery.
//
// Postnummer. 4 cifre. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Postnr(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "postnr", Values: s, Multi: true, Null: false})
	return q
}

// Kommunekode will add a parameter for 'kommunekode' to the AdresseQuery.
//
// Kommunekoden for den kommune som adressen skal ligge på. 4 cifre. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Kommunekode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "kommunekode", Values: s, Multi: true, Null: false})
	return q
}

// Ejerlavkode will add a parameter for 'ejerlavkode' to the AdresseQuery.
//
// Koden på det matrikulære ejerlav som adressen skal ligge på. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Ejerlavkode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "ejerlavkode", Values: s, Multi: true, Null: false})
	return q
}

// Zonekode will add a parameter for 'zonekode' to the AdresseQuery.
//
// Heltalskoden for den zone som adressen skal ligge i.
// Mulige værdier er 1 for byzone, 2 for sommerhusområde og 3 for landzone.
// (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Zonekode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "zonekode", Values: s, Multi: true, Null: false})
	return q
}

// Matrikelnr will add a parameter for 'matrikelnr' to the AdresseQuery.
//
// Matrikelnummer. Unikt indenfor et ejerlav. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Matrikelnr(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "matrikelnr", Values: s, Multi: true, Null: false})
	return q
}

// Esrejendomsnr will add a parameter for 'esrejendomsnr' to the AdresseQuery.
//
// ESR Ejendomsnummer. Indtil 7 cifre. (Flerværdisøgning mulig).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Esrejendomsnr(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "esrejendomsnr", Values: s, Multi: true, Null: false})
	return q
}

// Srid will add a parameter for 'srid' to the AdresseQuery.
//
// Angiver SRID for det koordinatsystem, som geospatiale parametre er angivet i. Default er 4326 (WGS84).
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Srid(s string) *AdresseQuery {
	q.add(textQuery{Name: "srid", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Polygon will add a parameter for 'polygon' to the AdresseQuery.
//
// Find de adresser, som ligger indenfor det angivne polygon.
// Polygonet specificeres som et array af koordinater på samme måde som koordinaterne
// specificeres i GeoJSON's polygon.
// Bemærk at polygoner skal være lukkede, dvs. at første og sidste koordinat skal være identisk.
// Som koordinatsystem kan anvendes (ETRS89/UTM32 eller) WGS84/geografisk.
// Dette angives vha. srid parameteren, se ovenover.
// Eksempel: Polygon("[[[10.3,55.3],[10.4,55.3],[10.4,55.31],[10.4,55.31],[10.3,55.3]]]")
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Polygon(s string) *AdresseQuery {
	q.add(textQuery{Name: "polygon", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Cirkel will add a parameter for 'cirkel' to the AdresseQuery.
//
// Find de adresser, som ligger indenfor den cirkel angivet af koordinatet (x,y) og radius r.
// Som koordinatsystem kan anvendes (ETRS89/UTM32 eller) WGS84/geografisk.
// Radius angives i meter. Cirkel("{x},{y},{r}")
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Cirkel(s string) *AdresseQuery {
	q.add(textQuery{Name: "cirkel", Values: []string{s}, Multi: false, Null: false})
	return q
}

// Regionskode will add a parameter for 'regionskode' to the AdresseQuery.
//
// Find de adresser som ligger indenfor regionen angivet ved regionkoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Regionskode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "regionskode", Values: s, Multi: true, Null: true})
	return q
}

// Sognekode will add a parameter for 'sognekode' to the AdresseQuery.
//
// Find de adresser som ligger indenfor sognet angivet ved sognkoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Sognekode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "sognekode", Values: s, Multi: true, Null: true})
	return q
}

// Opstillingskredskode will add a parameter for 'opstillingskredskode' to the AdresseQuery.
//
// Find de adresser som ligger indenfor opstillingskredsen angivet ved opstillingskredskoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Opstillingskredskode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "opstillingskredskode", Values: s, Multi: true, Null: true})
	return q
}

// Retskredskode will add a parameter for 'retskredskode' to the AdresseQuery.
//
// Find de adresser som ligger indenfor retskredsen angivet ved retskredskoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Retskredskode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "retskredskode", Values: s, Multi: true, Null: true})
	return q
}

// Politikredskode will add a parameter for 'politikredskode' to the AdresseQuery.
//
// Find de adresser som ligger indenfor politikredsen angivet ved politikredskoden.
// (Flerværdisøgning mulig). Søgning efter ingen værdi mulig.
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Politikredskode(s ...string) *AdresseQuery {
	q.add(textQuery{Name: "politikredskode", Values: s, Multi: true, Null: true})
	return q
}

// Side will add a parameter for 'side' to the AdresseQuery.
//
// Angiver hvilken siden som skal leveres. Se Paginering.
// http://dawa.aws.dk/generelt#paginering
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) Side(i int) *AdresseQuery {
	q.add(textQuery{Name: "side", Values: []string{strconv.Itoa(i)}, Multi: true, Null: true})
	return q
}

// PerSide will add a parameter for 'per_side' to the AdresseQuery.
//
// Antal resultater per side. Se Paginering.
// http://dawa.aws.dk/generelt#paginering
//
// See documentation at http://dawa.aws.dk/adressedok#adressesoegning
func (q *AdresseQuery) PerSide(i int) *AdresseQuery {
	q.add(textQuery{Name: "per_side", Values: []string{strconv.Itoa(i)}, Multi: true, Null: true})
	return q
}

// NoFormat will disable extra whitespace. Always enabled when querying
func (q *AdresseQuery) NoFormat() *AdresseQuery {
	q.add(textQuery{Name: "noformat", Multi: false, Null: true})
	return q
}
