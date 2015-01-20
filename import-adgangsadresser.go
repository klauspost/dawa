package dawa

import (
	"bufio"
	"encoding/csv"
	"github.com/ugorji/go/codec"
	"io"
	"strconv"
)

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

// ImportAdresserJSON will import "adresser" from a JSON input, supplied to the reader.
// An iterator will be returned that return all addresses.
func ImportAdgangsAdresserJSON(in io.Reader) (*AdgangsAdresseIter, error) {
	reader := bufio.NewReader(in)
	// Skip until after '['
	_, err := reader.ReadBytes('[')
	if err != nil {
		return nil, err
	}
	// Start decoder
	ret := &AdgangsAdresseIter{a: make(chan AdgangsAdresse, 100)}
	go func() {
		defer close(ret.a)
		var h codec.JsonHandle
		h.ErrorIfNoField = true
		for {
			var dec *codec.Decoder = codec.NewDecoder(reader, &h)
			a := AdgangsAdresse{}
			if err := dec.Decode(&a); err != nil {
				ret.err = err
				return
			}
			ret.a <- a

			// Skip comma
			if b, err := readByteSkippingSpace(reader); err != nil {
				ret.err = err
				return
			} else {
				switch b {
				case ',':
					continue
				case ']':
					ret.err = io.EOF
					return
				default:
					panic("Invalid character in JSON data: " + string([]byte{b}))
				}
			}

		}
	}()
	return ret, nil
}
