package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
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
	runPDFQuotationExample(ctx, thoth, logger)
	runWordQuotationExample(ctx, thoth, logger)
	runPDFSearchHighlightExample(ctx, thoth, logger)

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

func runPDFSearchHighlightExample(ctx context.Context, thoth *core.Thoth, logger zerolog.Logger) {
	logger.Info().Msg("--- PDF Search & Highlight Example ---")

	// 1. Create a PDF with highlighted text
	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	// Highlighted text using different colors
	yellowHighlight := document.NewCellStyleBuilder().Background("FFFF00").Build()
	_ = pdf.AddParagraph("This paragraph is highlighted in yellow.", yellowHighlight)

	cyanHighlight := document.NewCellStyleBuilder().Background("00FFFF").Build()
	_ = pdf.AddParagraph("This paragraph is highlighted in cyan.", cyanHighlight)

	outputFile := "highlight_search_output.pdf"
	if err := pdf.Export(outputFile); err != nil {
		logger.Error().Err(err).Msg("Failed to export PDF")
		return
	}

	// 2. Search in the created PDF
	searchDoc, err := thoth.PDF().Open(ctx, outputFile)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to open PDF for searching")
		return
	}
	defer searchDoc.Close()

	keywords := []string{"highlighted", "yellow", "cyan"}
	results, err := searchDoc.Search(keywords)
	if err != nil {
		logger.Error().Err(err).Msg("Search failed")
		return
	}

	for _, res := range results {
		fmt.Printf("- Found '%s' in PDF\n", res.Keyword)
	}

	// 3. Replace text
	replacements := map[string]string{
		"yellow": "GREEN",
	}
	if err := searchDoc.Replace(replacements); err != nil {
		logger.Error().Err(err).Msg("Replace failed")
		return
	}

	if err := searchDoc.Export("highlight_replaced_output.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to export replaced PDF")
	}

	// Cleanup
	_ = os.Remove(outputFile)
	_ = os.Remove("highlight_replaced_output.pdf")
}

