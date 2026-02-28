package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/gsoultan/thoth/core"
	"github.com/gsoultan/thoth/document"
	"github.com/rs/zerolog"
)

func main() {
	// Ensure test image exists for the example
	ensureTestImage()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	thoth := core.New(logger)
	ctx := context.Background()

	logger.Info().Msg("--- Final Production PDF Example ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	_ = pdf.SetPageSettings(document.NewPageSettingsBuilder().WithPaperType(document.PaperA4).Build())
	_ = pdf.SetHeader("Thoth Production Ready Report", document.NewCellStyleBuilder().Bold().Size(14).Color("4F81BD").Align("center", "center").Build())
	_ = pdf.SetFooter("Page {n} of {nb}", document.NewCellStyleBuilder().Align("right", "center").Size(9).Build())

	// 1. Aligned Body Images
	_ = pdf.AddParagraph("1. Aligned Body Images", document.NewCellStyleBuilder().Bold().Size(18).Name("Heading1").Build())
	_ = pdf.InsertImage("test.png", 50, 50, document.NewCellStyleBuilder().Align("left", "center").Build())
	_ = pdf.AddParagraph("Image left-aligned ^", document.NewCellStyleBuilder().Size(8).Build())
	_ = pdf.InsertImage("test.png", 50, 50, document.NewCellStyleBuilder().Align("center", "center").Build())
	_ = pdf.AddParagraph("Image center-aligned ^", document.NewCellStyleBuilder().Align("center", "center").Size(8).Build())
	_ = pdf.InsertImage("test.png", 50, 50, document.NewCellStyleBuilder().Align("right", "center").Build())
	_ = pdf.AddParagraph("Image right-aligned ^", document.NewCellStyleBuilder().Align("right", "center").Size(8).Build())

	_ = pdf.AddParagraph("2. Complex Table with Multi-Items and Wraps", document.NewCellStyleBuilder().Bold().Size(18).Name("Heading1").Build())

	// 2. Complex Table
	if tbl, err := pdf.AddTable(4, 3); err == nil {
		hStyle := document.NewCellStyleBuilder().Bold().Background("D9D9D9").Border().Align("center", "center").Build()
		_ = tbl.Row(0).Cell(0).AddParagraph("Service").Style(hStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Details (Multi-Paragraph & Wraps)").Style(hStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Status & Logo").Style(hStyle)

		dStyle := document.NewCellStyleBuilder().Border().Build()

		// Row 1: Multi-paragraph and wrapped text
		_ = tbl.Row(1).Cell(0).AddParagraph("Reporting").Style(dStyle)
		_ = tbl.Row(1).Cell(1).
			AddParagraph("Primary: Full monthly reports including all metrics and performance analysis.").
			Style(dStyle).
			AddParagraph("Secondary: Weekly summaries of key indicators.").
			Style(document.NewCellStyleBuilder().Italic().Size(9).Build()).
			AddParagraph("This is a very long paragraph that should wrap multiple times inside the table cell to demonstrate dynamic row height calculation in action.").
			Style(document.NewCellStyleBuilder().Size(8).Build())
		_ = tbl.Row(1).Cell(2).
			AddImage("test.png", 30, 30).
			AddParagraph("Active").
			Style(document.NewCellStyleBuilder().Color("00B050").Bold().Align("center", "center").Border().Build())

		// Row 2: Merged cell with image and text
		_ = tbl.MergeCells(2, 0, 1, 2)
		_ = tbl.Row(2).Cell(0).
			AddParagraph("Merged Cell (2 columns)").
			Style(document.NewCellStyleBuilder().Background("F2F2F2").Border().Build()).
			AddParagraph("Contains multiple items including an image:").
			Style(dStyle).
			AddImage("test.png", 20, 20)
		_ = tbl.Row(2).Cell(2).AddParagraph("N/A").Style(dStyle)

		// Row 3: Image alignment in cell
		_ = tbl.Row(3).Cell(0).AddParagraph("Right-Aligned Logo").Style(dStyle)
		_ = tbl.Row(3).Cell(1).AddImage("test.png", 40, 40).Style(document.NewCellStyleBuilder().Align("right", "center").Border().Build())
		_ = tbl.Row(3).Cell(2).AddImage("test.png", 40, 40).Style(document.NewCellStyleBuilder().Align("center", "center").Border().Build())
	}

	_ = pdf.AddParagraph("3. Large Table Spanning Pages", document.NewCellStyleBuilder().Bold().Size(18).Name("Heading1").Build())

	// 3. Large Table
	if tbl, err := pdf.AddTable(30, 2); err == nil {
		headerStyle := document.NewCellStyleBuilder().Bold().Background("4F81BD").Color("FFFFFF").Border().Align("center", "center").Build()
		_ = tbl.Row(0).Cell(0).AddParagraph("ID").Style(headerStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Value Description").Style(headerStyle)

		for r := 1; r < 30; r++ {
			_ = tbl.Row(r).Cell(0).AddParagraph(fmt.Sprintf("%03d", r)).Style(document.NewCellStyleBuilder().Border().Build())
			_ = tbl.Row(r).Cell(1).AddParagraph("Data point for production testing with dynamic heights.").Style(document.NewCellStyleBuilder().Border().Build())
		}
	}

	if err := pdf.Export("final_output.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved final_output.pdf")
	}
}

func ensureTestImage() {
	const path = "test.png"
	if _, err := os.Stat(path); err == nil {
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Create a blue square with a border
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			c := color.RGBA{0, 122, 255, 255} // Blue
			if x < 5 || x > 94 || y < 5 || y > 94 {
				c = color.RGBA{0, 0, 0, 255} // Black border
			}
			img.Set(x, y, c)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	_ = png.Encode(f, img)
}
