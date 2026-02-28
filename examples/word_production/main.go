package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word"
)

func main() {
	ctx := context.Background()

	// Create a new Word document
	doc := word.NewDocument().(document.WordProcessor)
	defer doc.Close()

	// 1. Production Metadata
	doc.SetMetadata(document.Metadata{
		Title:   "Production Quarterly Analysis",
		Author:  "Thoth Enterprise Engine",
		Subject: "Business Performance Report",
		Company: "Thoth Solutions Ltd.",
	})

	// 2. Production Header & Footer with Page Numbering
	headerStyle := document.CellStyle{Size: 9, Color: "808080", Italic: true, Horizontal: "right"}
	doc.SetHeader("CONFIDENTIAL - Q1 2026", headerStyle)

	footerStyle := document.CellStyle{Size: 10, Color: "444444", Horizontal: "center"}
	doc.SetFooter("Page {n} of {nb}", footerStyle)

	// 3. Title Page
	doc.AddHeading("Quarterly Business Analysis", 1, document.CellStyle{
		Bold: true, Size: 28, Color: "2E5481", Horizontal: "center", SpacingAfter: 20,
	})
	doc.AddParagraph("Prepared for the Executive Board", document.CellStyle{
		Size: 14, Italic: true, Horizontal: "center", SpacingAfter: 40,
	})

	doc.AddPageBreak()

	// 4. Table of Contents
	doc.AddHeading("Table of Contents", 1)
	_ = doc.AddTableOfContents()
	doc.AddPageBreak()

	// 5. Executive Summary
	doc.AddHeading("Executive Summary", 1)
	doc.AddParagraph("This report provides a detailed analysis of the performance metrics for Q1 2026. " +
		"We have seen a significant growth in the cloud services sector, while traditional licensing " +
		"remains steady.")

	// 6. Detailed Performance Table (Production Grade)
	doc.AddHeading("Performance Metrics", 2)
	table, _ := doc.AddTable(5, 4)
	table.SetStyle("TableGrid")
	table.SetColumnWidths(1.5, 1.2, 1.2, 1.2)
	table.SetHeaderRows(1)

	headerCell := document.CellStyle{
		Bold: true, Background: "2E5481", Color: "FFFFFF", Horizontal: "center", Vertical: "center", Padding: 5,
	}
	bodyCell := document.CellStyle{Padding: 5, Vertical: "center"}
	numCell := document.CellStyle{Padding: 5, Vertical: "center", Horizontal: "right"}

	table.Row(0).Cell(0).AddParagraph("Region", headerCell)
	table.Row(0).Cell(1).AddParagraph("Revenue (M)", headerCell)
	table.Row(0).Cell(2).AddParagraph("Growth (%)", headerCell)
	table.Row(0).Cell(3).AddParagraph("Status", headerCell)

	data := [][]string{
		{"North America", "45.2", "+12.5", "Exceeding"},
		{"Europe", "38.7", "+8.2", "Target"},
		{"Asia Pacific", "29.4", "+18.9", "Exceeding"},
		{"Latin America", "12.1", "-2.4", "Under"},
	}

	for i, rowData := range data {
		r := table.Row(i + 1)
		r.Cell(0).AddParagraph(rowData[0], bodyCell)
		r.Cell(1).AddParagraph(rowData[1], numCell)
		r.Cell(2).AddParagraph(rowData[2], numCell)

		statusStyle := bodyCell
		if rowData[3] == "Under" {
			statusStyle.Color = "FF0000"
			statusStyle.Bold = true
		}
		r.Cell(3).AddParagraph(rowData[3], statusStyle)
	}

	// 7. Advanced Lists (Real Numbering)
	doc.AddHeading("Strategic Initiatives", 2)
	doc.AddList([]string{
		"Expansion into South East Asian markets",
		"Consolidation of data center infrastructure",
		"Implementation of AI-driven customer support",
	}, true) // Numbered list

	doc.AddHeading("Risk Assessment", 2)
	doc.AddList([]string{
		"Currency fluctuations in emerging markets",
		"Cybersecurity threats on legacy systems",
		"Talent acquisition in specialized AI roles",
	}, false) // Bulleted list

	// Save the document
	f, err := os.Create("production_report.docx")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()

	if err := doc.Save(ctx, f); err != nil {
		fmt.Printf("Error saving document: %v\n", err)
		return
	}

	fmt.Println("Production report generated: production_report.docx")
}
