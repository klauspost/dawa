package dawa

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	//"github.com/gobs/pretty"
)

var adgangs_csv_data = `id,status,oprettet,ændret,vejkode,vejnavn,husnr,supplerendebynavn,postnr,postnrnavn,kommunekode,kommunenavn,ejerlavkode,ejerlavnavn,matrikelnr,esrejendomsnr,etrs89koordinat_øst,etrs89koordinat_nord,wgs84koordinat_bredde,wgs84koordinat_længde,nøjagtighed,kilde,tekniskstandard,tekstretning,adressepunktændringsdato,ddkn_m100,ddkn_km1,ddkn_km10,kvh,regionskode,regionsnavn,sognekode,sognenavn,politikredskode,politikredsnavn,retskredskode,retskredsnavn,opstillingskredskode,opstillingskredsnavn,zone
0a3f507a-3669-32b8-e044-0003ba298018,1,2000-02-05T20:17:59.000,2009-11-25T01:07:37.000,0004,Abel Cathrines Gade,3A,,1654,København V,0101,København,2000174,"Udenbys Vester Kvarter, København",377,9343,723743.16,6175322.16,55.6720594006065,12.5582458296225,A,5,TD,200,2002-04-07T00:00:00.000,100m_61753_7237,1km_6175_723,10km_617_72,01010004__3A,1084,Region Hovedstaden,9185,Vesterbro,1470,Københavns Politi,1101,Københavns Byret,0009,Vesterbro,Byzone
0a3f507a-3680-32b8-e044-0003ba298018,1,2000-12-14T10:07:42.000,2009-11-25T01:07:37.000,0004,Abel Cathrines Gade,28,,1654,København V,0101,København,2000174,"Udenbys Vester Kvarter, København",70æ,342514,723928.98,6175195.49,55.6708378041398,12.5610912697286,A,5,TD,200,2002-04-05T00:00:00.000,100m_61751_7239,1km_6175_723,10km_617_72,01010004__28,1084,Region Hovedstaden,9185,Vesterbro,1470,Københavns Politi,1101,Københavns Byret,0009,Vesterbro,Byzone
0a3f507a-367c-32b8-e044-0003ba298018,1,2000-02-05T20:27:56.000,2009-11-25T01:07:37.000,0004,Abel Cathrines Gade,24,,1654,København V,0101,København,2000174,"Udenbys Vester Kvarter, København",70æ,342514,723920.84,6175201.53,55.6708957228815,12.5609670270951,A,5,TD,200,2002-04-07T00:00:00.000,100m_61752_7239,1km_6175_723,10km_617_72,01010004__24,1084,Region Hovedstaden,9185,Vesterbro,1470,Københavns Politi,1101,Københavns Byret,0009,Vesterbro,Byzone
`

