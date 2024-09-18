package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4/middleware"
	datapkg "github.com/moedersvoormoeders/print-onthaal/data"

	"github.com/labstack/echo/v4"
	"github.com/mect/go-escpos"
)

var printMutex = sync.Mutex{}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Mvm Onthaal Printer")
	})

	e.POST("/print", handleMateriaalPrint)
	e.POST("/print-markt", handleMarktPrint)
	e.POST("/eenmaligen", handleEenmaligenPrint)
	e.POST("/sinterklaas", handleSinterklaasPrint)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleMateriaalPrint(c echo.Context) error {
	data := datapkg.MateriaalRequest{}
	c.Bind(&data)

	mainItems := []datapkg.MateriaalItem{}
	seperateItems := []datapkg.MateriaalItem{}
	for _, item := range data.Items {
		if item.SeperateReceipt {
			seperateItems = append(seperateItems, item)
		} else {
			mainItems = append(mainItems, item)
		}
	}

	if len(mainItems) > 0 {
		err := printMateriaalTicket(c, data, mainItems)
		if err != nil {
			return err
		}
	}

	if len(seperateItems) > 0 {
		for _, item := range seperateItems {
			err := printMateriaalTicket(c, data, []datapkg.MateriaalItem{item})
			if err != nil {
				return err
			}
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func printMateriaalTicket(c echo.Context, data datapkg.MateriaalRequest, items []datapkg.MateriaalItem) error {
	printMutex.Lock()

	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinterByPath("") // auto discover USB
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}
	defer p.Close()

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

	//p.Align(escpos.AlignCenter)
	if data.Klant.MVMNummer != "" {
		p.Barcode(strings.Replace(data.Klant.MVMNummer, "MVM", "", -1), escpos.BarcodeTypeCODE39)
		p.PrintLn("")
		p.PrintLn("")
	}

	p.Align(escpos.AlignLeft)
	p.Size(3, 3)
	if data.Klant.MVMNummer != "" {
		p.PrintLn(data.Klant.MVMNummer)
	}
	if data.Klant.EenmaligenNummer != "" {
		p.PrintLn(data.Klant.EenmaligenNummer)
	}
	p.Size(2, 2)
	p.PrintLn(fmt.Sprintf("%s %s", data.Klant.Voornaam, data.Klant.Naam))
	p.PrintLn("")
	p.PrintLn("Materiaal")

	totaal := 0.0

	for _, entry := range items {
		p.Size(1, 1)
		p.PrintLn("==========================================")
		p.Size(2, 2)

		totaal += entry.Prijs

		p.PrintLn(entry.Object)

		hasMaat := false
		if entry.Maat != "" && entry.Maat != "<geen>" {
			hasMaat = true
		}
		if entry.Ontvanger.Naam != "" {
			p.Size(1, 1)
			p.PrintLn(entry.Ontvanger.Naam)
			p.Size(2, 1)
			p.Print(entry.Ontvanger.Geslacht)
			p.Print(" ")
			if !hasMaat {
				p.PrintLn(fmt.Sprintf("%d jaar", entry.Ontvanger.Leeftijd))
			}
		}

		if hasMaat {
			p.PrintLn(fmt.Sprintf("maat: %s", strings.Replace(entry.Maat, " - ", "-", -1)))
		}
		if entry.Opmerking != "" {
			p.PrintLn(entry.Opmerking)
		}

		if entry.ExtraEscposData != "" {
			// base64 decode ExtraEscposData
			base64Decoded, err := base64.StdEncoding.DecodeString(entry.ExtraEscposData)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"status": "error", "error": err.Error()})
			}

			tmpl, err := template.New("print").Parse(string(base64Decoded))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"status": "error", "error": err.Error()})
			}
			escposData := bytes.NewBufferString("")
			err = tmpl.Execute(escposData, data)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"status": "error", "error": err.Error()})
			}

			p.Align(escpos.AlignLeft)
			p.Size(1, 1)
			p.Print(escposData.String())
		}
	}

	p.Size(1, 1)
	p.PrintLn("==========================================")
	p.PrintLn(fmt.Sprintf("\n\nTotaal: %.2f EUR", totaal))

	p.Cut()
	p.End()

	return nil
}

func handleEenmaligenPrint(c echo.Context) error {
	data := datapkg.RequestEenmaligen{}
	c.Bind(&data)

	printMutex.Lock()
	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinterByPath("") // auto discover USB
	defer p.Close()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

	p.Size(4, 4)
	p.PrintLn("VR")

	p.Size(2, 2)
	p.PrintLn(data.EenmaligenNummer)

	p.PrintLn(data.Naam)
	p.PrintLn("")
	p.PrintLn(data.Bericht)

	p.Cut()
	p.End()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func handleSinterklaasPrint(c echo.Context) error {
	data := datapkg.SinterklaasRequest{}
	c.Bind(&data)

	printMutex.Lock()
	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinterByPath("") // auto discover USB
	defer p.Close()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

	// p.Size(3, 3)
	// p.PrintLn(data.Speelgoed.MVMNummer)

	// p.Size(2, 2)
	// p.PrintLn(data.Speelgoed.Naam)
	// p.PrintLn("")

	// p.PrintLn("Sinterklaas")

	// for _, entry := range data.Speelgoed.Paketten {
	// 	p.PrintLn("")
	// 	p.PrintLn("----------------")
	// 	p.PrintLn("")

	// 	p.PrintLn(entry.Naam)
	// 	p.PrintLn(entry.Geslacht)
	// 	p.PrintLn(fmt.Sprintf("%.1f jaar", entry.Leeftijd))
	// 	p.PrintLn(entry.Opmerking)

	// 	p.PrintLn("")
	// 	p.PrintLn("----------------")
	// 	p.PrintLn("")
	// }

	// p.Cut()

	p.Size(4, 4)
	p.PrintLn(fmt.Sprintf("%d", data.Snoep.VolgNummer))

	p.Size(3, 3)
	p.PrintLn("")
	p.PrintLn(data.Snoep.MVMNummer)

	p.Size(2, 2)
	p.PrintLn(data.Snoep.Naam)
	p.PrintLn("")

	p.PrintLn("Sinterklaas Snoep")
	p.PrintLn("")
	p.PrintLn(fmt.Sprintf("%d personen", data.Snoep.Personen))

	p.Cut()

	p.End()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func handleMarktPrint(c echo.Context) error {
	data := datapkg.MarktRequest{}
	c.Bind(&data)

	printMutex.Lock()

	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinterByPath("") // auto discover USB
	defer p.Close()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

	p.Align(escpos.AlignLeft)
	p.Size(2, 2)
	if data.MVMNummer != "" {
		p.PrintLn("MVM" + data.MVMNummer)
	}
	p.Size(2, 2)
	p.PrintLn("")
	p.PrintLn(data.Naam)
	p.Size(2, 1)
	p.PrintLn(data.Beschrijving)
	p.Feed(2)

	for _, entry := range data.Kinderen {
		p.Size(1, 1)
		p.PrintLn("==========================================")
		p.Size(2, 2)

		if entry.Naam != "" {
			p.Size(1, 1)
			p.PrintLn(entry.Naam)
			p.Size(2, 1)
			p.Print(entry.Geslacht)
			p.Print(" ")
			p.PrintLn(fmt.Sprintf("%d jaar", entry.Leeftijd))
		}
	}

	p.Size(1, 1)
	p.PrintLn("==========================================")
	p.AztecViaImage(data.Barcode, 400, 400)

	p.Feed(2)

	p.Cut()
	p.End()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
