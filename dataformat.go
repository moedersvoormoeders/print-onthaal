package main

type MateriaalRequest struct {
	Klant struct {
		MVMNummer string `json:"mvmNummer"`
		Naam      string `json:"naam"`
		Voornaam  string `json:"voornaam"`
	} `json:"klant"`

	Items []struct {
		Object    string  `json:"object"`
		Prijs     float64 `json:"prijs"`
		Maat      string  `json:"maat"`
		Opmerking string  `json:"maat"`
		Ontvanger struct {
			Geslacht string `json:"geslacht"`
			Leeftijd int    `json:"leeftijd"`
			Naam     string `json:"naam"`
		} `json:"ontvanger"`
	} `json:"items`
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
		MVMNummer   string `json:"mvmNummer"`
		Naam        string `json:"naam"`
		Volwassenen int    `json:"volwassenen"`
		Kinderen    int    `json:"kinderen"`
	} `json:snoep`
}