func TestImportAdgangsAdresserCSV(t *testing.T) {
	// We test one entry only to match
	csv_expect := []AdgangsAdresse{
		AdgangsAdresse{
			DDKN: DDKN{
				Km1:  "1km_6175_723",
				Km10: "10km_617_72",
				M100: "100m_61753_7237",
			},
			Adgangspunkt: Adgangspunkt{
				Kilde: 0,
				Koordinater: []float64{
					55.6720594006065,
					12.5582458296225,
				},
				Nøjagtighed:     "A",
				Tekniskstandard: "TD",
				Tekstretning:    200,
				Ændret:          MustParseTime("2002-04-07T00:00:00.000"),
			},
			Ejerlav: Ejerlav{
				Kode: 2000174,
				Navn: "Udenbys Vester Kvarter, København",
			},
			EsrEjendomsNr: "9343",
			Historik: Historik{
				Oprettet: MustParseTime("2000-02-05T20:17:59.000"),
				Ændret:   MustParseTime("2009-11-25T01:07:37.000"),
			},
			Href:  "",
			Husnr: "3A",
			ID:    "0a3f507a-3669-32b8-e044-0003ba298018",
			Kommune: Kommune{
				Href: "",
				Kode: "0101",
				Navn: "København",
			},
			Kvh:        "",
			Matrikelnr: "377",
			Opstillingskreds: Opstillingskreds{
				Href: "",
				Kode: "0009",
				Navn: "Vesterbro",
			},
			Politikreds: Politikreds{
				Href: "",
				Kode: "1470",
				Navn: "Københavns Politi",
			},
			Postnummer: PostnummerRef{
				Href: "",
				Navn: "København V",
				Nr:   "1654",
			},
			Region: Region{
				Href: "",
				Kode: "1084",
				Navn: "Region Hovedstaden",
			},
			Retskreds: Retskreds{
				Href: "",
				Kode: "1101",
				Navn: "Københavns Byret",
			},
			Sogn: Sogn{
				Href: "",
				Kode: "9185",
				Navn: "Vesterbro",
			},
			Status:            1,
			SupplerendeBynavn: "",
			Vejstykke: VejstykkeRef{
				Href: "",
				Kode: "0004",
				Navn: "Abel Cathrines Gade",
			},
			Zone: "Byzone",
		},
		AdgangsAdresse{},
		AdgangsAdresse{},
	}
	b := bytes.NewBuffer([]byte(adgangs_csv_data))
	iter, err := ImportAdgangsAdresserCSV(b)
	if err != nil {
		t.Fatalf("ImportAdgangsAdresserCSV: %v", err)
	}
	for i, expect := range csv_expect {
		item, err := iter.Next()
		if err != nil {
			t.Fatalf("ImportAdgangsAdresserCSV, iter.Next(): %v", err)
		}
		if item == nil {
			t.Fatalf("ImportAdgangsAdresserCSV, iter.Next() returned nil value")
		}
		if i == 0 && !reflect.DeepEqual(*item, expect) {
			t.Fatalf("ImportAdgangsAdresserCSV, value mismatch.\nGot:\n%#v\nExpected:\n%#v\n", *item, expect)
		}
		// Test time parsing.
		if i == 0 {
			if item.Historik.Oprettet.Time().Unix() != 949778279 {
				t.Fatalf("Timestamp mismatch, expected 949778279, was %d", item.Historik.Oprettet.Time().Unix())
			}
		}
	}
	// We should now have read all entries
	_, err = iter.Next()
	if err != io.EOF {
		t.Fatalf("ImportAdresserCSV: Expected io.EOF, got:%v", err)
	}
}

