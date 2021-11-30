package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4/middleware"

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
	e.POST("/eenmaligen", handleEenmaligenPrint)
	e.POST("/sinterklaas", handleSinterklaasPrint)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleMateriaalPrint(c echo.Context) error {
	data := MateriaalRequest{}
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

	//p.Align(escpos.AlignCenter)
	p.Barcode(strings.Replace(data.Klant.MVMNummer, "MVM", "", -1), escpos.BarcodeTypeCODE39)
	p.PrintLn("")
	p.PrintLn("")

	p.Align(escpos.AlignLeft)
	p.Size(3, 3)
	p.PrintLn(data.Klant.MVMNummer)
	p.Size(2, 2)
	p.PrintLn(fmt.Sprintf("%s %s", data.Klant.Voornaam, data.Klant.Naam))
	p.PrintLn("")
	p.PrintLn("Materiaal")

	totaal := 0.0

	for _, entry := range data.Items {
		p.PrintLn("")
		p.PrintLn("----------------")
		p.PrintLn("")

		totaal += entry.Prijs

		p.PrintLn(entry.Object)

		hasMaat := false
		if entry.Maat != "" && entry.Maat != "<geen>" {
			hasMaat = true
		}
		if entry.Ontvanger.Naam != "" {
			p.PrintLn(entry.Ontvanger.Naam)
			p.PrintLn(entry.Ontvanger.Geslacht)
			if !hasMaat {
				p.PrintLn(fmt.Sprintf("%d jaar", entry.Ontvanger.Leeftijd))
			}
		}

		if hasMaat {
			p.PrintLn(fmt.Sprintf("Maat: %s", entry.Maat))
		}

		p.PrintLn(entry.Opmerking)
	}

	p.PrintLn("")
	p.PrintLn("----------------")
	p.PrintLn("")

	p.PrintLn(fmt.Sprintf("\n\nTotaal: %.2f EUR", totaal))

	p.Cut()
	p.End()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func handleEenmaligenPrint(c echo.Context) error {
	data := RequestEenmaligen{}
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
	data := SinterklaasRequest{}
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

	p.Size(3, 3)
	p.PrintLn(data.Speelgoed.MVMNummer)

	log.Printf("Speelgoed voor %s\n", data.Speelgoed.MVMNummer)

	p.Size(2, 2)
	p.PrintLn(data.Speelgoed.Naam)
	p.PrintLn("")

	p.PrintLn("Sinterklaas")

	for _, entry := range data.Speelgoed.Paketten {
		p.PrintLn("")
		p.PrintLn("----------------")
		p.PrintLn("")

		p.PrintLn(entry.Naam)
		p.PrintLn(entry.Geslacht)
		p.PrintLn(fmt.Sprintf("%.1f jaar", entry.Leeftijd))
		p.PrintLn(entry.Opmerking)

		p.PrintLn("")
		p.PrintLn("----------------")
		p.PrintLn("")

		log.Printf("%s is braaf geweest\n", entry.Naam)
	}

	p.Cut()

	// p.Size(3, 3)
	// p.PrintLn(data.Snoep.MVMNummer)

	// p.Size(2, 2)
	// p.PrintLn(data.Snoep.Naam)
	// p.PrintLn("")

	// p.PrintLn("Sinterklaas Snoep")

	// p.PrintLn(fmt.Sprintf("volwassenen: %d", data.Snoep.Volwassenen))
	// p.PrintLn(fmt.Sprintf("kinderen: %d", data.Snoep.Kinderen))

	// p.Cut()

	// p.End()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
