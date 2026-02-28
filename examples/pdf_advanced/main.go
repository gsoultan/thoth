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

	logger.Info().Msg("--- Advanced PDF Layout Example ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	_ = pdf.SetPageSettings(document.NewPageSettingsBuilder().WithPaperType(document.PaperA4).Build())

	// 1. Justified Text
	_ = pdf.AddHeading("1. Justified Text Alignment", 1)
	lorem := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur."
	_ = pdf.AddParagraph(lorem, document.NewCellStyleBuilder().
		Align("justify", "top").
		LineSpacing(1.3).
		SpacingAfter(20).
		Build())

	// 2. Indentation & Hanging Indents
	_ = pdf.AddHeading("2. Indentation & Hanging Indents", 2)
	_ = pdf.AddParagraph("This paragraph has a left indent of 40 points, moving the whole block to the right.",
		document.NewCellStyleBuilder().Indent(40).SpacingAfter(10).Build())

	_ = pdf.AddParagraph("This paragraph uses a hanging indent. The first line starts at the base position, but all subsequent lines are indented by 30 points, creating a classic bibliographical or list-like effect without using a formal list component.",
		document.NewCellStyleBuilder().Hanging(30).SpacingAfter(20).Build())

	// 3. Advanced Lists
	_ = pdf.AddHeading("3. Nested & Styled Lists", 2)
	_ = pdf.AddList([]string{"Top level item 1", "Top level item 2"}, true,
		document.NewCellStyleBuilder().SpacingAfter(5).Build())

	// Nested list via Indent
	_ = pdf.AddList([]string{"Nested item A", "Nested item B"}, false,
		document.NewCellStyleBuilder().Indent(20).SpacingAfter(10).Build())

	// 4. Table Cell Padding & Backgrounds
	_ = pdf.AddHeading("4. Table Styling: Padding & Backgrounds", 2)
	if tbl, err := pdf.AddTable(3, 3); err == nil {
		_ = tbl.SetColumnWidths(150, 150, 150)

		headerStyle := document.NewCellStyleBuilder().
			Bold().
			Background("4F81BD").
			Color("FFFFFF").
			Padding(10).
			Align("center", "center").
			Border().
			Build()

		_ = tbl.Row(0).Cell(0).AddParagraph("Category").Style(headerStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Description").Style(headerStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Status").Style(headerStyle)

		rowStyle := document.NewCellStyleBuilder().
			Padding(8).
			Border().
			Build()

		_ = tbl.Row(1).Cell(0).AddParagraph("Production").Style(rowStyle)
		_ = tbl.Row(1).Cell(1).AddParagraph("Full-cell background colors and custom padding for better readability.").Style(rowStyle)
		_ = tbl.Row(1).Cell(2).AddParagraph("Active").Style(document.NewCellStyleBuilder().
			Background("E2EFDA").
			Color("385723").
			Padding(8).
			Align("center", "center").
			Border().
			Build())

		_ = tbl.Row(2).Cell(0).AddParagraph("Research").Style(rowStyle)
		_ = tbl.Row(2).Cell(1).AddParagraph("Hanging indents also work inside table cells!").Style(rowStyle)
		_ = tbl.Row(2).Cell(2).AddList([]string{"Sub-task 1", "Sub-task 2"}, true,
			document.NewCellStyleBuilder().Size(8).Hanging(10).Border().Build())
	}

	if err := pdf.Export("advanced_layout.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved advanced_layout.pdf")
	}
}