func runPDFQuotationExample(ctx context.Context, thoth *core.Thoth, logger zerolog.Logger) {
	logger.Info().Msg("--- PDF Quotation Example (Complex Layout) ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	// Settings
	_ = pdf.SetPageSettings(document.NewPageSettingsBuilder().
		WithPaperType(document.PaperA4).
		WithMargins(document.Margins{Top: 20, Bottom: 20, Left: 20, Right: 20}).
		Build())

	// Header/Footer
	_ = pdf.SetFooter("Generated by Thoth PDF Engine - Page {n} of {nb}",
		document.NewCellStyleBuilder().Size(8).Color("808080").Align("center", "bottom").Build())

	// 1. Header with Title
	if tbl, err := pdf.AddTable(1, 2); err == nil {
		_ = tbl.SetColumnWidths(300, 200)
		_ = tbl.Row(0).Cell(0).AddParagraph("THOTH SOLUTIONS").
			Style(document.NewCellStyleBuilder().Bold().Size(20).Color("2E5481").Build())
		_ = tbl.Row(0).Cell(1).AddParagraph("QUOTATION").
			Style(document.NewCellStyleBuilder().Bold().Size(24).Color("A6A6A6").Align("right", "center").Build())
	}

	_ = pdf.AddParagraph("\n", document.NewCellStyleBuilder().Size(10).Build())

	// 2. Info Section
	if tbl, err := pdf.AddTable(1, 2); err == nil {
		_ = tbl.SetColumnWidths(250, 250)
		_ = tbl.Row(0).Cell(0).
			AddParagraph("FROM:").Style(document.NewCellStyleBuilder().Bold().Size(9).Color("808080").Build()).
			AddParagraph("Thoth Solutions Ltd.\n123 Tech Avenue\nLondon, UK").Style(document.NewCellStyleBuilder().Size(10).Build())

		_ = tbl.Row(0).Cell(1).
			AddParagraph("BILL TO:").Style(document.NewCellStyleBuilder().Bold().Size(9).Color("808080").Align("right", "top").Build()).
			AddParagraph("Global Corp Inc.\n456 Enterprise Way\nNew York, USA").Style(document.NewCellStyleBuilder().Size(10).Align("right", "top").Build())
	}

	_ = pdf.AddParagraph("\n", document.NewCellStyleBuilder().Size(10).Build())

	// 3. Line Items
	if tbl, err := pdf.AddTable(4, 4); err == nil {
		_ = tbl.SetColumnWidths(250, 50, 100, 100)
		hStyle := document.NewCellStyleBuilder().Bold().Background("2E5481").Color("FFFFFF").Border().Align("center", "center").Build()
		_ = tbl.Row(0).Cell(0).AddParagraph("Description").Style(hStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Qty").Style(hStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Unit Price").Style(hStyle)
		_ = tbl.Row(0).Cell(3).AddParagraph("Total").Style(hStyle)

		dStyle := document.NewCellStyleBuilder().Border().Padding(5).Build()
		_ = tbl.Row(1).Cell(0).AddParagraph("Enterprise License").Style(dStyle)
		_ = tbl.Row(1).Cell(1).AddParagraph("1").Style(dStyle)
		_ = tbl.Row(1).Cell(2).AddParagraph("$2000.00").Style(dStyle)
		_ = tbl.Row(1).Cell(3).AddParagraph("$2000.00").Style(dStyle)

		_ = tbl.Row(2).Cell(0).AddParagraph("Support Contract").Style(dStyle)
		_ = tbl.Row(2).Cell(1).AddParagraph("1").Style(dStyle)
		_ = tbl.Row(2).Cell(2).AddParagraph("$500.00").Style(dStyle)
		_ = tbl.Row(2).Cell(3).AddParagraph("$500.00").Style(dStyle)

		// Subtotal in merged cell
		_ = tbl.MergeCells(3, 0, 1, 3)
		_ = tbl.Row(3).Cell(0).AddParagraph("TOTAL").Style(document.NewCellStyleBuilder().Bold().Align("right", "center").Padding(5).Build())
		_ = tbl.Row(3).Cell(3).AddParagraph("$2500.00").Style(document.NewCellStyleBuilder().Bold().Border().Align("right", "center").Padding(5).Build())
	}

	if err := pdf.Export("quotation_main_example.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save Quotation PDF")
	} else {
		logger.Info().Msg("Successfully saved quotation_main_example.pdf")
	}
}