var adgangs_json_input = `[
{
  "href": "http://dawa.aws.dk/adgangsadresser/0a3f507a-3669-32b8-e044-0003ba298018",
  "id": "0a3f507a-3669-32b8-e044-0003ba298018",
  "kvh": "01010004__3A",
  "status": 1,
  "vejstykke": {
    "href": "http://dawa.aws.dk/vejstykker/101/4",
    "navn": "Abel Cathrines Gade",
    "kode": "0004"
  },
  "husnr": "3A",
  "supplerendebynavn": null,
  "postnummer": {
    "href": "http://dawa.aws.dk/postnumre/1654",
    "nr": "1654",
    "navn": "København V"
  },
  "kommune": {
    "href": "http://dawa.aws.dk/kommuner/101",
    "kode": "0101",
    "navn": "København"
  },
  "ejerlav": {
    "kode": 2000174,
    "navn": "Udenbys Vester Kvarter, København"
  },
  "esrejendomsnr": "9343",
  "matrikelnr": "377",
  "historik": {
    "oprettet": "2000-02-05T20:17:59.000",
    "ændret": "2009-11-25T01:07:37.000"
  },
  "adgangspunkt": {
    "koordinater": [
      12.5582458296225,
      55.6720594006065
    ],
    "nøjagtighed": "A",
    "kilde": 5,
    "tekniskstandard": "TD",
    "tekstretning": 200,
    "ændret": "2002-04-07T00:00:00.000"
  },
  "DDKN": {
    "m100": "100m_61753_7237",
    "km1": "1km_6175_723",
    "km10": "10km_617_72"
  },
  "sogn": {
    "kode": "9185",
    "navn": "Vesterbro",
    "href": "http://dawa.aws.dk/sogne/9185"
  },
  "region": {
    "kode": "1084",
    "navn": "Region Hovedstaden",
    "href": "http://dawa.aws.dk/regioner/1084"
  },
  "retskreds": {
    "kode": "1101",
    "navn": "Københavns Byret",
    "href": "http://dawa.aws.dk/retskredse/1101"
  },
  "politikreds": {
    "kode": "1470",
    "navn": "Københavns Politi",
    "href": "http://dawa.aws.dk/politikredse/1470"
  },
  "opstillingskreds": {
    "kode": "0009",
    "navn": "Vesterbro",
    "href": "http://dawa.aws.dk/opstillingskredse/9"
  },
  "zone": "Byzone"
}, {
  "href": "http://dawa.aws.dk/adgangsadresser/0a3f507a-3680-32b8-e044-0003ba298018",
  "id": "0a3f507a-3680-32b8-e044-0003ba298018",
  "kvh": "01010004__28",
  "status": 1,
  "vejstykke": {
    "href": "http://dawa.aws.dk/vejstykker/101/4",
    "navn": "Abel Cathrines Gade",
    "kode": "0004"
  },
  "husnr": "28",
  "supplerendebynavn": null,
  "postnummer": {
    "href": "http://dawa.aws.dk/postnumre/1654",
    "nr": "1654",
    "navn": "København V"
  },
  "kommune": {
    "href": "http://dawa.aws.dk/kommuner/101",
    "kode": "0101",
    "navn": "København"
  },
  "ejerlav": {
    "kode": 2000174,
    "navn": "Udenbys Vester Kvarter, København"
  },
  "esrejendomsnr": "342514",
  "matrikelnr": "70æ",
  "historik": {
    "oprettet": "2000-12-14T10:07:42.000",
    "ændret": "2009-11-25T01:07:37.000"
  },
  "adgangspunkt": {
    "koordinater": [
      12.5610912697286,
      55.6708378041398
    ],
    "nøjagtighed": "A",
    "kilde": 5,
    "tekniskstandard": "TD",
    "tekstretning": 200,
    "ændret": "2002-04-05T00:00:00.000"
  },
  "DDKN": {
    "m100": "100m_61751_7239",
    "km1": "1km_6175_723",
    "km10": "10km_617_72"
  },
  "sogn": {
    "kode": "9185",
    "navn": "Vesterbro",
    "href": "http://dawa.aws.dk/sogne/9185"
  },
  "region": {
    "kode": "1084",
    "navn": "Region Hovedstaden",
    "href": "http://dawa.aws.dk/regioner/1084"
  },
  "retskreds": {
    "kode": "1101",
    "navn": "Københavns Byret",
    "href": "http://dawa.aws.dk/retskredse/1101"
  },
  "politikreds": {
    "kode": "1470",
    "navn": "Københavns Politi",
    "href": "http://dawa.aws.dk/politikredse/1470"
  },
  "opstillingskreds": {
    "kode": "0009",
    "navn": "Vesterbro",
    "href": "http://dawa.aws.dk/opstillingskredse/9"
  },
  "zone": "Byzone"
}, {
  "href": "http://dawa.aws.dk/adgangsadresser/0a3f507a-367c-32b8-e044-0003ba298018",
  "id": "0a3f507a-367c-32b8-e044-0003ba298018",
  "kvh": "01010004__24",
  "status": 1,
  "vejstykke": {
    "href": "http://dawa.aws.dk/vejstykker/101/4",
    "navn": "Abel Cathrines Gade",
    "kode": "0004"
  },
  "husnr": "24",
  "supplerendebynavn": null,
  "postnummer": {
    "href": "http://dawa.aws.dk/postnumre/1654",
    "nr": "1654",
    "navn": "København V"
  },
  "kommune": {
    "href": "http://dawa.aws.dk/kommuner/101",
    "kode": "0101",
    "navn": "København"
  },
  "ejerlav": {
    "kode": 2000174,
    "navn": "Udenbys Vester Kvarter, København"
  },
  "esrejendomsnr": "342514",
  "matrikelnr": "70æ",
  "historik": {
    "oprettet": "2000-02-05T20:27:56.000",
    "ændret": "2009-11-25T01:07:37.000"
  },
  "adgangspunkt": {
    "koordinater": [
      12.5609670270951,
      55.6708957228815
    ],
    "nøjagtighed": "A",
    "kilde": 5,
    "tekniskstandard": "TD",
    "tekstretning": 200,
    "ændret": "2002-04-07T00:00:00.000"
  },
  "DDKN": {
    "m100": "100m_61752_7239",
    "km1": "1km_6175_723",
    "km10": "10km_617_72"
  },
  "sogn": {
    "kode": "9185",
    "navn": "Vesterbro",
    "href": "http://dawa.aws.dk/sogne/9185"
  },
  "region": {
    "kode": "1084",
    "navn": "Region Hovedstaden",
    "href": "http://dawa.aws.dk/regioner/1084"
  },
  "retskreds": {
    "kode": "1101",
    "navn": "Københavns Byret",
    "href": "http://dawa.aws.dk/retskredse/1101"
  },
  "politikreds": {
    "kode": "1470",
    "navn": "Københavns Politi",
    "href": "http://dawa.aws.dk/politikredse/1470"
  },
  "opstillingskreds": {
    "kode": "0009",
    "navn": "Vesterbro",
    "href": "http://dawa.aws.dk/opstillingskredse/9"
  },
  "zone": "Byzone"
}
]`

