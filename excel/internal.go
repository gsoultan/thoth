package excel

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

func (e *state) loadCore(ctx context.Context) error {
	// 0. Load Content Types
	var ct xmlstructs.ContentTypes
	if err := e.loadXML("[Content_Types].xml", &ct); err == nil {
		e.contentTypes = &ct
	} else {
		e.contentTypes = xmlstructs.NewContentTypes()
	}

	// 1. Load root relationships to find workbook
	var rootRels xmlstructs.Relationships
	if err := e.loadXML("_rels/.rels", &rootRels); err != nil {
		return fmt.Errorf("load root rels: %w", err)
	}
	e.rootRels = &rootRels
	e.rootRelsPath = "_rels/.rels"

	workbookPath := rootRels.TargetByType("http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument")
	if workbookPath == "" {
		workbookPath = "xl/workbook.xml" // fallback
	} else if strings.HasPrefix(workbookPath, "/") {
		workbookPath = workbookPath[1:]
	}

	// 2. Load workbook
	var wb xmlstructs.Workbook
	if err := e.loadXML(workbookPath, &wb); err != nil {
		return fmt.Errorf("load workbook: %w", err)
	}
	e.workbook = &wb

	// 3. Load workbook relationships to find other parts
	wbRelsPath := strings.Replace(workbookPath, "workbook.xml", "_rels/workbook.xml.rels", 1)
	e.wbRelsPath = wbRelsPath
	var wbRels xmlstructs.Relationships
	e.loadXML(wbRelsPath, &wbRels)
	e.workbookRels = &wbRels

	// Shared Strings
	ssPath := wbRels.TargetByType("http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings")
	if ssPath != "" {
		if !strings.HasPrefix(ssPath, "/") {
			ssPath = "xl/" + ssPath // simplified base path handling
		} else {
			ssPath = ssPath[1:]
		}
		var ss xmlstructs.SharedStrings
		if err := e.loadXML(ssPath, &ss); err == nil {
			e.sharedStrings = &ss
			for i, si := range ss.SI {
				e.sharedStringsIndex[si.T] = i
			}
		}
	}

	// Core Properties
	cpPath := rootRels.TargetByType("http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties")
	if cpPath != "" {
		if strings.HasPrefix(cpPath, "/") {
			cpPath = cpPath[1:]
		}
		var cp xmlstructs.CoreProperties
		if err := e.loadXML(cpPath, &cp); err == nil {
			e.coreProperties = &cp
		}
	}

	// Styles
	stylesPath := wbRels.TargetByType("http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles")
	if stylesPath != "" {
		if !strings.HasPrefix(stylesPath, "/") {
			stylesPath = "xl/" + stylesPath
		} else {
			stylesPath = stylesPath[1:]
		}
		var s xmlstructs.Styles
		if err := e.loadXML(stylesPath, &s); err == nil {
			e.styles = &s
		}
	}

	// Load sheets
	for _, s := range e.workbook.Sheets {
		target := ""
		for _, rel := range e.workbookRels.Rels {
			if rel.ID == s.RID {
				target = rel.Target
				break
			}
		}

		if target == "" {
			target = fmt.Sprintf("worksheets/sheet%s.xml", s.SheetID)
		}

		var path string
		if strings.HasPrefix(target, "/") {
			path = target[1:]
		} else {
			path = "xl/" + target
		}

		var ws xmlstructs.Worksheet
		if err := e.loadXML(path, &ws); err != nil {
			continue
		}
		e.sheets[s.Name] = &ws

		// Load sheet rels
		relPath := ""
		if strings.Contains(path, "/") {
			lastSlash := strings.LastIndex(path, "/")
			relPath = path[:lastSlash] + "/_rels/" + path[lastSlash+1:] + ".rels"
		} else {
			relPath = "_rels/" + path + ".rels"
		}

		var wRels xmlstructs.Relationships
		if err := e.loadXML(relPath, &wRels); err == nil {
			e.sheetRels[s.Name] = &wRels
		} else {
			e.sheetRels[s.Name] = &xmlstructs.Relationships{}
		}
	}

	return nil
}

func (e *state) loadXML(name string, target any) error {
	for _, f := range e.reader.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			return xml.NewDecoder(rc).Decode(target)
		}
	}
	return fmt.Errorf("file %s not found in zip", name)
}

