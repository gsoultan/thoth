package excel

import (
	"fmt"
	"strconv"

	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// sheetManager handles operations related to sheets.
type sheetManager struct{ *state }

func (e *state) GetSheets() ([]string, error) {
	if e.workbook == nil {
		return nil, fmt.Errorf("workbook not loaded")
	}
	sheets := make([]string, 0, len(e.workbook.Sheets))
	for _, s := range e.workbook.Sheets {
		sheets = append(sheets, s.Name)
	}
	return sheets, nil
}

func (e *sheetManager) addSheet(name string) error {
	if e.workbook == nil {
		e.workbook = &xmlstructs.Workbook{
			Sheets: make([]xmlstructs.Sheet, 0),
		}
	}

	// Check if sheet exists
	for _, s := range e.workbook.Sheets {
		if s.Name == name {
			return fmt.Errorf("sheet %s already exists", name)
		}
	}

	// 1. Generate new sheetId and rId
	maxSheetID := 0
	for _, s := range e.workbook.Sheets {
		sid, _ := strconv.Atoi(s.SheetID)
		if sid > maxSheetID {
			maxSheetID = sid
		}
	}
	sheetID := maxSheetID + 1
	rID := fmt.Sprintf("rId%d", sheetID+100)

	// 2. Add to workbook
	e.workbook.Sheets = append(e.workbook.Sheets, xmlstructs.Sheet{
		Name:    name,
		SheetID: strconv.Itoa(sheetID),
		RID:     rID,
	})

	// 3. Add to relationships
	if e.workbookRels == nil {
		e.workbookRels = &xmlstructs.Relationships{}
	}
	target := fmt.Sprintf("worksheets/sheet%d.xml", sheetID)
	e.workbookRels.Rels = append(e.workbookRels.Rels, xmlstructs.Relationship{
		ID:     rID,
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet",
		Target: target,
	})

	// 4. Initialize worksheet
	e.sheets[name] = &xmlstructs.Worksheet{
		SheetData: xmlstructs.SheetData{
			Rows: make([]xmlstructs.Row, 0),
		},
	}

	return nil
}

func (e *sheetManager) mergeCells(sheet, hRange string) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.MergeCells == nil {
		ws.MergeCells = &xmlstructs.MergeCells{Items: make([]xmlstructs.MergeCell, 0)}
	}

	ws.MergeCells.Items = append(ws.MergeCells.Items, xmlstructs.MergeCell{Ref: hRange})
	ws.MergeCells.Count = len(ws.MergeCells.Items)

	return nil
}

func (e *sheetManager) setColumnWidth(sheet string, col int, width float64) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.Cols == nil {
		ws.Cols = &xmlstructs.Cols{Items: make([]xmlstructs.Col, 0)}
	}

	// Check if column already has setting
	for i := range ws.Cols.Items {
		if ws.Cols.Items[i].Min <= col && ws.Cols.Items[i].Max >= col {
			if ws.Cols.Items[i].Min == col && ws.Cols.Items[i].Max == col {
				ws.Cols.Items[i].Width = width
				ws.Cols.Items[i].CustomWidth = 1
				return nil
			}
			// Column is part of a range. We'd need to split the range.
			// For simplicity in this advanced scenario, we'll just add it as a new rule if not exact.
		}
	}

	ws.Cols.Items = append(ws.Cols.Items, xmlstructs.Col{
		Min:         col,
		Max:         col,
		Width:       width,
		CustomWidth: 1,
	})

	return nil
}

func (e *sheetManager) autoFilter(sheet, ref string) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}
	ws.AutoFilter = &xmlstructs.AutoFilter{Ref: ref}
	return nil
}

func (e *sheetManager) freezePanes(sheet string, col, row int) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.SheetViews == nil {
		ws.SheetViews = &xmlstructs.SheetViews{
			Items: []xmlstructs.SheetView{{WorkbookViewID: 0, TabSelected: 1}},
		}
	}

	activePane := "bottomRight"
	if col > 0 && row == 0 {
		activePane = "topRight"
	} else if col == 0 && row > 0 {
		activePane = "bottomLeft"
	}

	ws.SheetViews.Items[0].Pane = &xmlstructs.Pane{
		XSplit:      col,
		YSplit:      row,
		TopLeftCell: fmt.Sprintf("%s%d", string(rune('A'+col)), row+1),
		ActivePane:  activePane,
		State:       "frozen",
	}

	return nil
}
