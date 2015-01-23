package dawa

import (
	"bufio"
	"encoding/csv"
	"github.com/ugorji/go/codec"
	"io"
	"strconv"
)

// En adgangsadresse er en struktureret betegnelse som angiver en særskilt
// adgang til et areal eller en bygning efter reglerne i adressebekendtgørelsen.
//
// Forskellen på en adresse og en adgangsadresse er at adressen rummer
// eventuel etage- og/eller dørbetegnelse. Det gør adgangsadressen ikke.
type AdgangsAdresse struct {
	DDKN              DDKN             `json:"DDKN"`              // Adressens placering i Det Danske Kvadratnet (DDKN).
	Adgangspunkt      Adgangspunkt     `json:"adgangspunkt"`      // Geografisk punkt, som angiver særskilt adgang fra navngiven vej ind på et areal eller bygning.
	Ejerlav           Ejerlav          `json:"ejerlav"`           // Det matrikulære ejerlav som adressen ligger i.
	EsrEjendomsNr     string           `json:"esrejendomsnr"`     // ESR Ejendomsnummer. Indtil 7 cifre.
	Historik          Historik         `json:"historik"`          // Væsentlige tidspunkter for adgangsadressen
	Href              string           `json:"href"`              // Adgangsadressens URL.
	Husnr             string           `json:"husnr"`             // Husnummer. Max 4 cifre eventuelt med et efterfølgende bogstav.
	ID                string           `json:"id"`                // Adgangsadressens unikke id, f.eks. 0a3f5095-45ec-32b8-e044-0003ba298018
	Kommune           Kommune          `json:"kommune"`           // Kommunen som adressen er beliggende i.
	Kvh               string           `json:"kvh"`               // KVH-nøgle. 12 tegn bestående af 4 cifre der repræsenterer kommunekode, 4 cifre der repræsenterer vejkode efterfulgt af 4 tegn der repræsenter husnr
	Matrikelnr        string           `json:"matrikelnr"`        // Matrikelnummer. Unikt indenfor et ejerlav.
	Opstillingskreds  Opstillingskreds `json:"opstillingskreds"`  // Opstillingskresen som adressen er beliggende i. Beregnes udfra adgangspunktet og opstillingskredsinddelingerne fra DAGI
	Politikreds       Politikreds      `json:"politikreds"`       // Politikredsen som adressen er beliggende i. Beregnes udfra adgangspunktet og politikredsinddelingerne fra DAGI
	Postnummer        PostnummerRef    `json:"postnummer"`        // Postnummeret som adressen er beliggende i.
	Region            Region           `json:"region"`            // Regionen som adressen er beliggende i. Beregnes udfra adgangspunktet og regionsinddelingerne fra DAGI
	Retskreds         Retskreds        `json:"retskreds"`         // Retskredsen som adressen er beliggende i. Beregnes udfra adgangspunktet og retskredsinddelingerne fra DAGI
	Sogn              Sogn             `json:"sogn"`              // Sognet som adressen er beliggende i. Beregnes udfra adgangspunktet og sogneinddelingerne fra DAGI
	Status            int              `json:"status"`            // Adressens status, som modtaget fra BBR. "1" angiver en endelig adresse og "3" angiver en foreløbig adresse". Adresser med status "2" eller "4" er ikke med i DAWA.
	SupplerendeBynavn string           `json:"supplerendebynavn"` // Et supplerende bynavn – typisk landsbyens navn – eller andet lokalt stednavn, der er fastsat af kommunen for at præcisere adressens beliggenhed indenfor postnummeret.
	Vejstykke         VejstykkeRef     `json:"vejstykke"`         // Vejstykket som adressen er knyttet til.
	Zone              string           `json:"zone"`              // Hvilken zone adressen ligger i. "Byzone", "Sommerhusområde" eller "Landzone". Beregnes udfra adgangspunktet og zoneinddelingerne fra PlansystemDK
}

// Adressens placering i Det Danske Kvadratnet (DDKN).
type DDKN struct {
	Km1  string `json:"km1"`
	Km10 string `json:"km10"`
	M100 string `json:"m100"`
}

