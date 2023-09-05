package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	file, _ := reader.ReadString('\n')
	pdfFileName := strings.TrimSuffix(file, "\n")
	reader = bufio.NewReader(os.Stdin)
	format, _ := reader.ReadString('\n')
	format = strings.TrimSuffix(format, "\n")

	pageNumber := 1

	ctx, err := pdfcpu.ReadFile(pdfFileName, nil)

	if err != nil {
		fmt.Printf("Fehler beim Einlesen der PDF-Datei: %v\n", err)
		return
	}

	extractedCtx, err := pdfcpu.ExtractPage(ctx, pageNumber)
	if err != nil {
		fmt.Printf("Fehler beim Extrahieren der Seite: %v\n", err)
		return
	}

	outputFileName := fmt.Sprintf("extrahierte-seite-%d.pdf", pageNumber)

	err = api.WriteContextFile(extractedCtx, outputFileName)
	if err != nil {
		fmt.Printf("Fehler beim Speichern der extrahierten Seite: %v\n", err)
		return
	}

	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err = mw.ReadImage("./extrahierte-seite-1.pdf")
	if err != nil {
		fmt.Println(err)
	}
	mw.SetIteratorIndex(0)
	err = mw.SetImageFormat(format)
	if err != nil {
		fmt.Println(err)
	}
	err = mw.WriteImage("test" + "." + format)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf(
		"Seite %d wurde erfolgreich extrahiert und als %s gespeichert.\n",
		pageNumber,
		outputFileName,
	)
}
