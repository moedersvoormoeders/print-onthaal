package data

type Ontvanger struct {
	Geslacht string `json:"geslacht"`
	Leeftijd int    `json:"leeftijd"`
	Naam     string `json:"naam"`
}

type MateriaalItem struct {
	Object          string    `json:"object"`
	Prijs           float64   `json:"prijs"`
	Maat            string    `json:"maat"`
	Opmerking       string    `json:"opmerking"`
	ExtraEscposData string    `json:"extraEscposData"` // this is raw ESC/POS data that will be printed
	SeperateReceipt bool      `json:"seperateReceipt"` // if true, this object will be printed on a seperate receipt
	Ontvanger       Ontvanger `json:"ontvanger"`
}
type MateriaalKlant struct {
	MVMNummer        string `json:"mvmNummer"`
	Naam             string `json:"naam"`
	Voornaam         string `json:"voornaam"`
	EenmaligenNummer string `json:"eenmaligenNummer"`
}

type MateriaalRequest struct {
	Klant MateriaalKlant `json:"klant"`

	Items []MateriaalItem `json:"items`
}

type RequestEenmaligen struct {
	EenmaligenNummer string `json:"eenmaligenNummer"`
	Naam             string `json:"naam"`
	Bericht          string `json:"bericht"`
}

type SinterklaasRequest struct {
	// fun fact: we designed these first to be seperate calls
	Speelgoed struct {
		MVMNummer string `json:"mvmNummer"`
		Naam      string `json:"naam"`
		Paketten  []struct {
			Naam      string  `json:"naam"`
			Geslacht  string  `json:"geslacht"`
			Leeftijd  float64 `json:"leeftijd"` // float64 really!
			Opmerking string  `json:"opmerking"`
		} `json:"paketten"`
	} `json:speelgoed`
	Snoep struct {
		MVMNummer  string `json:"mvmNummer"`
		VolgNummer int    `json:"volgNummer"`
		Naam       string `json:"naam"`
		Personen   int    `json:"personen"`
	} `json:snoep`
}

type MarktTicketKind struct {
	Naam     string `json:"naam"`
	Leeftijd int    `json:"leeftijd"`
	Geslacht string `json:"geslacht"`
}

type MarktRequest struct {
	Naam         string `json:"naam"`
	Beschrijving string `json:"beschrijving"`

	MVMNummer string `json:"mvmnummer"`

	Kinderen []MarktTicketKind `json:"kinderen"`

	MarktID uint `json:"markt_id"`

	Barcode string `json:"barcode"`
}
