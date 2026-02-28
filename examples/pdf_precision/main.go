package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gsoultan/thoth/core"
	"github.com/gsoultan/thoth/document"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	thoth := core.New(logger)
	ctx := context.Background()

	logger.Info().Msg("--- PDF Precision Styling & Auto-Sizing Example ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	_ = pdf.SetPageSettings(document.NewPageSettingsBuilder().WithPaperType(document.PaperA4).Build())

	// 1. Table Auto-Sizing
	_ = pdf.AddParagraph("1. Table Auto-Sizing (No widths provided)", document.NewCellStyleBuilder().Bold().Size(14).SpacingAfter(10).Build())
	if tbl, err := pdf.AddTable(3, 3); err == nil {
		hStyle := document.NewCellStyleBuilder().Bold().Background("EEEEEE").Border().Build()
		_ = tbl.Row(0).Cell(0).AddParagraph("ID").Style(hStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Short Name").Style(hStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("A Very Very Long Header Title to test wrapping and auto-sizing").Style(hStyle)

		dStyle := document.NewCellStyleBuilder().Border().Build()
		_ = tbl.Row(1).Cell(0).AddParagraph("1").Style(dStyle)
		_ = tbl.Row(1).Cell(1).AddParagraph("Alpha").Style(dStyle)
		_ = tbl.Row(1).Cell(2).AddParagraph("Content A").Style(dStyle)

		_ = tbl.Row(2).Cell(0).AddParagraph("2").Style(dStyle)
		_ = tbl.Row(2).Cell(1).AddParagraph("Beta Project").Style(dStyle)
		_ = tbl.Row(2).Cell(2).AddParagraph("Content B with some more text").Style(dStyle)
	}

	// 2. Precise Borders
	_ = pdf.AddParagraph("\n2. Precise Borders (Top/Bottom only, custom width/color)", document.NewCellStyleBuilder().Bold().Size(14).SpacingBefore(20).SpacingAfter(10).Build())
	if tbl, err := pdf.AddTable(3, 3); err == nil {
		// Only bottom border for header
		hStyle := document.NewCellStyleBuilder().Bold().BorderBottom().BorderWidth(2).BorderColor("4F81BD").Build()
		_ = tbl.Row(0).Cell(0).AddParagraph("Category").Style(hStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Description").Style(hStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Amount").Style(hStyle)

		// Grid with light horizontal lines only
		lineStyle := document.NewCellStyleBuilder().BorderBottom().BorderWidth(0.5).BorderColor("CCCCCC").Build()
		for r := 1; r < 3; r++ {
			_ = tbl.Row(r).Cell(0).AddParagraph("Exp").Style(lineStyle)
			_ = tbl.Row(r).Cell(1).AddParagraph("Travel costs").Style(lineStyle)
			_ = tbl.Row(r).Cell(2).AddParagraph("$100.00").Style(lineStyle)
		}
	}

	// 3. Multi-row Repeating Headers
	_ = pdf.AddParagraph("\n3. Multi-row Repeating Headers (spans 2 pages)", document.NewCellStyleBuilder().Bold().Size(14).SpacingBefore(20).SpacingAfter(10).Build())
	if tbl, err := pdf.AddTable(50, 3); err == nil {
		_ = tbl.SetHeaderRows(2) // Repeat first 2 rows

		hStyle := document.NewCellStyleBuilder().Bold().Background("4F81BD").Color("FFFFFF").Border().Align("center", "center").Build()
		_ = tbl.MergeCells(0, 0, 1, 3)
		_ = tbl.Row(0).Cell(0).AddParagraph("Company Financial Report 2026").Style(hStyle)

		hStyle2 := document.NewCellStyleBuilder().Bold().Background("D9D9D9").Border().Build()
		_ = tbl.Row(1).Cell(0).AddParagraph("Q1").Style(hStyle2)
		_ = tbl.Row(1).Cell(1).AddParagraph("Q2").Style(hStyle2)
		_ = tbl.Row(1).Cell(2).AddParagraph("Total").Style(hStyle2)

		dStyle := document.NewCellStyleBuilder().Border().Build()
		for r := 2; r < 50; r++ {
			_ = tbl.Row(r).Cell(0).AddParagraph(fmt.Sprintf("%d00", r)).Style(dStyle)
			_ = tbl.Row(r).Cell(1).AddParagraph(fmt.Sprintf("%d50", r)).Style(dStyle)
			_ = tbl.Row(r).Cell(2).AddParagraph(fmt.Sprintf("%d50", r*2)).Style(dStyle)
		}
	}

	if err := pdf.Export("precision_styling.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved precision_styling.pdf")
	}
}
