package dawa

type Adresse struct {
	Adgangsadresse    AdgangsAdresse `json:"adgangsadresse"`    // Adressens adgangsadresse
	Adressebetegnelse string         `json:"adressebetegnelse"` // (unknown)
	Dør               string         `json:"dør"`               // Dørbetegnelse. Tal fra 1 til 9999, små og store bogstaver samt tegnene / og -.
	Etage             string         `json:"etage"`             // Etagebetegnelse. Hvis værdi angivet kan den antage følgende værdier: tal fra 1 til 99, st, kl, kl2 op til kl9.
	Historik          Historik       `json:"historik"`          // Væsentlige tidspunkter for adressen
	Href              string         `json:"href"`              // Adgangsadressens URL.
	ID                string         `json:"id"`                // Adressens unikke id, f.eks. 0a3f5095-45ec-32b8-e044-0003ba298018.
	Kvhx              string         `json:"kvhx"`              // KVHX-nøgle. 19 tegn bestående af 4 cifre der repræsenterer kommunekode, 4 cifre der repræsenterer vejkode, 4 tegn der repræsenter husnr, 3 tegn der repræsenterer etage og 4 tegn der repræsenter dør.
	Status            int            `json:"status"`            // Adressens status. 1 indikerer en gældende adresse, 3 indikerer en foreløbig adresse.
}
