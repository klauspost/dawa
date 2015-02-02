package dawa

import (
	"net/url"
	"strings"
	"testing"
)

type qb struct {
	Got      string
	Expected string
}

var testParameters = []string{"", "1234", "abcdef", "æøå", `"?!"#¤=%&|&`, "234.0", `"abc"`}

func TestAddresseQueryParameters(t *testing.T) {
	generated := NewAdresseQuery().Postnr(testParameters...).URL()
	parsed, err := url.Parse(generated)
	if err != nil {
		t.Fatal(err)
	}
	values := parsed.Query()
	param, ok := values["postnr"]
	if !ok {
		t.Fatal("Unable to find expected field 'postnr'")
	}
	if len(param) != 1 {
		t.Fatal("Number of parameters expected to be 1, was %d", len(param))
	}
	expect := `|1234|abcdef|æøå|"?!"#¤=%&|&|234.0|"abc"`
	if param[0] != expect {
		t.Fatalf("Unexpected value of parameter:\n     Was:\t%s\nExpected:\t%s", param[0], expect)
	}
	// We don't test encoding, we expect Go to handle that.
}

// TODO:Test with other Host names
var singleParam = `Test&*æøåÆØÅ!"{.}$?=`
var singleEncoded = `Test%26%2A%C3%A6%C3%B8%C3%A5%C3%86%C3%98%C3%85%21%22%7B.%7D%24%3F%3D`
var multiParam = []string{`Mtest!"#222.%&=?`, "Seconday Param*"}
var multiEncoded = `Mtest%21%22%23222.%25%26%3D%3F|Seconday+Param%2A`
var intParam = 23453231
var intEncoded = "23453231"