func TestImportAdgangsAdresserJSON(t *testing.T) {
	// We test one entry only to match
	var json_expect = []AdgangsAdresse{
		AdgangsAdresse{
			DDKN: DDKN{
				Km1:  "1km_6175_723",
				Km10: "10km_617_72",
				M100: "100m_61753_7237",
			},
			Adgangspunkt: Adgangspunkt{
				Kilde: 5,
				Koordinater: []float64{
					12.5582458296225,
					55.6720594006065,
				},
				Nøjagtighed:     "A",
				Tekniskstandard: "TD",
				Tekstretning:    200,
				Ændret:          MustParseTime("2002-04-07T00:00:00.000"),
			},
			Ejerlav: Ejerlav{
				Kode: 2000174,
				Navn: "Udenbys Vester Kvarter, København",
			},
			EsrEjendomsNr: "9343",
			Historik: Historik{
				Oprettet: MustParseTime("2000-02-05T20:17:59.000"),
				Ændret:   MustParseTime("2009-11-25T01:07:37.000"),
			},
			Href:  "http://dawa.aws.dk/adgangsadresser/0a3f507a-3669-32b8-e044-0003ba298018",
			Husnr: "3A",
			ID:    "0a3f507a-3669-32b8-e044-0003ba298018",
			Kommune: Kommune{
				Href: "http://dawa.aws.dk/kommuner/101",
				Kode: "0101",
				Navn: "København",
			},
			Kvh:        "01010004__3A",
			Matrikelnr: "377",
			Opstillingskreds: Opstillingskreds{
				Href: "http://dawa.aws.dk/opstillingskredse/9",
				Kode: "0009",
				Navn: "Vesterbro",
			},
			Politikreds: Politikreds{
				Href: "http://dawa.aws.dk/politikredse/1470",
				Kode: "1470",
				Navn: "Københavns Politi",
			},
			Postnummer: PostnummerRef{
				Href: "http://dawa.aws.dk/postnumre/1654",
				Navn: "København V",
				Nr:   "1654",
			},
			Region: Region{
				Href: "http://dawa.aws.dk/regioner/1084",
				Kode: "1084",
				Navn: "Region Hovedstaden",
			},
			Retskreds: Retskreds{
				Href: "http://dawa.aws.dk/retskredse/1101",
				Kode: "1101",
				Navn: "Københavns Byret",
			},
			Sogn: Sogn{
				Href: "http://dawa.aws.dk/sogne/9185",
				Kode: "9185",
				Navn: "Vesterbro",
			},
			Status:            1,
			SupplerendeBynavn: "",
			Vejstykke: VejstykkeRef{
				Href: "http://dawa.aws.dk/vejstykker/101/4",
				Kode: "0004",
				Navn: "Abel Cathrines Gade",
			},
			Zone: "Byzone",
		},
		AdgangsAdresse{},
		AdgangsAdresse{},
	}

	b := bytes.NewBuffer([]byte(adgangs_json_input))
	iter, err := ImportAdgangsAdresserJSON(b)
	if err != nil {
		t.Fatalf("ImportAdgangsAdresserJSON: %v", err)
	}
	for i, expect := range json_expect {
		item, err := iter.Next()
		if err != nil {
			t.Fatalf("ImportAdgangsAdresserJSON, iter.Next(): %v", err)
		}
		if item == nil {
			t.Fatalf("ImportAdgangsAdresserJSON, iter.Next() returned nil value")
		}
		if i == 0 && !reflect.DeepEqual(*item, expect) {
			t.Fatalf("ImportAdgangsAdresserJSON, value mismatch.\nGot:\n%#v\nExpected:\n%#v\n", *item, expect)
		}
		// Since we leak time parsing abstraction, we need to test a value.
		if i == 0 {
			if item.Historik.Oprettet.Time().Unix() != 949778279 {
				t.Fatalf("Timestamp mismatch, expected 949778279, was %d", item.Historik.Oprettet.Time().Unix())
			}
		}
	}
	// We should now have read all entries
	_, err = iter.Next()
	if err != io.EOF {
		t.Fatalf("ImportAdgangsAdresserJSON: Expected io.EOF, got:%v", err)
	}
}
