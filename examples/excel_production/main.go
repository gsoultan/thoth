package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel"
)

func main() {
	ctx := context.Background()

	// Create a new Excel document
	doc := excel.NewDocument().(document.Spreadsheet)
	defer doc.Close()

	// 1. Metadata
	doc.SetMetadata(document.Metadata{
		Title:       "Production Sales Report",
		Author:      "Thoth Engine",
		Subject:     "Quarterly Performance",
		Description: "Comprehensive sales data with formulas and formatting.",
	})

	sheet, _ := doc.Sheet("Sheet1")

	// 2. Header Style
	headerStyle := document.CellStyle{
		Bold:       true,
		Background: "D9D9D9",
		Horizontal: "center",
		Border:     true,
	}

	// Set Headers
	sheet.Cell("A1").Set("ID").Style(headerStyle)
	sheet.Cell("B1").Set("Product").Style(headerStyle)
	sheet.Cell("C1").Set("Quantity").Style(headerStyle)
	sheet.Cell("D1").Set("Unit Price").Style(headerStyle)
	sheet.Cell("E1").Set("Total").Style(headerStyle)
	sheet.Cell("F1").Set("Date").Style(headerStyle)

	// 3. Data Rows
	data := []struct {
		ID    int
		Prod  string
		Qty   int
		Price float64
		Date  time.Time
	}{
		{101, "Enterprise License", 2, 1250.00, time.Now()},
		{102, "Support Contract", 1, 499.00, time.Now().AddDate(0, 0, -1)},
		{103, "Custom Module", 3, 250.00, time.Now().AddDate(0, 0, -5)},
	}

	dateStyle := document.CellStyle{NumberFormat: "mm-dd-yy", Horizontal: "right"}
	moneyStyle := document.CellStyle{NumberFormat: "#,##0.00", Horizontal: "right"}

	for i, item := range data {
		row := i + 2
		sheet.Cell(fmt.Sprintf("A%d", row)).Set(item.ID)
		sheet.Cell(fmt.Sprintf("B%d", row)).Set(item.Prod)
		sheet.Cell(fmt.Sprintf("C%d", row)).Set(item.Qty)
		sheet.Cell(fmt.Sprintf("D%d", row)).Set(item.Price).Style(moneyStyle)

		// 4. Formula: Quantity * Unit Price
		sheet.Cell(fmt.Sprintf("E%d", row)).Formula(fmt.Sprintf("C%d*D%d", row, row)).Style(moneyStyle)

		sheet.Cell(fmt.Sprintf("F%d", row)).Set(item.Date).Style(dateStyle)
	}

	// 5. Total Row
	lastRow := len(data) + 2
	sheet.Cell(fmt.Sprintf("D%d", lastRow)).Set("GRAND TOTAL:").Style(document.CellStyle{Bold: true, Horizontal: "right"})
	sheet.Cell(fmt.Sprintf("E%d", lastRow)).Formula(fmt.Sprintf("SUM(E2:E%d)", lastRow-1)).Style(document.CellStyle{Bold: true, Background: "FFFF00", NumberFormat: "#,##0.00", Border: true})

	// 6. Formatting & Utilities
	sheet.SetColumnWidth(1, 10)
	sheet.SetColumnWidth(2, 25)
	sheet.SetColumnWidth(3, 12)
	sheet.SetColumnWidth(4, 15)
	sheet.SetColumnWidth(5, 15)
	sheet.SetColumnWidth(6, 15)

	// Set header row height
	sheet.SetRowHeight(1, 25)

	sheet.AutoFilter("A1:F1")
	sheet.FreezePanes(0, 1)

	// Merge cells for a title at the bottom
	sheet.MergeCells("B8:E8")
	sheet.Cell("B8").Set("CONFIDENTIAL SALES REPORT").Style(document.CellStyle{
		Bold:       true,
		Size:       14,
		Horizontal: "center",
		Background: "FFCCCC",
	})

	// 7. Hyperlink & Image
	sheet.Cell("B10").Set("Documentation").Hyperlink("https://github.com/gsoultan/thoth").Style(document.CellStyle{Color: "0000FF"})

	// Insert Image (using transparent.png from project root)
	if _, err := os.Stat("../../transparent.png"); err == nil {
		sheet.InsertImage("../../transparent.png", 7, 10) // Column 7 (H), Row 10
	} else if _, err := os.Stat("transparent.png"); err == nil {
		sheet.InsertImage("transparent.png", 7, 10)
	}

	// Save
	f, err := os.Create("production_sales.xlsx")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()

	if err := doc.Save(ctx, f); err != nil {
		fmt.Printf("Error saving document: %v\n", err)
		return
	}

	fmt.Println("Excel production report generated: production_sales.xlsx")
}
