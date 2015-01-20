package dawa

type Postnummer struct {
	Href     string    `json:"href"`     // Postnummerets unikke URL.
	Kommuner []Kommune `json:"kommuner"` // De kommuner hvis areal overlapper postnumeret areal.
	Navn     string    `json:"navn"`     // Det navn der er knyttet til postnummeret, typisk byens eller bydelens navn. Repræsenteret ved indtil 20 tegn. Eksempel: ”København NV”.
	Nr       string    `json:"nr"`       // Unik identifikation af det postnummeret. Postnumre fastsættes af Post Danmark. Repræsenteret ved fire cifre. Eksempel: ”2400” for ”København NV”.
	// Never set to anything but null
	Stormodtageradresser interface{} `json:"stormodtageradresser"` // Hvis postnummeret er et stormodtagerpostnummer rummer feltet adresserne på stormodtageren.
}