func (e *state) writeXML(zw *zip.Writer, name string, data any) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	fmt.Fprint(w, xml.Header)
	return xml.NewEncoder(w).Encode(data)
}

func (e *state) copyFile(f *zip.File, zw *zip.Writer) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	w, err := zw.Create(f.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, rc)
	return err
}

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

func (e *state) getOrCreateCell(sheet, axis string) (*xmlstructs.Cell, error) {
	ws, ok := e.sheets[sheet]
	if !ok {
		return nil, fmt.Errorf("%w: %s", document.ErrSheetNotFound, sheet)
	}

	rowIdx, err := getRowFromAxis(axis)
	if err != nil {
		return nil, err
	}

	if e.cellCache[sheet] == nil {
		e.cellCache[sheet] = make(map[string]*xmlstructs.Cell)
	}
	if cell, ok := e.cellCache[sheet][axis]; ok {
		return cell, nil
	}

	var targetRow *xmlstructs.Row
	insertIdx := -1
	for i := range ws.SheetData.Rows {
		if ws.SheetData.Rows[i].R == rowIdx {
			targetRow = &ws.SheetData.Rows[i]
			break
		}
		if ws.SheetData.Rows[i].R > rowIdx {
			insertIdx = i
			break
		}
	}

	if targetRow == nil {
		newRow := xmlstructs.Row{R: rowIdx}
		if insertIdx == -1 {
			ws.SheetData.Rows = append(ws.SheetData.Rows, newRow)
			targetRow = &ws.SheetData.Rows[len(ws.SheetData.Rows)-1]
		} else {
			ws.SheetData.Rows = append(ws.SheetData.Rows[:insertIdx], append([]xmlstructs.Row{newRow}, ws.SheetData.Rows[insertIdx:]...)...)
			targetRow = &ws.SheetData.Rows[insertIdx]
		}
	}

	col := getColumnFromAxis(axis)
	cellInsertIdx := -1
	for i := range targetRow.Cells {
		currCol := getColumnFromAxis(targetRow.Cells[i].R)
		if targetRow.Cells[i].R == axis {
			cell := &targetRow.Cells[i]
			e.cellCache[sheet][axis] = cell
			return cell, nil
		}
		if compareColumns(currCol, col) > 0 {
			cellInsertIdx = i
			break
		}
	}

	newCell := xmlstructs.Cell{R: axis}
	if cellInsertIdx == -1 {
		targetRow.Cells = append(targetRow.Cells, newCell)
		cell := &targetRow.Cells[len(targetRow.Cells)-1]
		e.cellCache[sheet][axis] = cell
		return cell, nil
	}

	targetRow.Cells = append(targetRow.Cells[:cellInsertIdx], append([]xmlstructs.Cell{newCell}, targetRow.Cells[cellInsertIdx:]...)...)
	cell := &targetRow.Cells[cellInsertIdx]
	e.cellCache[sheet][axis] = cell
	return cell, nil
}

func (e *state) resolveValue(cell xmlstructs.Cell) string {
	if cell.IS != nil {
		if cell.IS.T != "" {
			return cell.IS.T
		}
		var sb strings.Builder
		for _, r := range cell.IS.R {
			sb.WriteString(r.T)
		}
		return sb.String()
	}
	if cell.T == "s" {
		idx, err := strconv.Atoi(cell.V)
		if err == nil && e.sharedStrings != nil && idx >= 0 && idx < len(e.sharedStrings.SI) {
			return e.sharedStrings.SI[idx].T
		}
	}
	return cell.V
}

func getColumnFromAxis(axis string) string {
	idx := strings.IndexFunc(axis, func(r rune) bool {
		return r >= '0' && r <= '9'
	})
	if idx == -1 {
		return axis
	}
	return axis[:idx]
}

func compareColumns(a, b string) int {
	if len(a) != len(b) {
		return len(a) - len(b)
	}
	return strings.Compare(a, b)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func getRowFromAxis(axis string) (int, error) {
	idx := strings.IndexFunc(axis, func(r rune) bool {
		return r >= '0' && r <= '9'
	})
	if idx == -1 || idx == 0 {
		return 0, fmt.Errorf("invalid axis: %s", axis)
	}
	prefix := axis[:idx]
	for _, r := range prefix {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return 0, fmt.Errorf("invalid axis prefix: %s", axis)
		}
	}
	return strconv.Atoi(axis[idx:])
}
