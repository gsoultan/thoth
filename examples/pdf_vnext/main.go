package main

import (
	"context"
	"os"

	"github.com/gsoultan/thoth/core"
	"github.com/gsoultan/thoth/document"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	thoth := core.New(logger)
	ctx := context.Background()

	logger.Info().Msg("--- PDF VNext Production Example ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	// 1. Metadata with Unicode and PDF 2.0 readiness
	_ = pdf.SetMetadata(document.Metadata{
		Title:       "Enterprise Unicode & PDF 2.0 Report 🚀",
		Author:      "Lead Architect",
		Subject:     "Advanced PDF features demonstration",
		Keywords:    []string{"Unicode", "PDF 2.0", "Attachments", "Zebra"},
		Description: "A comprehensive report showcasing UTF-8 support, file attachments, and advanced styling.",
	})

	// 2. Attachments (Associated Files - PDF 2.0)
	// We attach the project's go.mod as a source reference
	_ = pdf.AttachFile("go.mod", "project_source.mod", "The project's go module definition file")

	// 3. Advanced Page Settings
	pageSettings := document.PageSettings{
		PaperType:   document.PaperA4,
		Orientation: document.OrientationPortrait,
		Margins: document.Margins{
			Top:    70,
			Bottom: 70,
			Left:   50,
			Right:  50,
		},
		Columns:   2,
		ColumnGap: 20,
	}
	_ = pdf.SetPageSettings(pageSettings)

	// 2. Header and Footer
	_ = pdf.SetHeader("Enterprise Production Report 2026", document.CellStyle{
		Size:       10,
		Color:      "999999",
		Horizontal: "center",
	})

	_ = pdf.SetFooter("Confidential - Page {n} of {nb}", document.CellStyle{
		Size:       10,
		Color:      "999999",
		Horizontal: "right",
	})

	// 3. Table of Contents
	_ = pdf.AddHeading("Executive Summary", 1)
	_ = pdf.AddTableOfContents()

	// 4. Multi-column content with Justified text
	_ = pdf.AddHeading("1. Market Overview", 2)
	lorem := "The global market for document processing is evolving rapidly. Organizations are increasingly looking for modular and extensible libraries that can handle complex layouts with precision. This report outlines the current state of the industry and explores future trends. Our latest technology provides robust support for nested tables, advanced styling, and precise positioning, ensuring production-ready output for the most demanding enterprise use cases."
	_ = pdf.AddParagraph(lorem, document.NewCellStyleBuilder().
		Align("justify", "top").
		LineSpacing(1.2).
		SpacingAfter(15).
		Build())

	_ = pdf.AddParagraph("Key Takeaways:", document.NewCellStyleBuilder().Bold().Build())
	_ = pdf.AddList([]string{
		"High precision coordinate management",
		"Robust cell merging in tables",
		"Recursive nested table support",
		"Dynamic multi-column layouting",
	}, true, document.NewCellStyleBuilder().SpacingAfter(15).Build())

	// 5. Advanced Table: Row/Col Spans, Nesting and Zebra Stripes
	_ = pdf.AddHeading("2. Data Analysis (Advanced Tables)", 2)
	if tbl, err := pdf.AddTable(5, 4); err == nil {
		tbl.SetStyle("zebra:F9F9F9") // Alternating row backgrounds

		headerStyle := document.NewCellStyleBuilder().
			Bold().
			Background("4472C4").
			Color("FFFFFF").
			Padding(5).
			Align("center", "center").
			Border().
			Build()

		_ = tbl.Row(0).Cell(0).AddParagraph("Region").Style(headerStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Q1 Stats").Style(headerStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Q2 Stats").Style(headerStyle)
		_ = tbl.Row(0).Cell(3).AddParagraph("Trend").Style(headerStyle)

		// Region Cell spanning 3 rows
		_ = tbl.Row(1).Cell(0).AddParagraph("North America 🇺🇸").Style(document.NewCellStyleBuilder().
			Align("center", "center").
			Background("D9E1F2").
			Border().
			Build())
		tbl.MergeCells(1, 0, 3, 1)

		_ = tbl.Row(1).Cell(1).AddParagraph("1,240").Style(document.NewCellStyleBuilder().Border().Align("right", "center").Build())
		_ = tbl.Row(1).Cell(2).AddParagraph("1,350").Style(document.NewCellStyleBuilder().Border().Align("right", "center").Build())

		// Nested table in a cell
		nestedCell := tbl.Row(1).Cell(3)
		nTbl := nestedCell.AddTable(2, 2)
		nStyle := document.NewCellStyleBuilder().Size(8).Border().Build()
		_ = nTbl.Row(0).Cell(0).AddParagraph("ID").Style(nStyle)
		_ = nTbl.Row(0).Cell(1).AddParagraph("Value").Style(nStyle)
		_ = nTbl.Row(1).Cell(0).AddParagraph("A-1").Style(nStyle)
		_ = nTbl.Row(1).Cell(1).AddParagraph("98%").Style(nStyle)

		_ = tbl.Row(2).Cell(1).AddParagraph("2,100").Style(document.NewCellStyleBuilder().Border().Align("right", "center").Build())
		_ = tbl.Row(2).Cell(2).AddParagraph("2,250").Style(document.NewCellStyleBuilder().Border().Align("right", "center").Build())
		_ = tbl.Row(2).Cell(3).AddParagraph("Upward 📈").Style(document.NewCellStyleBuilder().Border().Size(9).Build())

		_ = tbl.Row(3).Cell(1).AddParagraph("3,100").Style(document.NewCellStyleBuilder().Border().Align("right", "center").Build())
		_ = tbl.Row(3).Cell(2).AddParagraph("3,450").Style(document.NewCellStyleBuilder().Border().Align("right", "center").Build())
		_ = tbl.Row(3).Cell(3).AddParagraph("Steady").Style(document.NewCellStyleBuilder().Border().Size(9).Build())

		_ = tbl.Row(4).Cell(1).AddParagraph("Sum").Style(document.NewCellStyleBuilder().Bold().Border().Build())
		_ = tbl.Row(4).Cell(2).AddParagraph("7,050").Style(document.NewCellStyleBuilder().Bold().Border().Align("right", "center").Build())
		tbl.MergeCells(4, 2, 1, 2)
	}

	// 6. Final Page with specialized layout
	_ = pdf.AddPageBreak()
	// Update page settings mid-document
	_ = pdf.SetPageSettings(document.PageSettings{
		PaperType: document.PaperA4,
		Columns:   1,
		Margins:   document.Margins{Top: 50, Bottom: 50, Left: 50, Right: 50},
	})

	_ = pdf.AddHeading("3. Conclusion", 1)
	_ = pdf.AddParagraph("The results demonstrate that our PDF engine is ready for production workloads. The combination of advanced layout precision and robust table handling makes it an ideal choice for complex enterprise reports.",
		document.NewCellStyleBuilder().SpacingAfter(20).Build())

	// 7. Absolute Positioning: Confidential Stamp
	stampStyle := document.NewCellStyleBuilder().
		Pos(400, 800).
		Size(30).
		Bold().
		Color("FF0000").
		Opacity(0.3).
		Build()
	_ = pdf.AddParagraph("CONFIDENTIAL", stampStyle)

	if err := pdf.Export("vnext_production.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved vnext_production.pdf")
	}
}