func runWordQuotationExample(ctx context.Context, thoth *core.Thoth, logger zerolog.Logger) {
	logger.Info().Msg("--- Word Quotation Example (Complex Layout) ---")

	ensureTestImageLocal()

	word, err := thoth.Word().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Word document")
		return
	}
	defer word.Close()

	// Settings
	_ = word.SetPageSettings(document.NewPageSettingsBuilder().
		WithPaperType(document.PaperA4).
		WithMargins(document.Margins{Top: 20, Bottom: 20, Left: 20, Right: 20}).
		Build())

	// Footer
	_ = word.SetFooter("Generated by Thoth Word Engine - Page {n} of {nb}",
		document.NewCellStyleBuilder().Size(8).Color("808080").Align("center", "bottom").Build())

	// 1. Header with Logo and Title
	if tbl, err := word.AddTable(1, 2); err == nil {
		_ = tbl.SetColumnWidths(350, 150)
		_ = tbl.Row(0).Cell(0).AddImage("logo_main.png", 50, 50)
		_ = tbl.Row(0).Cell(1).AddParagraph("QUOTATION").
			Style(document.NewCellStyleBuilder().Bold().Size(24).Color("2E5481").Align("right", "center").Build())
	}

	_ = word.AddParagraph("", document.NewCellStyleBuilder().Size(10).Build())

	// 2. Info Section
	if tbl, err := word.AddTable(1, 2); err == nil {
		_ = tbl.SetColumnWidths(250, 250)

		fromCell := tbl.Row(0).Cell(0)
		_ = fromCell.AddParagraph("FROM:").Style(document.NewCellStyleBuilder().Bold().Size(9).Color("808080").Build())
		_ = fromCell.AddParagraph("Thoth Solutions Ltd.").Style(document.NewCellStyleBuilder().Bold().Size(12).Build())
		_ = fromCell.AddParagraph("123 Document Street").Style(document.NewCellStyleBuilder().Size(10).Build())
		_ = fromCell.AddParagraph("London, UK").Style(document.NewCellStyleBuilder().Size(10).Build())

		toCell := tbl.Row(0).Cell(1)
		alignRight := document.NewCellStyleBuilder().Align("right", "top")
		_ = toCell.AddParagraph("BILL TO:").Style(alignRight.Bold().Size(9).Color("808080").Build())
		_ = toCell.AddParagraph("Global Corp Inc.").Style(alignRight.Bold().Size(12).Build())
		_ = toCell.AddParagraph("456 Enterprise Way").Style(alignRight.Size(10).Build())
		_ = toCell.AddParagraph("New York, USA").Style(alignRight.Size(10).Build())
	}

	_ = word.AddParagraph("", document.NewCellStyleBuilder().Size(10).Build())

	// 3. Line Items
	items := []struct {
		Desc  string
		Qty   int
		Price float64
	}{
		{"Word Automation Module", 1, 1500.00},
		{"Template Design Service", 2, 250.00},
		{"Priority Support", 1, 500.00},
	}

	if tbl, err := word.AddTable(len(items)+1, 4); err == nil {
		_ = tbl.SetColumnWidths(250, 50, 100, 100)
		hStyle := document.NewCellStyleBuilder().Bold().Background("2E5481").Color("FFFFFF").Border().Align("center", "center").Build()
		_ = tbl.Row(0).Cell(0).AddParagraph("Description").Style(hStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Qty").Style(hStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Unit Price").Style(hStyle)
		_ = tbl.Row(0).Cell(3).AddParagraph("Total").Style(hStyle)

		dStyle := document.NewCellStyleBuilder().Border().Padding(5).Build()
		numStyle := document.NewCellStyleBuilder().Border().Align("right", "center").Padding(5).Build()

		total := 0.0
		for i, item := range items {
			lineTotal := float64(item.Qty) * item.Price
			total += lineTotal
			r := i + 1
			_ = tbl.Row(r).Cell(0).AddParagraph(item.Desc).Style(dStyle)
			_ = tbl.Row(r).Cell(1).AddParagraph(fmt.Sprintf("%d", item.Qty)).Style(numStyle)
			_ = tbl.Row(r).Cell(2).AddParagraph(fmt.Sprintf("%.2f", item.Price)).Style(numStyle)
			_ = tbl.Row(r).Cell(3).AddParagraph(fmt.Sprintf("%.2f", lineTotal)).Style(numStyle)
		}

		_ = word.AddParagraph("", document.NewCellStyleBuilder().Size(10).Build())

		if summaryTbl, err := word.AddTable(1, 2); err == nil {
			_ = summaryTbl.SetColumnWidths(400, 100)
			_ = summaryTbl.Row(0).Cell(0).AddParagraph("GRAND TOTAL").Style(document.NewCellStyleBuilder().Bold().Align("right", "center").Padding(5).Build())
			_ = summaryTbl.Row(0).Cell(1).AddParagraph(fmt.Sprintf("%.2f", total)).Style(document.NewCellStyleBuilder().Bold().Background("F0F0F0").Border().Align("right", "center").Padding(5).Build())
		}
	}

	if err := word.Export("quotation_main_example.docx"); err != nil {
		logger.Error().Err(err).Msg("Failed to save Quotation Word")
	} else {
		logger.Info().Msg("Successfully saved quotation_main_example.docx")
	}

	// Cleanup image
	_ = os.Remove("logo_main.png")
}

func ensureTestImageLocal() {
	if _, err := os.Stat("logo_main.png"); err == nil {
		return
	}
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{46, 84, 129, 255})
		}
	}
	f, _ := os.Create("logo_main.png")
	defer f.Close()
	_ = png.Encode(f, img)
}
