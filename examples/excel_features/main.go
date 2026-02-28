package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel"
)

func main() {
	ctx := context.Background()
	doc := excel.NewDocument().(document.Spreadsheet)
	defer doc.Close()

	sheet, _ := doc.Sheet("Features")

	// 1. Rich Text
	richText := []document.TextSpan{
		{Text: "Normal "},
		{Text: "Bold ", Style: document.CellStyle{Bold: true}},
		{Text: "Red ", Style: document.CellStyle{Color: "FF0000"}},
		{Text: "Big", Style: document.CellStyle{Size: 20}},
	}
	sheet.Cell("A1").Set(richText)
	sheet.SetRowHeight(1, 30)
	sheet.SetColumnWidth(1, 30)

	// 2. Data Validation (Dropdown)
	sheet.Cell("B1").Set("Pick one:")
	sheet.SetDataValidation("B2", "Option 1", "Option 2", "Option 3")
	sheet.Cell("B2").Set("Option 1")

	// 3. Conditional Formatting
	// Highlight cells > 100 with green background
	greenStyle := document.CellStyle{Background: "00FF00", Bold: true}
	sheet.SetConditionalFormatting("C1:C10", greenStyle)

	sheet.Cell("C1").Set(50)
	sheet.Cell("C2").Set(150)
	sheet.Cell("C3").Set(200)

	// Save
	f, err := os.Create("excel_features.xlsx")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer f.Close()

	if err := doc.Save(ctx, f); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		return
	}

	fmt.Println("Excel features report generated: excel_features.xlsx")
}
