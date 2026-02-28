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

	// Create a new Excel document
	doc := excel.NewDocument().(document.Spreadsheet)
	defer doc.Close()

	// 1. Workbook Protection
	_ = doc.SetPassword("secret123")

	sheet, _ := doc.Sheet("Strategic Report")

	// 2. Page Settings
	sheet.SetPageSettings(document.PageSettings{
		PaperType:   document.PaperA4,
		Orientation: document.OrientationLandscape,
		Margins: document.Margins{
			Top:    72, // 1 inch
			Bottom: 72,
			Left:   54, // 0.75 inch
			Right:  54,
		},
	})

	// 3. Header & Footer
	sheet.SetHeader("&LThoth Enterprise&CConfidential&R&D")
	sheet.SetFooter("&LPage &P of &N&RStrategic Planning")

	// 4. Content with Grouping (Outlining)
	headerStyle := document.CellStyle{Bold: true, Background: "D9E1F2", Border: true}

	// Group 1: Financials
	sheet.Cell("A1").Set("1. Financial Summary").Style(headerStyle)
	sheet.Cell("A2").Set("Q1 Revenue")
	sheet.Cell("B2").Set(1250000.0)
	sheet.Cell("A3").Set("Q1 Expenses")
	sheet.Cell("B3").Set(850000.0)
	sheet.Cell("A4").Set("Q1 Profit")
	sheet.Cell("B4").Formula("B2-B3").Style(document.CellStyle{Bold: true})

	// Group rows 2-3 under row 1
	sheet.GroupRows(2, 3, 1)

	// Group 2: Projections
	sheet.Cell("A6").Set("2. 2026 Projections").Style(headerStyle)
	sheet.Cell("A7").Set("Market Growth")
	sheet.Cell("B7").Set(0.15).Style(document.CellStyle{NumberFormat: "0%"})
	sheet.Cell("A8").Set("Target Revenue")
	sheet.Cell("B8").Formula("B2*(1+B7)")

	sheet.GroupRows(7, 8, 1)

	// 5. Column Grouping
	sheet.Cell("C1").Set("Details").Style(headerStyle)
	sheet.Cell("C2").Set("Verified")
	sheet.Cell("C3").Set("Pending")
	sheet.GroupCols(3, 3, 1) // Group Column C

	// 6. Sheet Protection
	sheet.Protect("sheetpass")

	// Formatting
	sheet.SetColumnWidth(1, 20)
	sheet.SetColumnWidth(2, 15)

	// Save
	f, err := os.Create("advanced_production_report.xlsx")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()

	if err := doc.Save(ctx, f); err != nil {
		fmt.Printf("Error saving document: %v\n", err)
		return
	}

	fmt.Println("Advanced Excel production report generated: advanced_production_report.xlsx")
}
