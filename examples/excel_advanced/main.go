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

	sheet, _ := doc.Sheet("Settings")

	// 1. Set up some options for data validation
	sheet.Cell("A1").Set("Active")
	sheet.Cell("A2").Set("Inactive")
	sheet.Cell("A3").Set("Pending")

	// 2. Named Range for these options
	_ = doc.SetNamedRange("StatusList", "Settings!$A$1:$A$3")

	dataSheet, _ := doc.Sheet("Data")
	dataSheet.Cell("A1").Set("User")
	dataSheet.Cell("B1").Set("Status")

	headerStyle := document.CellStyle{Bold: true, Background: "CCCCCC", Border: true}
	dataSheet.Cell("A1").Style(headerStyle)
	dataSheet.Cell("B1").Style(headerStyle)

	// 3. Data Validation using the named range or literal list
	// Using literal list for now as per our implementation
	_ = dataSheet.SetDataValidation("B2:B10", "Active", "Inactive", "Pending")

	dataSheet.Cell("A2").Set("John Doe")
	dataSheet.Cell("B2").Set("Active")

	dataSheet.Cell("A3").Set("Jane Smith")
	dataSheet.Cell("B3").Set("Pending")

	// 4. Advanced Styling
	fancyStyle := document.CellStyle{
		Bold:        true,
		Italic:      true,
		Size:        14,
		Color:       "FFFFFF",
		Background:  "4472C4",
		Border:      true,
		BorderColor: "000000",
		BorderWidth: 2,
		Horizontal:  "center",
		Vertical:    "center",
	}
	dataSheet.Cell("D1").Set("ADVANCED REPORT").Style(fancyStyle)
	_ = dataSheet.MergeCells("D1:F1")

	// Save
	outputPath := "advanced_excel.xlsx"
	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()

	if err := doc.Save(ctx, f); err != nil {
		fmt.Printf("Error saving document: %v\n", err)
		return
	}

	fmt.Printf("Advanced Excel report generated: %s\n", outputPath)
}
