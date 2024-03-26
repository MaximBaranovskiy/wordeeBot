package pdf

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
	"strings"
	"wordeeBot/internal/model/db"
)

func MakeDictionaryPDF(dictionaryName string, words []db.Word, fields map[string]bool) (error, *tgbotapi.FileBytes) {
	colNum := calculateColNums(&fields)
	colNames := createColNames(&fields)
	pdf := createPDF(&colNum)
	pageWidth, _ := pdf.GetPageSize()

	pdf.Text(pageWidth/2-pdf.GetStringWidth(dictionaryName)/2, 20, strings.ToUpper(dictionaryName))

	colWidth := make([]float64, colNum)
	for i := 0; i < colNum; i++ {
		colWidth[i] = pageWidth / float64(colNum)
	}

	pdf.SetY(40)
	pdf.SetX(0)

	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	for i := 0; i < colNum; i++ {
		pdf.CellFormat(colWidth[i], 16, colNames[i], "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetX(0)

	for _, word := range words {

		writingLines := pdf.SplitText(word.Writing, colWidth[0])
		writingInd := 0

		translationLines := pdf.SplitText(word.Translation, colWidth[0])
		translationInd := 0

		transcriptionLines := pdf.SplitText(word.Transcription, colWidth[0])
		transcriptionInd := 0

		synonymsLines := pdf.SplitText(word.Synonyms, colWidth[0])
		synonymsInd := 0

		antonymsLines := pdf.SplitText(word.Antonyms, colWidth[0])
		antonymsInd := 0

		definitionLines := pdf.SplitText(word.Definition, colWidth[0])
		definitionInd := 0

		collocationsLines := pdf.SplitText(word.Collocations, colWidth[0])
		collocationsInd := 0

		idiomsLines := pdf.SplitText(word.Idioms, colWidth[0])
		idiomsInd := 0

		maxLen := max(max(max(len(writingLines), len(transcriptionLines)), max(len(translationLines), len(synonymsLines))),
			max(max(len(antonymsLines), len(definitionLines)), max(len(collocationsLines), len(idiomsLines))))

		count := 0
		for count != maxLen {

			for i := 0; i < colNum; i++ {
				cellText := ""
				switch colNames[i] {
				case "Слово":
					if writingInd < len(writingLines) {
						cellText = writingLines[writingInd]
						writingInd++
					}
				case "Перевод":
					if translationInd < len(translationLines) {
						cellText = translationLines[translationInd]
						translationInd++
					}
				case "Транскрипция":
					if transcriptionInd < len(transcriptionLines) {
						cellText = transcriptionLines[transcriptionInd]
						transcriptionInd++
					}
				case "Синонимы":
					if synonymsInd < len(synonymsLines) {
						cellText = synonymsLines[synonymsInd]
						synonymsInd++
					}
				case "Антонимы":
					if antonymsInd < len(antonymsLines) {
						cellText = antonymsLines[antonymsInd]
						antonymsInd++
					}
				case "Определение":
					if definitionInd < len(definitionLines) {
						cellText = definitionLines[definitionInd]
						definitionInd++
					}
				case "Коллокации":
					if collocationsInd < len(collocationsLines) {
						cellText = collocationsLines[collocationsInd]
						collocationsInd++
					}
				case "Идиомы":
					if idiomsInd < len(idiomsLines) {
						cellText = idiomsLines[idiomsInd]
						idiomsInd++
					}
				}
				borderStr := "LR"
				if count == 0 {
					borderStr = "LRT"
				}
				pdf.CellFormat(colWidth[i], 16, cellText, borderStr, 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetX(0)

			count++
		}

	}

	err := pdf.OutputFileAndClose(dictionaryName + ".pdf")
	if err != nil {
		return err, nil
	}

	fileBytes, err := ioutil.ReadFile(dictionaryName + ".pdf")
	if err != nil {
		return err, nil
	}

	return nil, &tgbotapi.FileBytes{Name: dictionaryName + ".pdf", Bytes: fileBytes}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func calculateColNums(fields *map[string]bool) int {
	colNum := 0

	for _, isField := range *fields {
		if isField {
			colNum++
		}
	}

	return colNum
}

func createColNames(fields *map[string]bool) []string {
	colNames := make([]string, 0)

	colNames = append(colNames, "Слово")
	for _, field := range db.DbWordFields {
		if (*fields)[field] {
			colNames = append(colNames, field)
		}
	}

	return colNames
}

func createPDF(colNum *int) *gofpdf.Fpdf {
	orientation := gofpdf.OrientationPortrait
	if *colNum > 4 {
		orientation = gofpdf.OrientationLandscape
	}

	sizeStr := "A4"
	if *colNum > 4 && *colNum <= 6 {
		sizeStr = "A3"
	} else if *colNum > 6 {
		sizeStr = "A2"
	}

	pdf := gofpdf.New(orientation, "mm", sizeStr, "")
	pdf.AddUTF8Font("DejaVu", "", "DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 16)
	pdf.AddPage()

	return pdf
}
