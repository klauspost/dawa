package dawa

// Et vejstykke er en vej, som er afgrænset af en kommune.
// Et vejstykke er identificeret ved en kommunekode og en vejkode og har desuden et navn.
// En vej som gennemløber mere end en kommune vil bestå af flere vejstykker.
// Det er p.t. ikke muligt at få information om hvilke vejstykker der er en del af den samme vej.
// Vejstykker er udstillet under /vejstykker
type Vejstykke struct {
	Adresseringsnavn string          `json:"adresseringsnavn"` //En evt. forkortet udgave af vejnavnet på højst 20 tegn, som bruges ved adressering på labels og rudekuverter og lign., hvor der ikke plads til det fulde vejnavn.
	Historik         Historik        `json:"historik"`         // Væsentlige tidspunkter for vejstykket
	Href             string          `json:"href"`             // Vejstykkets unikke URL.
	Kode             string          `json:"kode"`             // Identifikation af vejstykke. Er unikt indenfor den pågældende kommune. Repræsenteret ved fire cifre. Eksempel: I Københavns kommune er ”0004” lig ”Abel Cathrines Gade”.
	Kommune          Kommune         `json:"kommune"`          // Kommunen som vejstykket er beliggende i.
	Navn             string          `json:"navn"`             // Vejens navn som det er fastsat og registreret af kommunen. Repræsenteret ved indtil 40 tegn. Eksempel: ”Hvidkildevej”.
	Postnumre        []PostnummerRef `json:"postnumre"`        // Postnummrene som vejstykket er beliggende i.
}
