package word

import (
	"strings"
	"testing"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func TestDocument_AddParagraph(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	text := "Hello, World!"
	err := doc.AddParagraph(text)
	if err != nil {
		t.Fatalf("AddParagraph failed: %v", err)
	}

	content, err := doc.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent failed: %v", err)
	}

	if !strings.Contains(content, text) {
		t.Errorf("Expected content to contain %q, got %q", text, content)
	}
}

func TestDocument_InsertTable(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	handle, err := doc.AddTable(3, 2)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	if handle == nil {
		t.Fatal("Expected table handle, got nil")
	}

	var tableCount int
	for _, c := range doc.doc.Body.Content {
		if _, ok := c.(*xmlstructs.Table); ok {
			tableCount++
		}
	}

	if tableCount != 1 {
		t.Errorf("Expected 1 table, got %d", tableCount)
	}

	tbl, err := doc.getTable(0)
	if err != nil {
		t.Fatalf("getTable failed: %v", err)
	}
	if len(tbl.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(tbl.Rows))
	}
}

func TestDocument_AdvancedStyling(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	style := document.CellStyle{
		Bold:          true,
		Italic:        true,
		Size:          14,
		Color:         "FF0000",
		Background:    "FFFF00",
		SpacingBefore: 10,
		SpacingAfter:  10,
		Indent:        20,
	}

	err := doc.AddParagraph("Styled Text", style)
	if err != nil {
		t.Fatalf("AddParagraph with style failed: %v", err)
	}

	p, ok := doc.doc.Body.Content[0].(*xmlstructs.Paragraph)
	if !ok {
		t.Fatal("Expected paragraph, got other content type")
	}

	if p.PPr == nil {
		t.Fatal("Expected PPr, got nil")
	}
	if p.PPr.Ind == nil || p.PPr.Ind.Left != 400 { // 20 * 20
		t.Errorf("Expected Ind.Left to be 400, got %v", p.PPr.Ind)
	}
	if p.PPr.Spacing == nil || p.PPr.Spacing.Before != 200 { // 10 * 20
		t.Errorf("Expected Spacing.Before to be 200, got %v", p.PPr.Spacing.Before)
	}

	run, ok := p.Content[0].(*xmlstructs.Run)
	if !ok {
		t.Fatal("Expected run, got other content type")
	}
	if run.RPr == nil {
		t.Fatal("Expected RPr, got nil")
	}
	if run.RPr.Bold == nil {
		t.Error("Expected Bold to be set")
	}
	if run.RPr.Sz == nil || run.RPr.Sz.Val != 28 { // 14 * 2
		t.Errorf("Expected Sz.Val to be 28, got %v", run.RPr.Sz.Val)
	}
}

func TestDocument_TableAdvanced(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	table, err := doc.AddTable(2, 2)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	table.SetHeaderRows(1)
	table.Row(0).Cell(0).Style(document.CellStyle{Background: "D9D9D9", Vertical: "center"})

	tbl, _ := doc.getTable(0)
	if tbl.Rows[0].TrPr == nil || tbl.Rows[0].TrPr.TblHeader == nil {
		t.Error("Expected TblHeader to be set")
	}

	cell := tbl.Rows[0].Cells[0]
	if cell.TcPr == nil || cell.TcPr.Shd == nil || cell.TcPr.Shd.Fill != "D9D9D9" {
		t.Errorf("Expected cell shading Fill D9D9D9, got %v", cell.TcPr.Shd)
	}
	if cell.TcPr.VAlign == nil || cell.TcPr.VAlign.Val != "center" {
		t.Errorf("Expected cell vertical alignment center, got %v", cell.TcPr.VAlign)
	}
}

