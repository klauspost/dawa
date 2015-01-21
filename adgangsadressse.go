package dawa

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