// Geografisk punkt, som angiver særskilt adgang fra navngiven vej ind på et areal eller bygning.
type Adgangspunkt struct {
	Kilde           int       `json:"kilde"`           // Kode der angiver kilden til adressepunktet. Et tegn. ”1” = oprettet maskinelt fra teknisk kort; ”2” = Oprettet maskinelt fra af matrikelnummer tyngdepunkt; ”3” = Eksternt indberettet af konsulent på vegne af kommunen; ”4” = Eksternt indberettet af kommunes kortkontor o.l. ”5” = Oprettet af teknisk forvaltning."
	Koordinater     []float64 `json:"koordinater"`     // Adgangspunktets koordinater som array [x,y].  *sic*
	Nøjagtighed     string    `json:"nøjagtighed"`     // Kode der angiver nøjagtigheden for adressepunktet. Et tegn. ”A” betyder at adressepunktet er absolut placeret på et detaljeret grundkort, tyisk med en nøjagtighed bedre end +/- 2 meter. ”B” betyder at adressepunktet er beregnet – typisk på basis af matrikelkortet, således at adressen ligger midt på det pågældende matrikelnummer. I så fald kan nøjagtigheden være ringere en end +/- 100 meter afhængig af forholdene. ”U” betyder intet adressepunkt.
	Tekniskstandard string    `json:"tekniskstandard"` // Kode der angiver den specifikation adressepunktet skal opfylde. 2 tegn. ”TD” = 3 meter inde i bygningen ved det sted hvor indgangsdør e.l. skønnes placeret; ”TK” = Udtrykkelig TK-standard: 3 meter inde i bygning, midt for længste side mod vej; ”TN” Alm. teknisk standard: bygningstyngdepunkt eller blot i bygning; ”UF” = Uspecificeret/foreløbig: ikke nødvendigvis placeret i bygning."
	Tekstretning    float64   `json:"tekstretning"`    // Angiver en evt. retningsvinkel for adressen i ”gon” dvs. hvor hele cirklen er 400 gon og 200 er vandret. Værdier 0.00-400.00: Eksempel: ”128.34”.
	Ændret          AwsTime   `json:"ændret"`          // Dato for sidste ændring i adressepunktet, som registreret af BBR.
}

type Ejerlav struct {
	Kode int    `json:"kode"` // Unik identifikation af det matrikulære ”ejerlav”, som adressen ligger i. Repræsenteret ved indtil 7 cifre. Eksempel: ”170354” for ejerlavet ”Eskebjerg By, Bregninge”.
	Navn string `json:"navn"` // Det matrikulære ”ejerlav”s navn. Eksempel: ”Eskebjerg By, Bregninge”.
}

type Historik struct {
	Oprettet AwsTime `json:"oprettet"` // Dato og tid for data oprettelse,
	Ændret   AwsTime `json:"ændret"`   // Dato og tid hvor der sidst er ændret i data,
}

// Kommunen som adressen er beliggende i.
type Kommune struct {
	Href string `json:"href"` // Kommunens unikke URL.
	Kode string `json:"kode"` // Kommunekoden. 4 cifre.
	Navn string `json:"navn"` // Kommunens navn.
}

type Opstillingskreds struct {
	Href string `json:"href"` // Opstillingskredsens unikke URL
	Kode string `json:"kode"` // Identifikation af opstillingskredsen.
	Navn string `json:"navn"` // Opstillingskredsens navn.
}

type Politikreds struct {
	Href string `json:"href"` // Politikredsens unikke URL
	Kode string `json:"kode"` // Identifikation af politikredsen
	Navn string `json:"navn"` // Politikredsens navn
}

type PostnummerRef struct {
	Href string `json:"href"` // Postnummerets unikke URL
	Navn string `json:"navn"` // Det navn der er knyttet til postnummeret, typisk byens eller bydelens navn. Repræsenteret ved indtil 20 tegn. Eksempel: ”København NV”.
	Nr   string `json:"nr"`   // Postnummer. 4 cifre
}

type Region struct {
	Href string `json:"href"` // Regionens unikke URL
	Kode string `json:"kode"` // Identifikation af regionen
	Navn string `json:"navn"` // Regionens navn
}

type Retskreds struct {
	Href string `json:"href"` // Retskredsens unikke URL
	Kode string `json:"kode"` // Identifikation af retskredsen
	Navn string `json:"navn"` // Retskredsens navn
}

type Sogn struct {
	Href string `json:"href"` // Sognets unikke URL
	Kode string `json:"kode"` // Identifikation af sognet
	Navn string `json:"navn"` // Sognets navn
}

type VejstykkeRef struct {
	Href string `json:"href"`
	Kode string `json:"kode"` // Vejkoden. 4 cifre.
	Navn string `json:"navn"` // Vejnavn. Der skelnes mellem store og små bogstaver.
}

// AdgangsAdresse is an Iterator that enable you to get individual entries.
type AdgangsAdresseIter struct {
	a   chan AdgangsAdresse
	err error
}

// Next will return addresses.
// It will return an error if that has been encountered.
// When there are not more entries nil, io.EOF will be returned.
func (a *AdgangsAdresseIter) Next() (*AdgangsAdresse, error) {
	v, ok := <-a.a
	if ok {
		return &v, nil
	}
	return nil, a.err
}

