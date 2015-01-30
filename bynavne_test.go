package dawa

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

var suppl_bynavn_json_input = `
[
{
  "href": "http://dawa.aws.dk/supplerendebynavne/%C3%85vang",
  "navn": "Åvang",
  "postnumre": [
    {
      "href": "http://dawa.aws.dk/postnumre/4320",
      "nr": "4320",
      "navn": "Lejre"
    }
  ],
  "kommuner": [
    {
      "href": "http://dawa.aws.dk/kommuner/350",
      "kode": "0350",
      "navn": "Lejre"
    }
  ]
}, {
  "href": "http://dawa.aws.dk/supplerendebynavne/Aavang",
  "navn": "Aavang",
  "postnumre": [
    {
      "href": "http://dawa.aws.dk/postnumre/4320",
      "nr": "4320",
      "navn": "Lejre"
    }
  ],
  "kommuner": [
    {
      "href": "http://dawa.aws.dk/kommuner/350",
      "kode": "0350",
      "navn": "Lejre"
    }
  ]
}, {
  "href": "http://dawa.aws.dk/supplerendebynavne/%C3%85ved",
  "navn": "Åved",
  "postnumre": [
    {
      "href": "http://dawa.aws.dk/postnumre/6780",
      "nr": "6780",
      "navn": "Skærbæk"
    }
  ],
  "kommuner": [
    {
      "href": "http://dawa.aws.dk/kommuner/550",
      "kode": "0550",
      "navn": "Tønder"
    }
  ]
}]
`

func TestImportSupplBynavnJSON(t *testing.T) {
	// We test one entry only to match
	var json_expect = []SupplBynavn{
		SupplBynavn{Navn: "Åvang", Href: "http://dawa.aws.dk/supplerendebynavne/%C3%85vang", Kommuner: []KommuneRef{KommuneRef{Href: "http://dawa.aws.dk/kommuner/350", Kode: "0350", Navn: "Lejre"}}, Postnumre: []PostnummerRef{PostnummerRef{Href: "http://dawa.aws.dk/postnumre/4320", Navn: "Lejre", Nr: "4320"}}},
		SupplBynavn{Navn: "Aavang", Href: "http://dawa.aws.dk/supplerendebynavne/Aavang", Kommuner: []KommuneRef{KommuneRef{Href: "http://dawa.aws.dk/kommuner/350", Kode: "0350", Navn: "Lejre"}}, Postnumre: []PostnummerRef{PostnummerRef{Href: "http://dawa.aws.dk/postnumre/4320", Navn: "Lejre", Nr: "4320"}}},
		SupplBynavn{Navn: "Åved", Href: "http://dawa.aws.dk/supplerendebynavne/%C3%85ved", Kommuner: []KommuneRef{KommuneRef{Href: "http://dawa.aws.dk/kommuner/550", Kode: "0550", Navn: "Tønder"}}, Postnumre: []PostnummerRef{PostnummerRef{Href: "http://dawa.aws.dk/postnumre/6780", Navn: "Skærbæk", Nr: "6780"}}},
	}

	b := bytes.NewBuffer([]byte(suppl_bynavn_json_input))
	iter, err := ImportSupplBynavnJSON(b)
	if err != nil {
		t.Fatalf("ImportSupplBynavnJSON: %v", err)
	}
	for _, expect := range json_expect {
		item, err := iter.Next()
		if err != nil {
			t.Fatalf("ImportSupplBynavnJSON, iter.Next(): %v", err)
		}
		if item == nil {
			t.Fatalf("ImportSupplBynavnJSON, iter.Next() returned nil value")
		}
		if !reflect.DeepEqual(*item, expect) {
			t.Fatalf("ImportSupplBynavnJSON, value mismatch.\nGot:\n%#v\nExpected:\n%#v\n", *item, expect)
		}
	}
	// We should now have read all entries
	_, err = iter.Next()
	if err != io.EOF {
		t.Fatalf("ImportSupplBynavnJSON: Expected io.EOF, got:%v", err)
	}
}
