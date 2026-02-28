package excel

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

func TestDocument_Search(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	keywords := []string{"test", "hello"}
	results, err := doc.Search(keywords)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Since it's a mock/skeleton, it currently returns nil, nil.
	// In a real implementation, we would check the results.
	if results != nil {
		t.Errorf("Expected nil results for mock, got %v", results)
	}
}

func TestDocument_Replace(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	replacements := map[string]string{"old": "new"}
	err := doc.Replace(replacements)
	if err != nil {
		t.Errorf("Replace failed: %v", err)
	}
}

func TestDocument_SetCellStyle(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	doc.addSheet("Sheet1")

	style := document.CellStyle{
		Bold:  true,
		Size:  12,
		Color: "FF0000",
	}
	err := doc.setCellStyle("Sheet1", "A1", style)
	if err != nil {
		t.Fatalf("SetCellStyle failed: %v", err)
	}

	if doc.state.styles == nil {
		t.Fatal("Expected styles to be initialized")
	}

	if len(doc.state.styles.Fonts.Items) < 2 {
		t.Errorf("Expected at least 2 fonts (default + new), got %d", len(doc.state.styles.Fonts.Items))
	}
}

func TestDocument_AdvancedFeatures(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	// 1. Workbook Password
	err := doc.SetPassword("workbookpass")
	if err != nil {
		t.Fatalf("SetPassword failed: %v", err)
	}

	sheet, err := doc.Sheet("Advanced")
	if err != nil {
		t.Fatalf("Sheet failed: %v", err)
	}

	// 2. Page Settings
	sheet.SetPageSettings(document.PageSettings{
		PaperType:   document.PaperA4,
		Orientation: document.OrientationLandscape,
		Margins: document.Margins{
			Top: 72, Bottom: 72, Left: 54, Right: 54,
		},
	})

	// 3. Header/Footer
	sheet.SetHeader("Header Text")
	sheet.SetFooter("Footer Text")

	// 4. Protection
	sheet.Protect("sheetpass")

	// 5. Grouping
	sheet.GroupRows(1, 10, 1)
	sheet.GroupCols(1, 5, 2)

	// 6. Print Settings
	sheet.SetPrintArea("A1:F20")
	sheet.SetPrintTitles("1:1", "A:A")

	// Verify internal state
	ws := doc.sheets["Advanced"]

	if ws.SheetProtection == nil || ws.SheetProtection.Password == "" {
		t.Error("Expected sheet protection to be set")
	}

	if ws.PageSetup == nil || ws.PageSetup.Orientation != "landscape" {
		t.Error("Expected page orientation to be landscape")
	}

	if ws.HeaderFooter == nil || ws.HeaderFooter.OddHeader != "Header Text" {
		t.Error("Expected header to be set")
	}

	if doc.workbook.WorkbookProtection == nil || doc.workbook.WorkbookProtection.WorkbookPassword == "" {
		t.Error("Expected workbook protection to be set")
	}

	// Check grouping
	foundRowGroup := false
	for _, row := range ws.SheetData.Rows {
		if row.R >= 1 && row.R <= 10 && row.OutlineLevel == 1 {
			foundRowGroup = true
			break
		}
	}
	if !foundRowGroup {
		t.Error("Expected row grouping at level 1")
	}

	// Check print settings
	if doc.workbook.DefinedNames == nil || len(doc.workbook.DefinedNames.Items) == 0 {
		t.Fatal("Expected defined names for print settings")
	}

	foundPrintArea := false
	foundPrintTitles := false
	for _, dn := range doc.workbook.DefinedNames.Items {
		if dn.Name == "_xlnm.Print_Area" {
			foundPrintArea = true
			if dn.LocalSheetID == nil || *dn.LocalSheetID != 0 {
				t.Errorf("Expected LocalSheetID 0 for Print_Area, got %v", dn.LocalSheetID)
			}
		}
		if dn.Name == "_xlnm.Print_Titles" {
			foundPrintTitles = true
		}
	}
	if !foundPrintArea {
		t.Error("Print Area not found in defined names")
	}
	if !foundPrintTitles {
		t.Error("Print Titles not found in defined names")
	}
}