// ImportAdresserCSV will import "adresser" from a CSV file, supplied to the reader.
// An iterator will be returned that return all addresses.
func ImportAdgangsAdresserCSV(in io.Reader) (*AdgangsAdresseIter, error) {
	r := csv.NewReader(in)
	r.Comma = ','

	// Read first line as headers
	name, err := r.Read()
	if err != nil {
		return nil, err
	}

	ret := &AdgangsAdresseIter{a: make(chan AdgangsAdresse, 100)}
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(ret.a)
		v := make(map[string]string, len(name))
		for {
			records, err := r.Read()
			if err != nil {
				ret.err = err
				return
			}
			// Map to indexes, so we don't rely on index numbers, but on column names.
			for j := range records {
				v[name[j]] = records[j]
			}
			// PROCESS: id,status,oprettet,ændret,vejkode,vejnavn,husnr,etage,dør,supplerendebynavn
			a := AdgangsAdresse{}
			a.ID = v["id"]
			a.Status, err = strconv.Atoi(v["status"])
			if err != nil {
				ret.err = err
				return
			}

			// Example 2000-02-16T21:58:33.000
			o, err := ParseTime(v["oprettet"])
			if err != nil {
				ret.err = err
				return
			}
			a.Historik.Oprettet = *o

			o, err = ParseTime(v["ændret"])
			if err != nil {
				ret.err = err
				return
			}
			a.Historik.Ændret = *o

			a.Vejstykke.Kode = v["vejkode"]
			a.Vejstykke.Navn = v["vejnavn"]
			a.Husnr = v["husnr"]
			a.SupplerendeBynavn = v["supplerendebynavn"]
			a.Postnummer.Nr = v["postnr"]
			a.Postnummer.Navn = v["postnrnavn"]
			a.Kommune.Kode = v["kommunekode"]
			a.Kommune.Navn = v["kommunenavn"]
			a.Ejerlav.Kode, _ = strconv.Atoi(v["ejerlavkode"])
			a.Ejerlav.Navn = v["ejerlavnavn"]
			a.Matrikelnr = v["matrikelnr"]
			a.EsrEjendomsNr = v["esrejendomsnr"]
			// ????
			// x,_ = strconv.ParseFloat("etrs89koordinat_øst")
			// x,_ = strconv.ParseFloat("etrs89koordinat_nord")
			a.Adgangspunkt.Koordinater = make([]float64, 2)
			a.Adgangspunkt.Koordinater[0], _ = strconv.ParseFloat(v["wgs84koordinat_bredde"], 64)
			a.Adgangspunkt.Koordinater[1], _ = strconv.ParseFloat(v["wgs84koordinat_længde"], 64)

			a.Adgangspunkt.Nøjagtighed = v["nøjagtighed"]
			a.Adgangspunkt.Kilde, _ = strconv.Atoi(v["nøjagtighed"])
			a.Adgangspunkt.Tekniskstandard = v["tekniskstandard"]
			a.Adgangspunkt.Tekstretning, _ = strconv.ParseFloat(v["tekstretning"], 64)
			a.DDKN.M100 = v["ddkn_m100"]
			a.DDKN.Km1 = v["ddkn_km1"]
			a.DDKN.Km10 = v["ddkn_km10"]
			o, err = ParseTime(v["adressepunktændringsdato"])
			if err != nil {
				ret.err = err
				return
			}
			a.Adgangspunkt.Ændret = *o
			a.Region.Kode = v["regionskode"]
			a.Region.Navn = v["regionsnavn"]
			a.Sogn.Kode = v["sognekode"]
			a.Sogn.Navn = v["sognenavn"]
			a.Politikreds.Kode = v["politikredskode"]
			a.Politikreds.Navn = v["politikredsnavn"]
			a.Retskreds.Kode = v["retskredskode"]
			a.Retskreds.Navn = v["retskredsnavn"]

			// opstilli	ngskredskode,opstillingskredsnavn,zone
			a.Opstillingskreds.Kode = v["opstillingskredskode"]
			a.Opstillingskreds.Navn = v["opstillingskredsnavn"]
			a.Zone = v["zone"]
			ret.a <- a
		}
	}()
	return ret, nil
}

// ImportAdgangsAdresserJSON will import "adgangsadresser" from a JSON input, supplied to the reader.
// An iterator will be returned that return all items.
func ImportAdgangsAdresserJSON(in io.Reader) (*AdgangsAdresseIter, error) {
	var h codec.JsonHandle
	h.DecodeOptions.ErrorIfNoField = JSONStrictFieldCheck
	// use a buffered reader for efficiency
	if _, ok := in.(io.ByteScanner); !ok {
		in = bufio.NewReader(in)
	}
	ret := &AdgangsAdresseIter{a: make(chan AdgangsAdresse, 100)}
	go func() {
		defer close(ret.a)
		var dec *codec.Decoder = codec.NewDecoder(in, &h)
		ret.err = dec.Decode(&ret.a)
		if ret.err == nil {
			ret.err = io.EOF
		}
	}()

	return ret, nil
}
