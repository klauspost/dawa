package dawa

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

var postnumre_json_input = `
[
 {
  "href": "http://dawa.aws.dk/postnumre/9981",
  "nr": "9981",
  "navn": "Jerup",
  "stormodtageradresser": null,
  "kommuner": [
    {
      "href": "http://dawa.aws.dk/kommuner/813",
      "kode": "0813",
      "navn": "Frederikshavn"
    }
  ]
}, {
  "href": "http://dawa.aws.dk/postnumre/9982",
  "nr": "9982",
  "navn": "Ålbæk",
  "stormodtageradresser": null,
  "kommuner": [
    {
      "href": "http://dawa.aws.dk/kommuner/813",
      "kode": "0813",
      "navn": "Frederikshavn"
    },
    {
      "href": "http://dawa.aws.dk/kommuner/860",
      "kode": "0860",
      "navn": "Hjørring"
    }
  ]
}, {
  "href": "http://dawa.aws.dk/postnumre/9990",
  "nr": "9990",
  "navn": "Skagen",
  "stormodtageradresser": null,
  "kommuner": [
    {
      "href": "http://dawa.aws.dk/kommuner/813",
      "kode": "0813",
      "navn": "Frederikshavn"
    }
  ]
}
]
`

func TestImportPostnumreJSON(t *testing.T) {
	// We test all entries
	var json_expect = []Postnummer{
		Postnummer{
			Href: "http://dawa.aws.dk/postnumre/9981",
			Kommuner: []KommuneRef{
				KommuneRef{
					Href: "http://dawa.aws.dk/kommuner/813",
					Kode: "0813",
					Navn: "Frederikshavn",
				},
			},
			Navn:                 "Jerup",
			Nr:                   "9981",
			Stormodtageradresser: nil,
		},
		Postnummer{
			Href: "http://dawa.aws.dk/postnumre/9982",
			Kommuner: []KommuneRef{
				KommuneRef{
					Href: "http://dawa.aws.dk/kommuner/813",
					Kode: "0813",
					Navn: "Frederikshavn",
				},
				KommuneRef{
					Href: "http://dawa.aws.dk/kommuner/860",
					Kode: "0860",
					Navn: "Hjørring",
				},
			},
			Navn:                 "Ålbæk",
			Nr:                   "9982",
			Stormodtageradresser: nil,
		},
		Postnummer{
			Href: "http://dawa.aws.dk/postnumre/9990",
			Kommuner: []KommuneRef{
				KommuneRef{
					Href: "http://dawa.aws.dk/kommuner/813",
					Kode: "0813",
					Navn: "Frederikshavn",
				},
			},
			Navn:                 "Skagen",
			Nr:                   "9990",
			Stormodtageradresser: nil,
		},
	}

	b := bytes.NewBuffer([]byte(postnumre_json_input))
	iter, err := ImportPostnumreJSON(b)
	if err != nil {
		t.Fatalf("ImportPostnumreJSON: %v", err)
	}
	for _, expect := range json_expect {
		item, err := iter.Next()
		if err != nil {
			t.Fatalf("ImportPostnumreJSON, iter.Next(): %v", err)
		}
		if item == nil {
			t.Fatalf("ImportPostnumreJSON, iter.Next() returned nil value")
		}

		if !reflect.DeepEqual(*item, expect) {
			t.Fatalf("ImportPostnumreJSON, value mismatch.\nGot:\n%#v\nExpected:\n%#v\n", *item, expect)
		}

	}
	// We should now have read all entries
	_, err = iter.Next()
	if err != io.EOF {
		t.Fatalf("ImportPostnumreJSON: Expected io.EOF, got:%v", err)
	}
}