func TestDocument_ProductionFeatures(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	// 1. TOC
	err := doc.AddTableOfContents()
	if err != nil {
		t.Fatalf("AddTableOfContents failed: %v", err)
	}

	foundTOC := false
	for _, c := range doc.doc.Body.Content {
		if p, ok := c.(*xmlstructs.Paragraph); ok {
			for _, rc := range p.Content {
				if r, ok := rc.(*xmlstructs.Run); ok && r.InstrText != nil && strings.Contains(r.InstrText.Text, "TOC") {
					foundTOC = true
					break
				}
			}
		}
	}
	if !foundTOC {
		t.Error("TOC field not found in document")
	}

	// 2. Header with fields
	err = doc.SetHeader("Page {n} of {nb}")
	if err != nil {
		t.Fatalf("SetHeader failed: %v", err)
	}

	if len(doc.headers) != 1 {
		t.Fatal("Expected 1 header")
	}

	var hdr *xmlstructs.Header
	for _, h := range doc.headers {
		hdr = h
		break
	}

	foundField := false
	for _, c := range hdr.Content {
		if p, ok := c.(*xmlstructs.Paragraph); ok {
			for _, rc := range p.Content {
				if r, ok := rc.(*xmlstructs.Run); ok && r.InstrText != nil && r.InstrText.Text == " PAGE " {
					foundField = true
					break
				}
			}
		}
	}
	if !foundField {
		t.Error("PAGE field not found in header")
	}

	// 3. Table Cell Padding
	table, _ := doc.AddTable(1, 1)
	table.Row(0).Cell(0).Style(document.CellStyle{Padding: 10})

	tbl, _ := doc.getTable(0)
	cell := tbl.Rows[0].Cells[0]
	if cell.TcPr.TcMar == nil || cell.TcPr.TcMar.Top == nil || cell.TcPr.TcMar.Top.W != 200 {
		t.Errorf("Expected cell padding 200, got %v", cell.TcPr.TcMar)
	}

	// 4. Table Cell List
	table.Row(0).Cell(0).AddList([]string{"Item 1", "Item 2"}, true)
	tbl, _ = doc.getTable(0)
	cell = tbl.Rows[0].Cells[0]

	// Cell already has one empty paragraph from AddTable, plus 1 from AddList(Item1) and 1 from AddList(Item2)
	// But AddList should have replaced the empty one if it was empty.
	if len(cell.Content) != 2 {
		t.Errorf("Expected 2 items in cell content, got %d", len(cell.Content))
	}
	if p, ok := cell.Content[0].(*xmlstructs.Paragraph); !ok || p.PPr.NumPr == nil || p.PPr.NumPr.NumID.Val != 2 {
		t.Error("Expected numbered list in table cell")
	}
}

func TestDocument_NestedTables(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	table, err := doc.AddTable(1, 1)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	cell := table.Row(0).Cell(0)
	nestedTable := cell.AddTable(2, 2)
	if nestedTable == nil {
		t.Fatal("Expected nested table, got nil")
	}

	tbl, _ := doc.getTable(0)
	cellContent := tbl.Rows[0].Cells[0].Content
	if len(cellContent) != 1 {
		t.Fatalf("Expected 1 item in cell content, got %d", len(cellContent))
	}

	if _, ok := cellContent[0].(*xmlstructs.Table); !ok {
		t.Errorf("Expected Table in cell content, got %T", cellContent[0])
	}
}

func TestDocument_NestedTableModification(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	table, _ := doc.AddTable(1, 1)
	cell := table.Row(0).Cell(0)
	nestedTable := cell.AddTable(1, 1)

	text := "Nested Content"
	nestedTable.Row(0).Cell(0).AddParagraph(text)

	// Verify
	tbl, _ := doc.getTable(0)
	innerTbl := tbl.Rows[0].Cells[0].Content[0].(*xmlstructs.Table)
	innerCell := innerTbl.Rows[0].Cells[0]

	if len(innerCell.Content) != 1 {
		t.Fatalf("Expected 1 item in nested cell content, got %d", len(innerCell.Content))
	}

	found := false
	for _, c := range innerCell.Content {
		if p, ok := c.(*xmlstructs.Paragraph); ok {
			for _, rc := range p.Content {
				if r, ok := rc.(*xmlstructs.Run); ok && r.T == text {
					found = true
					break
				}
			}
		}
	}

	if !found {
		t.Error("Nested table content not found")
	}
}
