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

	logger.Info().Msg("--- PDF Search & Highlight Example ---")

	// 1. Create a PDF with highlighted text
	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	// Add a title
	_ = pdf.AddParagraph("PDF Search and Highlight Test", document.NewCellStyleBuilder().Bold().Size(20).Build())

	_ = pdf.AddParagraph("This document demonstrates the text highlighting feature and the search capabilities of Thoth.", document.NewCellStyleBuilder().Build())

	// Highlighted text using different colors
	yellowHighlight := document.NewCellStyleBuilder().Background("FFFF00").Build()
	_ = pdf.AddParagraph("This paragraph is highlighted in yellow.", yellowHighlight)

	cyanHighlight := document.NewCellStyleBuilder().Background("00FFFF").Build()
	_ = pdf.AddParagraph("This paragraph is highlighted in cyan.", cyanHighlight)

	_ = pdf.AddParagraph("Searching for keywords will work even if they are highlighted.", document.NewCellStyleBuilder().Build())

	outputFile := "highlight_search_output.pdf"
	if err := pdf.Export(outputFile); err != nil {
		logger.Error().Err(err).Msg("Failed to export PDF")
		return
	}
	logger.Info().Str("file", outputFile).Msg("PDF created successfully with highlights")

	// 2. Search in the created PDF
	// To search, we open the exported file which populates the internal object state
	searchDoc, err := thoth.PDF().Open(ctx, outputFile)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to open PDF for searching")
		return
	}
	defer searchDoc.Close()

	keywords := []string{"highlighted", "yellow", "cyan", "Thoth", "missing"}
	results, err := searchDoc.Search(keywords)
	if err != nil {
		logger.Error().Err(err).Msg("Search failed")
		return
	}

	fmt.Printf("\nSearch Results for %v:\n", keywords)
	if len(results) == 0 {
		fmt.Println("No results found.")
	}
	for _, res := range results {
		fmt.Printf("- Found '%s'\n", res.Keyword)
	}

	// 3. Read all content
	fullText, err := searchDoc.ReadContent()
	if err != nil {
		logger.Error().Err(err).Msg("ReadContent failed")
		return
	}
	fmt.Printf("\nFull extracted text:\n\"%s\"\n", fullText)

	// 4. Replace text
	logger.Info().Msg("Testing Replace feature")
	replacements := map[string]string{
		"yellow": "GREEN",
		"cyan":   "MAGENTA",
	}
	if err := searchDoc.Replace(replacements); err != nil {
		logger.Error().Err(err).Msg("Replace failed")
		return
	}

	replacedFile := "replaced_output.pdf"
	if err := searchDoc.Export(replacedFile); err != nil {
		logger.Error().Err(err).Msg("Failed to export replaced PDF")
		return
	}
	logger.Info().Str("file", replacedFile).Msg("PDF with replaced text created")

	// Verify replacement
	finalText, _ := searchDoc.ReadContent()
	fmt.Printf("\nExtracted text after replacement:\n\"%s\"\n", finalText)

	// Cleanup
	_ = os.Remove(outputFile)
	_ = os.Remove(replacedFile)
}
