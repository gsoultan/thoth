package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gsoultan/thoth/core"
	"github.com/gsoultan/thoth/document"
	"github.com/rs/zerolog"
)

const (
	// MsgStarting is the Thoth application starting message.
	MsgStarting = "Starting Thoth application"
	// MsgFinished is the Thoth application finished message.
	MsgFinished = "Thoth application finished successfully"
	// MsgFailed is the Thoth application failed message.
	MsgFailed = "Thoth processing failed"
)

func main() {
	// 1. Initialize structured logging
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().Timestamp().Logger()

	logger.Info().Msg(MsgStarting)

	// 2. Create Thoth instance
	thoth := core.New(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. Run format-specific examples
	runExcelExample(ctx, thoth, logger)
	runWordExample(ctx, thoth, logger)
	runPDFExample(ctx, thoth, logger)

	// 4. Run infrastructure/storage example
	runStorageExample(ctx, logger)

	logger.Info().Msg(MsgFinished)
}

func runExcelExample(ctx context.Context, thoth *core.Thoth, logger zerolog.Logger) {
	logger.Info().Msg("--- Excel Example (Fluent API) ---")

	// 1. Create a new Excel spreadsheet using the fluent entry point
	ss, err := thoth.Excel().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Excel document")
		return
	}
	defer ss.Close()

	// 2. Fluent Styling
	headerStyle := document.NewCellStyleBuilder().
		Bold().
		Background("EFEFEF").
		Border().
		Color("0000FF").
		Align("center", "center").
		Build()

	currencyStyle := document.NewCellStyleBuilder().
		NumberFormat("$#,##0.00").
		Build()

	// 3. Chainable operations on a scoped sheet handle
	sheet, err := ss.Sheet("Sales Data")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get sheet")
		return
	}

	sheet.Cell("A1").Set("Monthly Sales Report").Style(headerStyle)
	sheet.MergeCells("A1:B1").
		SetColumnWidth(1, 24).
		SetColumnWidth(2, 12)

	sheet.Cell("A2").Set("Item Name").Style(headerStyle)
	sheet.Cell("B2").Set("Quantity").Style(headerStyle)

	sheet.Cell("A3").Set("Product A")
	sheet.Cell("B3").Set(150)

	sheet.Cell("A4").Set("Price")
	sheet.Cell("B4").Set(29.99).Style(currencyStyle)

	// New: AutoFilter and Freeze Panes
	sheet.AutoFilter("A2:B4").
		FreezePanes(0, 2)

	// 4. Handle errors if needed (terminal check)
	if err := sheet.Err(); err != nil {
		logger.Error().Err(err).Msg("Error during Excel operations")
	}

	// 5. Save using fluent Export
	if err := ss.Export("sales_report.xlsx"); err != nil {
		logger.Error().Err(err).Msg("Failed to save Excel document")
	} else {
		logger.Info().Msg("Successfully saved sales_report.xlsx")
	}
}

func runWordExample(ctx context.Context, thoth *core.Thoth, logger zerolog.Logger) {
	logger.Info().Msg("--- Word Example (Fluent API) ---")

	// 1. Create a new Word document using the fluent entry point
	wp, err := thoth.Word().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Word document")
		return
	}
	defer wp.Close()

	// 2. Scoped content addition
	_ = wp.AddParagraph("Internal Memorandum", document.NewCellStyleBuilder().Name("Heading1").Bold().Size(14).Align("center", "center").Build())
	_ = wp.AddParagraph("This is a justified paragraph demonstration of the Word processor. It ensures that the text stretches across the page correctly for a professional look.", document.NewCellStyleBuilder().Align("justify", "center").Build())
	_ = wp.AddParagraph("And here is a right-aligned paragraph for the signature or date.", document.NewCellStyleBuilder().Align("right", "center").Build())
	_ = wp.AddPageBreak()

	// 3. Systematic Table API (Row/Cell Scoped)
	if tbl, err := wp.AddTable(3, 2); err == nil {
		tbl.MergeCells(0, 0, 1, 2).
			SetStyle("TableGrid")

		// Fluent Row/Cell chaining
		tbl.Row(0).Cell(0).
			AddParagraph("Consolidated Report Header").
			Style(document.NewCellStyleBuilder().Background("D9D9D9").Border().Build())

		tbl.Row(1).Cell(0).AddParagraph("Data Point A").Style(document.NewCellStyleBuilder().Border().Build())
		tbl.Row(1).Cell(1).AddParagraph("Value A").Style(document.NewCellStyleBuilder().Border().Build())

		if err := tbl.Err(); err != nil {
			logger.Error().Err(err).Msg("Table operations failed")
		}
	}

	// 4. Save the document
	if err := wp.Export("memo_final.docx"); err != nil {
		logger.Error().Err(err).Msg("Failed to save Word document")
	} else {
		logger.Info().Msg("Successfully saved memo_final.docx")
	}
}

func runPDFExample(ctx context.Context, thoth *core.Thoth, logger zerolog.Logger) {
	logger.Info().Msg("--- PDF Example (Fluent API) ---")

	// 1. Create a new PDF document using the fluent entry point
	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	// 2. Standard document setup
	settings := document.NewPageSettingsBuilder().
		WithOrientation(document.OrientationPortrait).
		WithPaperType(document.PaperA4).
		Build()
	_ = pdf.SetPageSettings(settings)

	// 3. Fluent Table Scoping (Same as Word)
	_ = pdf.AddParagraph("Automated PDF Generation with Fluent API", document.NewCellStyleBuilder().Bold().Color("FF0000").Build())
	if tbl, err := pdf.AddTable(10, 2); err == nil {
		// Header Row
		headerStyle := document.NewCellStyleBuilder().Bold().Background("4F81BD").Color("FFFFFF").Border().Align("center", "center").Build()
		tbl.Row(0).Cell(0).AddParagraph("Metric").Style(headerStyle)
		tbl.Row(0).Cell(1).AddParagraph("Value").Style(headerStyle)

		// Data Rows
		dataStyle := document.NewCellStyleBuilder().Border().Build()
		for r := 1; r < 10; r++ {
			tbl.Row(r).Cell(0).AddParagraph(fmt.Sprintf("Metric %d", r)).Style(dataStyle)
			tbl.Row(r).Cell(1).AddParagraph(fmt.Sprintf("%d.00", r*10)).Style(dataStyle)
		}
	}

	// 4. Save the document
	if err := pdf.Export("output_example.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved output_example.pdf")
	}
}

func runStorageExample(ctx context.Context, logger zerolog.Logger) {
	logger.Info().Msg("--- Storage & S3 Example ---")

	// 1. Configure Thoth with S3 support
	s3Thoth := core.New(logger).WithS3Config(core.S3Config{
		Endpoint:        "s3.amazonaws.com",
		AccessKeyID:     "AKIA...", // Production: Use environment variables
		SecretAccessKey: "SECRET...",
		UseSSL:          true,
		Region:          "us-east-1",
	})

	// 2. Open document from S3 (demonstration will fail with fake credentials)
	logger.Info().Msg("Attempting to open from S3 (expected fail with fake credentials)")
	_, err := s3Thoth.Excel().Open(ctx, "s3://my-bucket/template.xlsx")
	if err != nil {
		logger.Warn().Err(err).Msg("S3 access denied as expected")
	}

	// 3. Open document from URL
	logger.Info().Msg("Attempting to open from URL (demonstration)")
	urlThoth := core.New(logger)
	_, err = urlThoth.Word().Open(ctx, "https://example.com/sample.docx")
	if err != nil {
		logger.Warn().Err(err).Msg("URL open failed as expected")
	}
}
