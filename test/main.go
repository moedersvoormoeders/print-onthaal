package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mect/go-escpos"
	datapkg "github.com/moedersvoormoeders/print-onthaal/data"
)

type BufferReadWriteCloser struct {
	buffer *bytes.Buffer
}

func (b BufferReadWriteCloser) Write(p []byte) (n int, err error) {
	return b.buffer.Write(p)
}

func (b BufferReadWriteCloser) Read(p []byte) (n int, err error) {
	return b.buffer.Read(p)
}

func (b BufferReadWriteCloser) Close() error {
	return nil
}

func main() {
	// make MateriaalRequest to localhost:8080

	gordijnenEscpos := bytes.NewBuffer([]byte{})

	e, _ := escpos.NewPrinterByRW(BufferReadWriteCloser{buffer: gordijnenEscpos})

	e.Align(escpos.AlignLeft)
	e.Size(1, 1)
	e.Feed(1)
	e.PrintLn("In te vullen door de naaikamer:")
	for i := 0; i < 3; i++ {
		e.PrintLn("Formaat: ......m x .....m")
		e.PrintLn("Aantal: ......")
		e.PrintLn("Prijs: ......")
		e.PrintLn("Gemaakt: [ ]")

		e.Feed(1)
	}

	e.PrintLn("Deze bestelling is:")
	e.PrintLn("Betaald: [ ]")
	e.PrintLn("Afgehaald: [ ]")

	e.Feed(2)
	e.Size(1, 1)
	e.PrintLn("==========================================")

	e.Size(2, 2)
	e.PrintLn("MVM{{.Klant.MVMNummer}}")
	e.Size(2, 1)
	e.PrintLn("{{.Klant.Voornaam}} {{.Klant.Naam}}")
	e.Barcode("{{.Klant.MVMNummer}}", escpos.BarcodeTypeCODE39)
	e.Feed(1)
	e.PrintLn("Mededeling aan onthaal:")
	e.PrintLn("Bestelling gordijnen is afgehaald!")

	// base64 encode the escpos data
	escposencoded := base64.StdEncoding.EncodeToString(gordijnenEscpos.Bytes())
	fmt.Println(escposencoded)

	jsonData := datapkg.MateriaalRequest{
		Klant: datapkg.MateriaalKlant{
			MVMNummer: "1958",
			Naam:      "Eyskens",
			Voornaam:  "Maartje",
		},
		Items: []datapkg.MateriaalItem{
			{
				Object: "Pakket Winter",
				Maat:   "2 jaar",
				Ontvanger: datapkg.Ontvanger{
					Naam:     "Rosalien Baankaart Van Haj",
					Leeftijd: 2,
					Geslacht: "Vrouw",
				},
			},
			{
				Object:          "Overgordijnen",
				SeperateReceipt: true,
				ExtraEscposData: escposencoded,
			},
		},
	}

	// send jsonData to localhost:8080
	data, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/print", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	log.Println("response Body:", string(body))

}
