package word

import (
	"context"
	"io"
	"testing"

	"github.com/gsoultan/thoth/document"
)

func TestDocument_AdvancedWordFeatures(t *testing.T) {
	doc := NewDocument()
	defer doc.Close()
	ctx := context.Background()

	w, ok := doc.(document.WordProcessor)
	if !ok {
		t.Fatal("Document does not implement WordProcessor")
	}

	// 1. Page Settings
	_ = w.SetPageSettings(document.PageSettings{
		PaperType:   document.PaperA4,
		Orientation: document.OrientationPortrait,
		Margins: document.Margins{
			Top:    72, // 1 inch
			Bottom: 72,
			Left:   72,
			Right:  72,
		},
	})

	// 2. Metadata
	_ = w.SetMetadata(document.Metadata{
		Title:   "Production Test Report",
		Author:  "Thoth AI",
		Subject: "Testing Advanced Features",
	})

	// 3. Headings and TOC
	_ = w.AddHeading("Table of Contents", 1)
	_ = w.AddTableOfContents()
	_ = w.AddPageBreak()

	// 4. Content with Styles
	_ = w.AddHeading("1. Introduction", 1)
	_ = w.AddParagraph("This is a production test document demonstrating advanced Word features.",
		document.NewCellStyleBuilder().Italic().Color("333333").Build())

	// 5. Lists
	_ = w.AddHeading("1.1 Features", 2)
	_ = w.AddList([]string{"Tables", "Images", "Shapes", "Fields"}, true)

	// 6. Hyperlinks and Bookmarks
	_ = w.AddBookmark("section2")
	_ = w.AddHeading("2. Advanced Elements", 1)
	_ = w.AddHyperlink("Visit Thoth Repository", "https://github.com/gsoultan/thoth")

	// 7. Drawings (Shapes)
	_ = w.AddParagraph("Below is a vector line drawing:", document.CellStyle{})
	_ = w.DrawLine(50, 50, 200, 50, document.NewCellStyleBuilder().Color("FF0000").BorderWidth(2).Build())
	_ = w.DrawRect(50, 70, 100, 50, document.NewCellStyleBuilder().Background("EFEFEF").BorderColor("0000FF").Build())

	// 8. Tables with Merged Cells
	tbl, _ := w.AddTable(3, 3)
	_ = tbl.MergeCells(0, 0, 1, 3) // Merge first row
	_ = tbl.Row(0).Cell(0).AddParagraph("Merged Header").Style(document.NewCellStyleBuilder().Bold().Align("center", "center").Build())
	_ = tbl.Row(1).Cell(0).AddParagraph("Cell 1,0")
	_ = tbl.Row(1).Cell(1).AddParagraph("Cell 1,1")

	// 9. Form Fields
	_ = w.AddHeading("3. Forms", 1)
	_ = w.AddParagraph("Please fill the form below:", document.CellStyle{})
	_ = w.AddTextField("UserName", 0, 0, 100, 20)
	_ = w.AddCheckbox("Agreed", 0, 0)

	// 10. Watermark and Footnotes
	_ = w.SetWatermark("DRAFT", document.NewCellStyleBuilder().Color("C0C0C0").Build())
	_ = w.AddFootnote("This is a test footnote.")

	// Export to check if it's corrupt-free (optional check in manual run)
	err := w.Save(ctx, io.Discard)
	if err != nil {
		t.Errorf("Failed to save document: %v", err)
	}
}
