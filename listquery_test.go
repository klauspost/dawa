package dawa

import (
	"net/url"
	"strings"
	"testing"
)

func TestListQueryParameters(t *testing.T) {
	generated := NewListQuery("opstillingskredse", false).Kode(testParameters...).URL()
	parsed, err := url.Parse(generated)
	if err != nil {
		t.Fatal(err)
	}
	values := parsed.Query()
	param, ok := values["kode"]
	if !ok {
		t.Fatal("Unable to find expected field 'kode'")
	}
	if len(param) != 1 {
		t.Fatalf("Number of parameters expected to be 1, was %d", len(param))
	}
	expect := `|1234|abcdef|æøå|"?!"#¤=%&|&|234.0|"abc"`
	if param[0] != expect {
		t.Fatalf("Unexpected value of parameter:\n     Was:\t%s\nExpected:\t%s", param[0], expect)
	}
	// We don't test encoding, we expect Go to handle that.
}

// TODO:Test with other Hostnames

var ListURL = []qb{
	// No parameters.

	qb{NewListQuery("regioner", false).URL(), DefaultHost + "/regioner"},
	qb{NewListQuery("sogne", false).URL(), DefaultHost + "/sogne"},
	qb{NewListQuery("retskredse", false).URL(), DefaultHost + "/retskredse"},
	qb{NewListQuery("politikredse", false).URL(), DefaultHost + "/politikredse"},
	qb{NewListQuery("opstillingskredse", false).URL(), DefaultHost + "/opstillingskredse"},
	qb{NewListQuery("valglandsdele", false).URL(), DefaultHost + "/valglandsdele"},
	qb{NewListQuery("ejerlav", false).URL(), DefaultHost + "/ejerlav"},
	qb{NewListQuery("adgangsadresser", false).URL(), DefaultHost + "/adgangsadresser"},
	qb{NewListQuery("adresser", false).URL(), DefaultHost + "/adresser"},
	qb{NewListQuery("postnumre", false).URL(), DefaultHost + "/postnumre"},

	// It should pass through unknown types.
	qb{NewListQuery("unknown", false).URL(), DefaultHost + "/unknown"},

	// Autocomplete
	qb{NewListQuery("regioner", true).URL(), DefaultHost + "/regioner/autocomplete"},

	// Test params
	qb{NewListQuery("regioner", false).Q(singleParam).URL(), DefaultHost + "/regioner?q=" + singleEncoded + ""},
	qb{NewListQuery("regioner", false).Kode(multiParam...).URL(), DefaultHost + "/regioner?kode=" + multiEncoded + ""},
	qb{NewListQuery("regioner", false).Navn(singleParam).URL(), DefaultHost + "/regioner?navn=" + singleEncoded + ""},
	qb{NewListQuery("regioner", false).NoFormat().URL(), DefaultHost + "/regioner?noformat="},

	// Test multiparam
	qb{NewListQuery("regioner", false).Q(singleParam).Kode(multiParam...).Navn(singleParam).NoFormat().URL(),
		DefaultHost + "/regioner?q=" + singleEncoded + "&kode=" + multiEncoded + "&navn=" + singleEncoded + "&noformat="},
}

func TestListQueryURL(t *testing.T) {
	for _, q := range ListURL {
		if q.Got != q.Expected {
			u := q.Got[19:]
			u = strings.Replace(u, multiEncoded, `"+ multiEncoded+"`, -1)
			u = strings.Replace(u, singleEncoded, `"+ singleEncoded+"`, -1)
			u = strings.Replace(u, intEncoded, `"+ intEncoded+"`, -1)
			t.Fatalf("Unexpected value of parameter:\n     Was:\t%s\nExpected:\t%s\nUrl:%s", q.Got, q.Expected, u)
		}
	}

}

var listTypes = []string{"regioner", "sogne", "retskredse", "politikredse", "opstillingskredse", "valglandsdele", "ejerlav", "adgangsadresser", "adresser", "postnumre"}

func TestListQueryTypes(t *testing.T) {
	for _, name := range listTypes {
		l := NewListQuery(name, false)
		if l.Type() == nil {
			t.Fatalf("list type did not resolve:%s", name)
		}
	}
	// Test that an unknown type is truly unknown
	l := NewListQuery("unknown", false).Type()
	if l != nil {
		t.Fatalf("expected unknown type to return nil, but returned %T", l)
	}

}
