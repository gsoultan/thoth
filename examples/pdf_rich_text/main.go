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

	logger.Info().Msg("--- Rich Text & Advanced Layout Example ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	_ = pdf.SetPageSettings(document.NewPageSettingsBuilder().WithPaperType(document.PaperA4).Build())

	// 1. Rich Paragraph with mixed styles
	logger.Info().Msg("Adding rich paragraph...")
	spans := []document.TextSpan{
		{Text: "This paragraph contains ", Style: document.CellStyle{Size: 12}},
		{Text: "BOLD ", Style: document.CellStyle{Bold: true, Size: 12, Color: "FF0000"}},
		{Text: "and ", Style: document.CellStyle{Size: 12}},
		{Text: "ITALIC ", Style: document.CellStyle{Italic: true, Size: 12, Color: "0000FF"}},
		{Text: "text, as well as ", Style: document.CellStyle{Size: 12}},
		{Text: "highlighted background", Style: document.CellStyle{Background: "FFFF00", Size: 12}},
		{Text: ".", Style: document.CellStyle{Size: 12}},
	}
	_ = pdf.AddRichParagraph(spans)

	// 2. Superscript and Subscript
	logger.Info().Msg("Adding superscripts and subscripts...")
	formula := []document.TextSpan{
		{Text: "Chemical Formula: H", Style: document.CellStyle{Size: 14}},
		{Text: "2", Style: document.CellStyle{Subscript: true, Size: 14}},
		{Text: "O. Mathematical: E = mc", Style: document.CellStyle{Size: 14}},
		{Text: "2", Style: document.CellStyle{Superscript: true, Size: 14}},
	}
	_ = pdf.AddRichParagraph(formula)

	// 3. Keep-With-Next demonstration
	logger.Info().Msg("Demonstrating Keep-With-Next...")
	// Add space to reach bottom of page
	for range 35 {
		_ = pdf.AddParagraph("Filler line to push heading to bottom", document.NewCellStyleBuilder().Size(10).Color("CCCCCC").Build())
	}

	// This heading should move to next page because it has KeepWithNext and not enough space for it and next paragraph
	hStyle := document.NewCellStyleBuilder().Bold().Size(16).KeepWithNext().Build()
	_ = pdf.AddParagraph("I should be on Page 2 (KeepWithNext)", hStyle)
	_ = pdf.AddParagraph("I am the content following the heading. We must stay together.", document.NewCellStyleBuilder().Size(11).Build())

	// 4. Rich Text in Table Cells
	logger.Info().Msg("Adding table with rich text...")
	if tbl, err := pdf.AddTable(2, 2); err == nil {
		_ = tbl.Row(0).Cell(0).AddParagraph("Service").Style(document.NewCellStyleBuilder().Bold().Background("EEEEEE").Border().Build())
		_ = tbl.Row(0).Cell(1).AddParagraph("Rich Description").Style(document.NewCellStyleBuilder().Bold().Background("EEEEEE").Border().Build())

		richCell := []document.TextSpan{
			{Text: "Important: ", Style: document.CellStyle{Bold: true, Color: "FF0000"}},
			{Text: "Check the status ", Style: document.CellStyle{}},
			{Text: "online", Style: document.CellStyle{Link: "https://example.com", Color: "0000FF"}},
			{Text: ".", Style: document.CellStyle{}},
		}
		_ = tbl.Row(1).Cell(0).AddParagraph("Cloud API").Style(document.NewCellStyleBuilder().Border().Build())
		_ = tbl.Row(1).Cell(1).AddRichParagraph(richCell).Style(document.NewCellStyleBuilder().Border().Build())
	}

	if err := pdf.Export("rich_text_output.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved rich_text_output.pdf")
	}
}
