package main

import (
	"fmt"
	"os"

	"github.com/gosimple/slug"
	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run main.go base.pdf \"text\"...\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	for _, textStr := range os.Args[2:] {
		outputPath := slug.Make(textStr) + ".pdf"
		err := addTextToPdf(inputPath, outputPath, textStr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Complete, see output file: %s\n", outputPath)
	}
}

const defaultFontSize = float64(16)
const cellWidth = 149
const paddingX = 15

const xCount = 4
const yCount = 7
const fontFile = "font.ttf"

func addTextToPdf(inputPath string, outputPath string, text string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	c := creator.New()
	// Create a mock to get real Height value
	mockC := creator.New()

	font, err := pdf.NewPdfFontFromTTFFile(fontFile)
	p := creator.NewParagraph(text)
	if err == nil {
		p.SetFont(font)
	}
	p.SetTextAlignment(creator.TextAlignmentCenter)
	fontSize := defaultFontSize
	p.SetWidth(cellWidth - 2*paddingX)
	p.SetFontSize(fontSize)
	p.SetPos(10, 10)
	_ = mockC.Draw(p)
	if p.Height() >= 2*defaultFontSize {
		fontSize = defaultFontSize - 4
		p.SetFontSize(fontSize)
		_ = mockC.Draw(p)
		if p.Height() >= 2*(fontSize) {
			fontSize = fontSize - 2
			p.SetFontSize(fontSize)
		}
	}
	height := p.Height()
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		err = c.AddPage(page)
		if err != nil {
			return err
		}

		for x := 0; x < xCount; x++ {
			for y := 0; y < yCount; y++ {
				xPos := float64(x)*cellWidth + paddingX
				yPos := float64(25 + defaultFontSize - height + float64(y)*112)
				p.SetPos(xPos, yPos)

				err := c.Draw(p)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
			}
		}

	}

	err = c.WriteToFile(outputPath)
	return err
}