func TestDocument_Tables(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	sheet, _ := doc.Sheet("TableSheet")
	_ = sheet.Cell("A1").Set("Header1")
	_ = sheet.Cell("B1").Set("Header2")
	_ = sheet.Cell("A2").Set("Value1")
	_ = sheet.Cell("B2").Set("Value2")

	sheet.AddTable("A1:B2", "MyTable")
	if err := sheet.Err(); err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	// Verify table was created
	if len(doc.tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(doc.tables))
	}

	tablePath := "xl/tables/table1.xml"
	table, ok := doc.tables[tablePath]
	if !ok {
		t.Fatalf("Table not found at %s", tablePath)
	}

	if table.Name != "MyTable" {
		t.Errorf("Expected table name 'MyTable', got '%s'", table.Name)
	}

	if table.Ref != "A1:B2" {
		t.Errorf("Expected ref 'A1:B2', got '%s'", table.Ref)
	}

	if len(table.TableColumns.Items) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(table.TableColumns.Items))
	}

	if table.TableColumns.Items[0].Name != "Header1" {
		t.Errorf("Expected first column header 'Header1', got '%s'", table.TableColumns.Items[0].Name)
	}

	// Verify worksheet reference
	ws := doc.sheets["TableSheet"]
	if ws.TableParts == nil || len(ws.TableParts.Items) != 1 {
		t.Fatal("Worksheet missing TableParts reference")
	}

	// Verify relationship
	rels := doc.sheetRels["TableSheet"]
	if rels == nil || len(rels.Rels) != 1 {
		t.Fatal("Worksheet missing relationship to table")
	}

	foundRel := false
	for _, rel := range rels.Rels {
		if rel.Type == "http://schemas.openxmlformats.org/officeDocument/2006/relationships/table" &&
			rel.Target == "../tables/table1.xml" {
			foundRel = true
			break
		}
	}
	if !foundRel {
		t.Error("Table relationship not found or incorrect")
	}
}

func TestWorksheetOrder(t *testing.T) {
	ws := &xmlstructs.Worksheet{
		SheetPr:               &xmlstructs.SheetPr{},
		Dimension:             &xmlstructs.Dimension{Ref: "A1:B2"},
		SheetViews:            &xmlstructs.SheetViews{Items: []xmlstructs.SheetView{{}}},
		SheetFormatPr:         &xmlstructs.SheetFormatPr{DefaultRowHeight: 15},
		Cols:                  &xmlstructs.Cols{Items: []xmlstructs.Col{{Min: 1, Max: 1, Width: 10}}},
		SheetData:             xmlstructs.SheetData{Rows: []xmlstructs.Row{{R: 1}}},
		SheetProtection:       &xmlstructs.SheetProtection{Password: "abc"},
		AutoFilter:            &xmlstructs.AutoFilter{Ref: "A1:B1"},
		MergeCells:            &xmlstructs.MergeCells{Count: 1, Items: []xmlstructs.MergeCell{{Ref: "A1:A2"}}},
		ConditionalFormatting: []xmlstructs.ConditionalFormatting{{Sqref: "A1"}},
		DataValidations:       &xmlstructs.DataValidations{Count: 1},
		Hyperlinks:            &xmlstructs.Hyperlinks{Items: []xmlstructs.Hyperlink{{Ref: "B2"}}},
		PageMargins:           &xmlstructs.PageMargins{},
		PageSetup:             &xmlstructs.PageSetup{},
		HeaderFooter:          &xmlstructs.HeaderFooter{},
		Drawing:               &xmlstructs.WsDrawing{RID: "rId1"},
		TableParts:            &xmlstructs.TableParts{Count: 1},
	}

	data, err := xml.Marshal(ws)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	xmlStr := string(data)

	// Check order of occurrences
	expectedOrder := []string{
		"sheetPr",
		"dimension",
		"sheetViews",
		"sheetFormatPr",
		"cols",
		"sheetData",
		"sheetProtection",
		"autoFilter",
		"mergeCells",
		"conditionalFormatting",
		"dataValidations",
		"hyperlinks",
		"pageMargins",
		"pageSetup",
		"headerFooter",
		"drawing",
		"tableParts",
	}

	lastIdx := -1
	for _, tag := range expectedOrder {
		idx := strings.Index(xmlStr, "<"+tag)
		if idx == -1 {
			t.Errorf("Tag %s not found in XML", tag)
			continue
		}
		if idx < lastIdx {
			t.Errorf("Tag %s is out of order", tag)
		}
		lastIdx = idx
	}
}

func TestWorkbookOrder(t *testing.T) {
	wb := &xmlstructs.Workbook{
		WorkbookPr:         &xmlstructs.WorkbookPr{Date1904: 0},
		WorkbookProtection: &xmlstructs.WorkbookProtection{WorkbookPassword: "abc"},
		WorkbookViews:      &xmlstructs.WorkbookViews{Items: []xmlstructs.WorkbookView{{}}},
		Sheets:             []xmlstructs.Sheet{{Name: "Sheet1", SheetID: "1", RID: "rId1"}},
		DefinedNames:       &xmlstructs.DefinedNames{Items: []xmlstructs.DefinedName{{Name: "test", Ref: "A1"}}},
		CalcPr:             &xmlstructs.CalcPr{FullCalcOnLoad: 1},
	}

	data, err := xml.Marshal(wb)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	xmlStr := string(data)

	expectedOrder := []string{
		"workbookPr",
		"workbookProtection",
		"bookViews",
		"sheets",
		"definedNames",
		"calcPr",
	}

	lastIdx := -1
	for _, tag := range expectedOrder {
		idx := strings.Index(xmlStr, "<"+tag)
		if idx == -1 {
			t.Errorf("Tag %s not found in XML", tag)
			continue
		}
		if idx < lastIdx {
			t.Errorf("Tag %s is out of order", tag)
		}
		lastIdx = idx
	}
}