var AdresseURL = []qb{
	// No parameters.
	qb{NewAdresseQuery().URL(), DefaultHost + "/adresser"},
	qb{NewAdresseComplete().URL(), DefaultHost + "/adresser/autocomplete"},

	// Single parameter
	qb{NewAdresseQuery().Cirkel(singleParam).URL(), DefaultHost + "/adresser?cirkel=" + singleEncoded + ""},
	qb{NewAdresseQuery().Dør(multiParam...).URL(), DefaultHost + "/adresser?d%C3%B8r=" + multiEncoded + ""},
	qb{NewAdresseQuery().Ejerlavkode(multiParam...).URL(), DefaultHost + "/adresser?ejerlavkode=" + multiEncoded + ""},
	qb{NewAdresseQuery().Esrejendomsnr(multiParam...).URL(), DefaultHost + "/adresser?esrejendomsnr=" + multiEncoded + ""},
	qb{NewAdresseQuery().Etage(multiParam...).URL(), DefaultHost + "/adresser?etage=" + multiEncoded + ""},
	qb{NewAdresseQuery().Husnr(multiParam...).URL(), DefaultHost + "/adresser?husnr=" + multiEncoded + ""},
	qb{NewAdresseQuery().ID(multiParam...).URL(), DefaultHost + "/adresser?id=" + multiEncoded + ""},
	qb{NewAdresseQuery().Kommunekode(multiParam...).URL(), DefaultHost + "/adresser?kommunekode=" + multiEncoded + ""},
	qb{NewAdresseQuery().Kvhx(singleParam).URL(), DefaultHost + "/adresser?kvhx=" + singleEncoded + ""},
	qb{NewAdresseQuery().Matrikelnr(multiParam...).URL(), DefaultHost + "/adresser?matrikelnr=" + multiEncoded + ""},
	qb{NewAdresseQuery().NoFormat().URL(), DefaultHost + "/adresser?noformat="},
	qb{NewAdresseQuery().Opstillingskredskode(multiParam...).URL(), DefaultHost + "/adresser?opstillingskredskode=" + multiEncoded + ""},
	qb{NewAdresseQuery().PerSide(intParam).URL(), DefaultHost + "/adresser?per_side=" + intEncoded + ""},
	qb{NewAdresseQuery().Politikredskode(multiParam...).URL(), DefaultHost + "/adresser?politikredskode=" + multiEncoded + ""},
	qb{NewAdresseQuery().Polygon(singleParam).URL(), DefaultHost + "/adresser?polygon=" + singleEncoded + ""},
	qb{NewAdresseQuery().Postnr(multiParam...).URL(), DefaultHost + "/adresser?postnr=" + multiEncoded + ""},
	qb{NewAdresseQuery().Q(singleParam).URL(), DefaultHost + "/adresser?q=" + singleEncoded + ""},
	qb{NewAdresseQuery().Regionskode(multiParam...).URL(), DefaultHost + "/adresser?regionskode=" + multiEncoded + ""},
	qb{NewAdresseQuery().Retskredskode(multiParam...).URL(), DefaultHost + "/adresser?retskredskode=" + multiEncoded + ""},
	qb{NewAdresseQuery().Side(intParam).URL(), DefaultHost + "/adresser?side=" + intEncoded + ""},
	qb{NewAdresseQuery().Sognekode(multiParam...).URL(), DefaultHost + "/adresser?sognekode=" + multiEncoded + ""},
	qb{NewAdresseQuery().Srid(singleParam).URL(), DefaultHost + "/adresser?srid=" + singleEncoded + ""},
	qb{NewAdresseQuery().Status(intParam).URL(), DefaultHost + "/adresser?status=" + intEncoded + ""},
	qb{NewAdresseQuery().SupplerendeBynavn(multiParam...).URL(), DefaultHost + "/adresser?supplerendebynavn=" + multiEncoded + ""},
	qb{NewAdresseQuery().Vejkode(multiParam...).URL(), DefaultHost + "/adresser?vejkode=" + multiEncoded + ""},
	qb{NewAdresseQuery().Vejnavn(multiParam...).URL(), DefaultHost + "/adresser?vejnavn=" + multiEncoded + ""},
	qb{NewAdresseQuery().Zonekode(multiParam...).URL(), DefaultHost + "/adresser?zonekode=" + multiEncoded + ""},

	// Combined parameters
	qb{NewAdresseQuery().Cirkel(singleParam).Etage(multiParam...).Husnr(multiParam...).Kvhx(singleParam).URL(),
		DefaultHost + "/adresser?cirkel=" + singleEncoded + "&etage=" + multiEncoded + "&husnr=" + multiEncoded + "&kvhx=" + singleEncoded + ""},

	qb{NewAdresseQuery().ID(multiParam...).Dør(multiParam...).Kommunekode(multiParam...).Opstillingskredskode(multiParam...).URL(),
		DefaultHost + "/adresser?id=" + multiEncoded + "&d%C3%B8r=" + multiEncoded + "&kommunekode=" + multiEncoded + "&opstillingskredskode=" + multiEncoded + ""},

	qb{NewAdresseQuery().Ejerlavkode(multiParam...).Matrikelnr(multiParam...).NoFormat().Postnr(multiParam...).PerSide(intParam).Politikredskode(multiParam...).URL(),
		DefaultHost + "/adresser?ejerlavkode=" + multiEncoded + "&matrikelnr=" + multiEncoded + "&noformat=&postnr=" + multiEncoded + "&per_side=" + intEncoded + "&politikredskode=" + multiEncoded + ""},

	qb{NewAdresseQuery().Esrejendomsnr(multiParam...).Polygon(singleParam).Q(singleParam).Regionskode(multiParam...).Retskredskode(multiParam...).URL(),
		DefaultHost + "/adresser?esrejendomsnr=" + multiEncoded + "&polygon=" + singleEncoded + "&q=" + singleEncoded + "&regionskode=" + multiEncoded + "&retskredskode=" + multiEncoded + ""},

	// Parameter merging
	qb{NewAdresseQuery().Vejkode(multiParam...).Zonekode(multiParam...).Vejkode("mergeme").URL(),
		DefaultHost + "/adresser?vejkode=" + multiEncoded + "|mergeme&zonekode=" + multiEncoded + ""},

	// Merge non multi
	// We expect the second Srid to be dropped
	qb{NewAdresseQuery().Srid(singleParam).Zonekode(multiParam...).Srid("dropme").URL(),
		DefaultHost + "/adresser?srid=" + singleEncoded + "&zonekode=" + multiEncoded + ""},
}

func TestAddresseQueryURL(t *testing.T) {
	for _, q := range AdresseURL {
		if q.Got != q.Expected {
			u := q.Got[19:]
			u = strings.Replace(u, multiEncoded, `"+ multiEncoded+"`, -1)
			u = strings.Replace(u, singleEncoded, `"+ singleEncoded+"`, -1)
			u = strings.Replace(u, intEncoded, `"+ intEncoded+"`, -1)
			t.Fatalf("Unexpected value of parameter:\n     Was:\t%s\nExpected:\t%s\nUrl:%s", q.Got, q.Expected, u)
		}
	}

}
