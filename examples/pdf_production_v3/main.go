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

	logger.Info().Msg("--- PDF Production v3: High-Precision & Interactive ---")

	pdf, err := thoth.PDF().New(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create PDF")
		return
	}
	defer pdf.Close()

	// 1. Metadata & Security
	_ = pdf.SetMetadata(document.Metadata{
		Title:    "Interactive Q1 2026 Audit Report",
		Author:   "Thoth Security Auditor",
		Subject:  "Form Fields & Transparency",
		Keywords: []string{"Interactive", "Form", "Transparency", "Precision"},
	})
	_ = pdf.SetPassword("thoth2026")

	// 2. Page Settings
	_ = pdf.SetPageSettings(document.PageSettings{
		PaperType:   document.PaperA4,
		Orientation: document.OrientationPortrait,
		Margins:     document.Margins{Top: 50, Bottom: 50, Left: 50, Right: 50},
	})

	// 3. Transparent Design Elements
	_ = pdf.AddHeading("Audit Summary & Findings", 1)

	// Overlapping transparent circles
	_ = pdf.DrawEllipse(100, 650, 80, 80, document.CellStyle{Background: "FF0000", Opacity: 0.5})
	_ = pdf.DrawEllipse(140, 650, 80, 80, document.CellStyle{Background: "00FF00", Opacity: 0.5})
	_ = pdf.DrawEllipse(120, 610, 80, 80, document.CellStyle{Background: "0000FF", Opacity: 0.5})

	_ = pdf.AddParagraph("This document demonstrates high-precision layout with transparent design elements. The circles above use 50% opacity and overlapping ExtGStates.",
		document.CellStyle{SpacingBefore: 120, Italic: true})

	// 4. Interactive Form Section
	_ = pdf.AddHeading("Internal Audit Form (Interactive)", 2)
	_ = pdf.AddParagraph("Please complete the following details for the record:", document.CellStyle{SpacingAfter: 15})

	// Use absolute positioning for form labels and fields via Paragraph indentation/spacing
	// Or just manual positioning if we had it. Since we are flow-based, we'll use small labels.

	_ = pdf.AddParagraph("Full Name:", document.CellStyle{Bold: true})
	_ = pdf.AddTextField("full_name", 150, 480, 300, 20)

	_ = pdf.AddParagraph("Department:", document.CellStyle{Bold: true, SpacingBefore: 10})
	_ = pdf.AddComboBox("dept", 150, 450, 200, 20, "Finance", "Engineering", "Legal", "Operations")

	_ = pdf.AddParagraph("Approval Status:", document.CellStyle{Bold: true, SpacingBefore: 10})
	_ = pdf.AddCheckbox("approved", 150, 420)
	_ = pdf.AddParagraph("      Approved", document.CellStyle{Size: 9})

	_ = pdf.AddParagraph("Risk Level:", document.CellStyle{Bold: true, SpacingBefore: 10})
	_ = pdf.AddRadioButton("risk", 150, 390, "Low", "Medium", "High")
	_ = pdf.AddParagraph("      Low      Medium      High", document.CellStyle{Size: 9})

	// 5. High-Precision Table with Translucency
	_ = pdf.AddHeading("Risk Assessment Matrix", 2)
	if tbl, err := pdf.AddTable(3, 3); err == nil {
		_ = tbl.SetColumnWidths(150, 150, 150)
		headerStyle := document.CellStyle{Bold: true, Background: "4472C4", Color: "FFFFFF", Padding: 5, Border: true}

		_ = tbl.Row(0).Cell(0).AddParagraph("Factor").Style(headerStyle)
		_ = tbl.Row(0).Cell(1).AddParagraph("Impact").Style(headerStyle)
		_ = tbl.Row(0).Cell(2).AddParagraph("Mitigation").Style(headerStyle)

		rowStyle := document.CellStyle{Padding: 5, Border: true}
		_ = tbl.Row(1).Cell(0).AddParagraph("Security").Style(rowStyle)
		_ = tbl.Row(1).Cell(1).AddParagraph("High").Style(document.CellStyle{Background: "FF0000", Opacity: 0.3, Border: true, Padding: 5})
		_ = tbl.Row(1).Cell(2).AddParagraph("Encryption").Style(rowStyle)

		_ = tbl.Row(2).Cell(0).AddParagraph("Compliance").Style(rowStyle)
		_ = tbl.Row(2).Cell(1).AddParagraph("Medium").Style(document.CellStyle{Background: "FFFF00", Opacity: 0.3, Border: true, Padding: 5})
		_ = tbl.Row(2).Cell(2).AddParagraph("Auditing").Style(rowStyle)
	}

	// 6. Final Footer with Page Numbers
	_ = pdf.SetFooter("Confidential Audit Report - Page {n} of {nb}", document.CellStyle{Size: 8, Color: "999999", Horizontal: "center"})

	if err := pdf.Export("audit_v3_production.pdf"); err != nil {
		logger.Error().Err(err).Msg("Failed to save PDF")
	} else {
		logger.Info().Msg("Successfully saved audit_v3_production.pdf")
	}
}
