package dawa

// Et supplerende bynavn – typisk landsbyens navn – eller andet lokalt stednavn
// der er fastsat af kommunen for at præcisere adressens beliggenhed indenfor postnummeret.
//
// Indgår som en del af den officielle adressebetegnelse.
type SupplBynavn struct {
	Navn      string          `json:"navn"`      // Det supplerende bynavn. Indtil 34 tegn. Eksempel: ”Sønderholm”.
	Href      string          `json:"href"`      // Det supplerende bynavns unikke URL
	Kommuner  []Kommune       `json:"kommuner"`  // Kommuner, som det supplerende bynavn er beliggende i.
	Postnumre []PostnummerRef `json:"postnumre"` // Postnumre, som det supplerende bynavn er beliggende i.
}
