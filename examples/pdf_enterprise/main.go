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
	ensureTestImage()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	thoth := core.New(logger)
	ctx := context.Background()

	logger.Info().Msg("--- PDF Enterprise Features Showcase ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	// 1. Production Security: AES-256 Encryption
	_ = pdf.SetPassword("EnterprisePassword123")

	_ = pdf.SetPageSettings(document.NewPageSettingsBuilder().WithPaperType(document.PaperA4).Build())
	_ = pdf.SetHeader("Enterprise PDF Solution", document.NewCellStyleBuilder().Bold().Size(14).Color("2E75B6").Align("center", "center").Build())
	_ = pdf.SetWatermark("CONFIDENTIAL", document.NewCellStyleBuilder().Color("FF0000").Size(80).Build())

	// 2. Interactive Forms (Advanced)
	_ = pdf.AddHeading("1. Interactive Form Fields", 1)
	_ = pdf.AddParagraph("Fill out the following fields (Visible in most PDF readers):", document.NewCellStyleBuilder().Italic().Build())

	_ = pdf.AddParagraph("Full Name:", document.NewCellStyleBuilder().Size(10).Build())
	_ = pdf.AddTextField("full_name", 150, 0, 200, 20)

	_ = pdf.AddParagraph("Department:", document.NewCellStyleBuilder().Size(10).Build())
	_ = pdf.AddComboBox("dept", 150, 0, 200, 20, "Engineering", "Finance", "HR", "Legal")

	_ = pdf.AddParagraph("Security Level:", document.NewCellStyleBuilder().Size(10).Build())
	_ = pdf.AddRadioButton("security_level", 250, 0, "Level 1", "Level 2", "Level 3")

	_ = pdf.AddParagraph("I agree to the terms and conditions:", document.NewCellStyleBuilder().Size(10).Build())
	_ = pdf.AddCheckbox("terms_agreed", 250, 0)

	// 2. Vertical Alignment in Tables
	_ = pdf.AddHeading("2. Vertical Alignment in Tables", 1)
	if tbl, err := pdf.AddTable(3, 3); err == nil {
		_ = tbl.SetColumnWidths(150, 150, 150)

		hStyle := document.NewCellStyleBuilder().Bold().Background("E0E0E0").Border().Align("center", "center").Build()
		_ = tbl.Row(0).Cell(0).AddParagraph("Top Aligned").Style(hStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Center Aligned").Style(hStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Bottom Aligned").Style(hStyle)

		// Set a large row height by adding more content to one cell
		_ = tbl.Row(1).Cell(0).AddParagraph("Base\nHeight\nCell").Style(document.NewCellStyleBuilder().Border().Build())

		// Demonstrate vertical alignments
		_ = tbl.Row(1).Cell(0).AddParagraph("Top").Style(document.NewCellStyleBuilder().Align("", "top").Border().Build())
		_ = tbl.Row(1).Cell(1).AddParagraph("Center").Style(document.NewCellStyleBuilder().Align("", "center").Border().Build())
		_ = tbl.Row(1).Cell(2).AddParagraph("Bottom").Style(document.NewCellStyleBuilder().Align("", "bottom").Border().Build())

		// Mixed content in vertically centered cell
		_ = tbl.Row(2).Cell(1).
			AddImage("test.png", 30, 30).
			AddParagraph("Centered Image & Text").
			Style(document.NewCellStyleBuilder().Align("center", "center").Border().Build())
	}

	// 3. Lists with Styles
	_ = pdf.AddHeading("3. Styled Lists", 1)
	_ = pdf.AddList([]string{
		"Enterprise-grade PDF generation",
		"Interactive AcroForms support",
		"Precise layout with custom font metrics",
		"Accessible Tagged PDF (PDF/UA)",
	}, false, document.NewCellStyleBuilder().Color("4472C4").Build())

	// 4. Multi-column Support
	_ = pdf.AddPageBreak()
	_ = pdf.SetPageSettings(document.NewPageSettingsBuilder().WithPaperType(document.PaperA4).WithColumns(2, 20).Build())
	_ = pdf.AddHeading("4. Multi-column Support", 1)
	for i := range 10 {
		_ = pdf.AddParagraph(fmt.Sprintf("This is paragraph %d in a multi-column layout. It demonstrates how text flows between columns and pages automatically when the content exceeds the available space.", i+1), document.NewCellStyleBuilder().Build())
	}

	if err := pdf.Export("enterprise_output.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved enterprise_output.pdf")
	}
}

func ensureTestImage() {
	const path = "test.png"
	if _, err := os.Stat(path); err == nil {
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			c := color.RGBA{46, 117, 182, 255}
			if x < 2 || x > 97 || y < 2 || y > 97 {
				c = color.RGBA{0, 0, 0, 255}
			}
			img.Set(x, y, c)
		}
	}
	f, _ := os.Create(path)
	defer f.Close()
	_ = png.Encode(f, img)
}
