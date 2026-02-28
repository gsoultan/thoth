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

	logger.Info().Msg("--- PDF Enterprise v2 Production Example ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	// 1. Metadata with advanced fields
	_ = pdf.SetMetadata(document.Metadata{
		Title:       "Enterprise Q1 2026 Strategy Report",
		Author:      "Jane Doe, Lead Architect",
		Subject:     "Strategic Planning & PDF Module Capabilities",
		Keywords:    []string{"Strategy", "Enterprise", "PDF", "Accessibility"},
		Description: "A comprehensive report demonstrating the advanced features of the Thoth PDF engine, including accessibility and vector graphics.",
	})

	// 2. Page Settings
	_ = pdf.SetPageSettings(document.PageSettings{
		PaperType:   document.PaperA4,
		Orientation: document.OrientationPortrait,
		Margins: document.Margins{
			Top:    60,
			Bottom: 60,
			Left:   50,
			Right:  50,
		},
	})

	// 3. Tagged Content (Headings & Alt Text)
	_ = pdf.AddHeading("Enterprise Strategy Report 2026", 1)
	_ = pdf.AddParagraph("This document is generated with full accessibility support (Tagged PDF). All headings and paragraphs are properly structured for screen readers.",
		document.NewCellStyleBuilder().Italic().SpacingAfter(15).Build())

	// 4. Advanced Vector Graphics
	_ = pdf.AddHeading("1. Financial Ecosystem Visualization", 2)

	// Draw a styled diagram with ellipses and dashed lines
	// Central node
	_ = pdf.DrawEllipse(250, 620, 100, 60, document.NewCellStyleBuilder().
		Background("D9E1F2").
		Border().
		BorderWidth(2).
		BorderColor("4472C4").
		Build())
	_ = pdf.AddParagraph("Core Growth Hub", document.NewCellStyleBuilder().Align("center", "top").Indent(0).SpacingAfter(0).Build()) // This positioning is currently global, but for diagrams we'd need better absolute text positioning.
	// In the current processor, AddParagraph appends. Diagrams usually need absolute pos.
	// Thoth's PDF processor is mostly flow-based, but Draw methods are absolute.

	// Connecting dashed lines
	dashedStyle := document.NewCellStyleBuilder().
		Color("666666").
		BorderWidth(1.5).
		Build()
	dashedStyle.DashPattern = []float64{5, 5}

	_ = pdf.DrawLine(300, 620, 150, 520, dashedStyle)
	_ = pdf.DrawLine(300, 620, 450, 520, dashedStyle)

	// Sub-nodes
	_ = pdf.DrawEllipse(100, 480, 80, 50, document.NewCellStyleBuilder().Background("E2EFDA").Border().Build())
	_ = pdf.DrawEllipse(420, 480, 80, 50, document.NewCellStyleBuilder().Background("FFF2CC").Border().Build())

	// 5. Image with Alt Text
	_ = pdf.AddHeading("2. Visual Assets", 2)
	// Assuming transparent.png exists from previous sessions
	if _, err := os.Stat("transparent.png"); err == nil {
		_ = pdf.InsertImage("transparent.png", 200, 150, document.NewCellStyleBuilder().
			Alt("Enterprise Logo with transparent background for high-fidelity rendering").
			SpacingAfter(10).
			Build())
	}

	// 6. Accessible Lists
	_ = pdf.AddHeading("3. Strategic Objectives", 2)
	_ = pdf.AddList([]string{
		"Global infrastructure expansion with high availability",
		"Enhanced document processing with Thoth Enterprise Edition",
		"Full compliance with PDF/UA accessibility standards",
	}, true, document.NewCellStyleBuilder().SpacingAfter(20).Build())

	// 7. Robust Table with specialized styling
	_ = pdf.AddHeading("4. Market Performance Data", 2)
	if tbl, err := pdf.AddTable(3, 3); err == nil {
		_ = tbl.SetColumnWidths(150, 150, 150)
		headerStyle := document.NewCellStyleBuilder().Bold().Background("4472C4").Color("FFFFFF").Padding(5).Align("center", "center").Border().Build()

		_ = tbl.Row(0).Cell(0).AddParagraph("Quarter").Style(headerStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Revenue").Style(headerStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Growth").Style(headerStyle)

		rowStyle := document.NewCellStyleBuilder().Padding(5).Border().Build()
		_ = tbl.Row(1).Cell(0).AddParagraph("Q1 2026").Style(rowStyle)
		_ = tbl.Row(1).Cell(1).AddParagraph("$4.2M").Style(rowStyle)
		_ = tbl.Row(1).Cell(2).AddParagraph("+12.5%").Style(rowStyle)

		_ = tbl.Row(2).Cell(0).AddParagraph("Q2 2026 (Est)").Style(rowStyle)
		_ = tbl.Row(2).Cell(1).AddParagraph("$4.8M").Style(rowStyle)
		_ = tbl.Row(2).Cell(2).AddParagraph("+14.2%").Style(rowStyle)
	}

	if err := pdf.Export("enterprise_v2_production.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved enterprise_v2_production.pdf")
	}
}
