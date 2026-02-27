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

	logger.Info().Msg("--- Production PDF Example ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	_ = pdf.SetPageSettings(document.NewPageSettingsBuilder().WithPaperType(document.PaperA4).Build())
	_ = pdf.SetFooter("Page {n} of {nb}", document.NewCellStyleBuilder().Align("center", "center").Size(9).Build())

	// 1. Title (will be a bookmark)
	_ = pdf.AddParagraph("Thoth Production-Ready PDF", document.NewCellStyleBuilder().Bold().Size(24).Name("Title").Build())

	// 2. Section 1
	_ = pdf.AddParagraph("1. Introduction", document.NewCellStyleBuilder().Bold().Size(16).Name("Heading1").Build())
	_ = pdf.AddParagraph("This document demonstrates production-ready features of the Thoth PDF module.", document.NewCellStyleBuilder().Build())

	// 3. Hyperlink
	_ = pdf.AddParagraph("Visit our website: https://github.com/gsoultan/thoth", document.NewCellStyleBuilder().Color("0000FF").Link("https://github.com/gsoultan/thoth").Build())

	// 4. Section 2 with many pages to test {nb} and compression
	for i := 2; i <= 5; i++ {
		_ = pdf.AddParagraph(fmt.Sprintf("%d. Section %d", i, i), document.NewCellStyleBuilder().Bold().Size(16).Name("Heading1").Build())
		for j := 0; j < 40; j++ {
			_ = pdf.AddParagraph(fmt.Sprintf("This is some dummy text to fill the page. Line %d of Section %d. This will help us test compression and pagination across multiple pages.", j, i), document.NewCellStyleBuilder().Build())
		}
		_ = pdf.AddPageBreak()
	}

	if err := pdf.Export("production_output.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved production_output.pdf")
	}
}
