package dawa

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

var vejstykker_json_input = `
[
{
  "href": "http://dawa.aws.dk/vejstykker/563/9369",
  "kode": "9369",
  "navn": "Vesten Bavnen",
  "adresseringsnavn": "Vesten Bavnen",
  "kommune": {
    "href": "http://dawa.aws.dk/kommuner/563",
    "kode": "0563",
    "navn": "Fanø"
  },
  "postnumre": [
    {
      "href": "http://dawa.aws.dk/postnumre/6720",
      "nr": "6720",
      "navn": "Fanø"
    }
  ],
  "historik": {
    "oprettet": "2010-01-17T11:19:52.237",
    "ændret": "2010-01-17T11:19:52.237"
  }
}, {
  "href": "http://dawa.aws.dk/vejstykker/563/9379",
  "kode": "9379",
  "navn": "Vesten Sandene",
  "adresseringsnavn": "Vesten Sandene",
  "kommune": {
    "href": "http://dawa.aws.dk/kommuner/563",
    "kode": "0563",
    "navn": "Fanø"
  },
  "postnumre": [
    {
      "href": "http://dawa.aws.dk/postnumre/6720",
      "nr": "6720",
      "navn": "Fanø"
    }
  ],
  "historik": {
    "oprettet": "2010-01-17T11:19:52.237",
    "ændret": "2010-01-17T11:19:52.237"
  }
}, {
  "href": "http://dawa.aws.dk/vejstykker/563/9389",
  "kode": "9389",
  "navn": "Vesterengen",
  "adresseringsnavn": "Vesterengen",
  "kommune": {
    "href": "http://dawa.aws.dk/kommuner/563",
    "kode": "0563",
    "navn": "Fanø"
  },
  "postnumre": [
    {
      "href": "http://dawa.aws.dk/postnumre/6720",
      "nr": "6720",
      "navn": "Fanø"
    }
  ],
  "historik": {
    "oprettet": "2010-01-17T11:19:52.253",
    "ændret": "2010-01-17T11:19:52.253"
  }
}
]
`

func TestImportVejstykkerJSON(t *testing.T) {
	// We test one entry only to match
	var json_expect = []Vejstykke{
		Vejstykke{
			Adresseringsnavn: "Vesten Bavnen",
			Historik: Historik{
				Oprettet: MustParseTime("2010-01-17T11:19:52.237"),
				Ændret:   MustParseTime("2010-01-17T11:19:52.237"),
			},
			Href: "http://dawa.aws.dk/vejstykker/563/9369",
			Kode: "9369",
			Kommune: Kommune{
				Href: "http://dawa.aws.dk/kommuner/563",
				Kode: "0563",
				Navn: "Fanø",
			},
			Navn: "Vesten Bavnen",
			Postnumre: []PostnummerRef{
				PostnummerRef{
					Href: "http://dawa.aws.dk/postnumre/6720",
					Navn: "Fanø",
					Nr:   "6720",
				},
			},
		},
		Vejstykke{
			Adresseringsnavn: "Vesten Sandene",
			Historik: Historik{
				Oprettet: MustParseTime("2010-01-17T11:19:52.237"),
				Ændret:   MustParseTime("2010-01-17T11:19:52.237"),
			},
			Href: "http://dawa.aws.dk/vejstykker/563/9379",
			Kode: "9379",
			Kommune: Kommune{
				Href: "http://dawa.aws.dk/kommuner/563",
				Kode: "0563",
				Navn: "Fanø",
			},
			Navn: "Vesten Sandene",
			Postnumre: []PostnummerRef{
				PostnummerRef{
					Href: "http://dawa.aws.dk/postnumre/6720",
					Navn: "Fanø",
					Nr:   "6720",
				},
			},
		},
		Vejstykke{
			Adresseringsnavn: "Vesterengen",
			Historik: Historik{
				Oprettet: MustParseTime("2010-01-17T11:19:52.253"),
				Ændret:   MustParseTime("2010-01-17T11:19:52.253"),
			},
			Href: "http://dawa.aws.dk/vejstykker/563/9389",
			Kode: "9389",
			Kommune: Kommune{
				Href: "http://dawa.aws.dk/kommuner/563",
				Kode: "0563",
				Navn: "Fanø",
			},
			Navn: "Vesterengen",
			Postnumre: []PostnummerRef{
				PostnummerRef{
					Href: "http://dawa.aws.dk/postnumre/6720",
					Navn: "Fanø",
					Nr:   "6720",
				},
			},
		},
	}

	b := bytes.NewBuffer([]byte(vejstykker_json_input))
	iter, err := ImportVejstykkerJSON(b)
	if err != nil {
		t.Fatalf("ImportVejstykkerJSON: %v", err)
	}
	for _, expect := range json_expect {
		item, err := iter.Next()
		if err != nil {
			t.Fatalf("ImportVejstykkerJSON, iter.Next(): %v", err)
		}
		if item == nil {
			t.Fatalf("ImportVejstykkerJSON, iter.Next() returned nil value")
		}
		if !reflect.DeepEqual(*item, expect) {
			t.Fatalf("ImportVejstykkerJSON, value mismatch.\nGot:\n%#v\nExpected:\n%#v\n", *item, expect)
		}
	}
	// We should now have read all entries
	_, err = iter.Next()
	if err != io.EOF {
		t.Fatalf("ImportVejstykkerJSON: Expected io.EOF, got:%v", err)
	}
}
