package excel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

type sheetProcessor struct{ *state }

func (e *sheetProcessor) addSheet(name string) error {
	if name == "" {
		return fmt.Errorf("sheet name cannot be empty")
	}
	if len(name) > 31 {
		return fmt.Errorf("sheet name cannot exceed 31 characters")
	}
	if e.workbook == nil {
		e.workbook = &xmlstructs.Workbook{
			Sheets: make([]xmlstructs.Sheet, 0),
		}
	}

	for _, s := range e.workbook.Sheets {
		if s.Name == name {
			return fmt.Errorf("sheet %s already exists", name)
		}
	}

	maxSheetID := 0
	for _, s := range e.workbook.Sheets {
		sid, _ := strconv.Atoi(s.SheetID)
		maxSheetID = max(maxSheetID, sid)
	}
	sheetID := maxSheetID + 1
	rID := fmt.Sprintf("rId%d", sheetID+100)

	e.workbook.Sheets = append(e.workbook.Sheets, xmlstructs.Sheet{
		Name:    name,
		SheetID: strconv.Itoa(sheetID),
		RID:     rID,
	})

	if e.workbookRels == nil {
		e.workbookRels = &xmlstructs.Relationships{}
	}
	target := fmt.Sprintf("worksheets/sheet%d.xml", sheetID)
	e.workbookRels.Rels = append(e.workbookRels.Rels, xmlstructs.Relationship{
		ID:     rID,
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet",
		Target: target,
	})

	tabSelected := 0
	if len(e.workbook.Sheets) == 1 {
		tabSelected = 1
	}

	e.sheets[name] = &xmlstructs.Worksheet{
		XMLNS_R: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
		SheetViews: &xmlstructs.SheetViews{
			Items: []xmlstructs.SheetView{
				{WorkbookViewID: 0, TabSelected: tabSelected},
			},
		},
		SheetData: xmlstructs.SheetData{
			Rows: make([]xmlstructs.Row, 0),
		},
	}

	return nil
}

func (e *sheetProcessor) mergeCells(sheet, hRange string) error {
	if sheet == "" || hRange == "" {
		return fmt.Errorf("sheet and range cannot be empty")
	}
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

func (e *sheetProcessor) setColumnWidth(sheet string, col int, width float64) error {
	if sheet == "" || col < 1 || width < 0 {
		return fmt.Errorf("invalid parameters for column width")
	}
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.Cols == nil {
		ws.Cols = &xmlstructs.Cols{Items: make([]xmlstructs.Col, 0)}
	}

	for i := range ws.Cols.Items {
		if ws.Cols.Items[i].Min <= col && ws.Cols.Items[i].Max >= col {
			if ws.Cols.Items[i].Min == col && ws.Cols.Items[i].Max == col {
				ws.Cols.Items[i].Width = width
				ws.Cols.Items[i].CustomWidth = 1
				return nil
			}
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

func (e *sheetProcessor) setRowHeight(sheet string, row int, height float64) error {
	if sheet == "" || row < 1 || height < 0 {
		return fmt.Errorf("invalid parameters for row height")
	}
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	for i := range ws.SheetData.Rows {
		if ws.SheetData.Rows[i].R == row {
			ws.SheetData.Rows[i].Ht = height
			ws.SheetData.Rows[i].CustomHeight = 1
			return nil
		}
	}

	// Not found, we can't create it here because it might be empty
	// but we could pre-create rows if we want to set height for an empty row.
	// However, usually we set height for rows with data.
	// But let's support empty rows too.
	newRow := xmlstructs.Row{
		R:            row,
		Ht:           height,
		CustomHeight: 1,
	}
	// Insert in correct position
	insertIdx := -1
	for i := range ws.SheetData.Rows {
		if ws.SheetData.Rows[i].R > row {
			insertIdx = i
			break
		}
	}
	if insertIdx == -1 {
		ws.SheetData.Rows = append(ws.SheetData.Rows, newRow)
	} else {
		ws.SheetData.Rows = append(ws.SheetData.Rows[:insertIdx], append([]xmlstructs.Row{newRow}, ws.SheetData.Rows[insertIdx:]...)...)
	}

	return nil
}

func (e *sheetProcessor) autoFilter(sheet, ref string) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}
	ws.AutoFilter = &xmlstructs.AutoFilter{Ref: ref}
	return nil
}

func (e *sheetProcessor) freezePanes(sheet string, col, row int) error {
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

func (e *sheetProcessor) setNamedRange(name, ref string) error {
	if e.workbook == nil {
		return fmt.Errorf("workbook not initialized")
	}

	if e.workbook.DefinedNames == nil {
		e.workbook.DefinedNames = &xmlstructs.DefinedNames{Items: make([]xmlstructs.DefinedName, 0)}
	}

	// Check if name already exists
	for i, dn := range e.workbook.DefinedNames.Items {
		if dn.Name == name {
			e.workbook.DefinedNames.Items[i].Ref = ref
			return nil
		}
	}

	e.workbook.DefinedNames.Items = append(e.workbook.DefinedNames.Items, xmlstructs.DefinedName{
		Name: name,
		Ref:  ref,
	})

	return nil
}

func (e *sheetProcessor) setPageSettings(sheet string, settings document.PageSettings) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	ws.PageMargins = &xmlstructs.PageMargins{
		Top:    settings.Margins.Top / 72.0, // Points to inches
		Bottom: settings.Margins.Bottom / 72.0,
		Left:   settings.Margins.Left / 72.0,
		Right:  settings.Margins.Right / 72.0,
		Header: 0.3,
		Footer: 0.3,
	}

	orientation := "portrait"
	if settings.Orientation == document.OrientationLandscape {
		orientation = "landscape"
	}

	ws.PageSetup = &xmlstructs.PageSetup{
		Orientation: orientation,
		PaperSize:   paperSizeToInt(settings.PaperType),
	}

	return nil
}

func (e *sheetProcessor) protect(sheet string, password string) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	ws.SheetProtection = &xmlstructs.SheetProtection{
		Sheet:          1,
		Objects:        1,
		Scenarios:      1,
		SelectLocked:   1,
		SelectUnlocked: 1,
	}

	if password != "" {
		ws.SheetProtection.Password = excelPasswordHash(password)
	}

	return nil
}

func (e *sheetProcessor) groupRows(sheet string, start, end int, level int) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.SheetPr == nil {
		ws.SheetPr = &xmlstructs.SheetPr{OutlinePr: &xmlstructs.OutlinePr{SummaryBelow: 1}}
	} else if ws.SheetPr.OutlinePr == nil {
		ws.SheetPr.OutlinePr = &xmlstructs.OutlinePr{SummaryBelow: 1}
	} else {
		ws.SheetPr.OutlinePr.SummaryBelow = 1
	}

	if ws.SheetFormatPr == nil {
		ws.SheetFormatPr = &xmlstructs.SheetFormatPr{DefaultRowHeight: 15.0}
	}
	if uint8(level) > ws.SheetFormatPr.OutlineLevelRow {
		ws.SheetFormatPr.OutlineLevelRow = uint8(level)
	}

	for r := start; r <= end; r++ {
		found := false
		for i := range ws.SheetData.Rows {
			if ws.SheetData.Rows[i].R == r {
				ws.SheetData.Rows[i].OutlineLevel = uint8(level)
				found = true
				break
			}
		}
		if !found {
			// Create empty row to hold grouping
			newRow := xmlstructs.Row{
				R:            r,
				OutlineLevel: uint8(level),
			}
			// Insert in correct position
			insertIdx := -1
			for i := range ws.SheetData.Rows {
				if ws.SheetData.Rows[i].R > r {
					insertIdx = i
					break
				}
			}
			if insertIdx == -1 {
				ws.SheetData.Rows = append(ws.SheetData.Rows, newRow)
			} else {
				ws.SheetData.Rows = append(ws.SheetData.Rows[:insertIdx], append([]xmlstructs.Row{newRow}, ws.SheetData.Rows[insertIdx:]...)...)
			}
		}
	}

	return nil
}

func (e *sheetProcessor) groupCols(sheet string, start, end int, level int) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.SheetPr == nil {
		ws.SheetPr = &xmlstructs.SheetPr{OutlinePr: &xmlstructs.OutlinePr{SummaryRight: 1}}
	} else if ws.SheetPr.OutlinePr == nil {
		ws.SheetPr.OutlinePr = &xmlstructs.OutlinePr{SummaryRight: 1}
	} else {
		ws.SheetPr.OutlinePr.SummaryRight = 1
	}

	if ws.SheetFormatPr == nil {
		ws.SheetFormatPr = &xmlstructs.SheetFormatPr{DefaultRowHeight: 15.0}
	}
	if uint8(level) > ws.SheetFormatPr.OutlineLevelCol {
		ws.SheetFormatPr.OutlineLevelCol = uint8(level)
	}

	if ws.Cols == nil {
		ws.Cols = &xmlstructs.Cols{Items: make([]xmlstructs.Col, 0)}
	}

	ws.Cols.Items = append(ws.Cols.Items, xmlstructs.Col{
		Min:          start,
		Max:          end,
		OutlineLevel: uint8(level),
		Width:        9.14, // default width if not set
	})

	return nil
}

func excelPasswordHash(password string) string {
	var hash uint16
	if len(password) > 0 {
		for i := len(password) - 1; i >= 0; i-- {
			hash = ((hash >> 14) & 0x01) | ((hash << 1) & 0x7fff)
			hash ^= uint16(password[i])
		}
		hash = ((hash >> 14) & 0x01) | ((hash << 1) & 0x7fff)
		hash ^= uint16(len(password))
		hash ^= 0xCE4B
	}
	return fmt.Sprintf("%X", hash)
}

func paperSizeToInt(p document.PaperType) int {
	switch p {
	case document.PaperA4:
		return 9
	case document.PaperLetter:
		return 1
	case document.PaperLegal:
		return 5
	default:
		return 9 // default A4
	}
}

func (e *sheetProcessor) setHeader(sheet string, text string) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.HeaderFooter == nil {
		ws.HeaderFooter = &xmlstructs.HeaderFooter{}
	}
	ws.HeaderFooter.OddHeader = text
	return nil
}

func (e *sheetProcessor) setFooter(sheet string, text string) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.HeaderFooter == nil {
		ws.HeaderFooter = &xmlstructs.HeaderFooter{}
	}
	ws.HeaderFooter.OddFooter = text
	return nil
}

func (e *sheetProcessor) setDataValidation(sheet, ref string, options ...string) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.DataValidations == nil {
		ws.DataValidations = &xmlstructs.DataValidations{Items: make([]xmlstructs.DataValidation, 0)}
	}

	dv := xmlstructs.DataValidation{
		Type:         "list",
		AllowBlank:   1,
		ShowInputMsg: 1,
		ShowErrorMsg: 1,
		Sqref:        ref,
		Formula1:     fmt.Sprintf("\"%s\"", strings.Join(options, ",")),
	}

	ws.DataValidations.Items = append(ws.DataValidations.Items, dv)
	ws.DataValidations.Count = len(ws.DataValidations.Items)

	return nil
}

func (e *sheetProcessor) SetPageSettings(settings document.PageSettings) error {
	if e.workbook == nil {
		e.workbook = &xmlstructs.Workbook{
			Sheets: make([]xmlstructs.Sheet, 0),
		}
	}

	for _, ws := range e.sheets {
		if ws.PageSetup == nil {
			ws.PageSetup = &xmlstructs.PageSetup{}
		}

		if settings.Orientation == document.OrientationLandscape {
			ws.PageSetup.Orientation = "landscape"
		} else {
			ws.PageSetup.Orientation = "portrait"
		}

		switch settings.PaperType {
		case document.PaperA4:
			ws.PageSetup.PaperSize = 9
		case document.PaperLetter:
			ws.PageSetup.PaperSize = 1
		}

		ws.PageMargins = &xmlstructs.PageMargins{
			Top:    settings.Margins.Top / 72.0,
			Bottom: settings.Margins.Bottom / 72.0,
			Left:   settings.Margins.Left / 72.0,
			Right:  settings.Margins.Right / 72.0,
		}
	}
	return nil
}

func (e *sheetProcessor) setPrintArea(sheet, ref string) error {
	idx := e.getSheetIndex(sheet)
	if idx == -1 {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if e.workbook.DefinedNames == nil {
		e.workbook.DefinedNames = &xmlstructs.DefinedNames{Items: make([]xmlstructs.DefinedName, 0)}
	}

	// Reference must be absolute and include sheet name
	fullRef := fmt.Sprintf("'%s'!%s", sheet, strings.ReplaceAll(ref, "$", ""))
	// Excel usually expects absolute references with $ in Print Area
	// But let's assume user might provide A1:B2 or $A$1:$B$2
	// Actually we should probably normalize it to $A$1:$B$2 format if possible.
	// For now let's just use what user provides but prefix with sheet name.

	name := "_xlnm.Print_Area"
	for i, dn := range e.workbook.DefinedNames.Items {
		if dn.Name == name && dn.LocalSheetID != nil && *dn.LocalSheetID == idx {
			e.workbook.DefinedNames.Items[i].Ref = fullRef
			return nil
		}
	}

	e.workbook.DefinedNames.Items = append(e.workbook.DefinedNames.Items, xmlstructs.DefinedName{
		Name:         name,
		LocalSheetID: &idx,
		Ref:          fullRef,
	})
	return nil
}

func (e *sheetProcessor) setPrintTitles(sheet, rowRef, colRef string) error {
	idx := e.getSheetIndex(sheet)
	if idx == -1 {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if e.workbook.DefinedNames == nil {
		e.workbook.DefinedNames = &xmlstructs.DefinedNames{Items: make([]xmlstructs.DefinedName, 0)}
	}

	var titles []string
	if colRef != "" {
		titles = append(titles, fmt.Sprintf("'%s'!%s", sheet, colRef))
	}
	if rowRef != "" {
		titles = append(titles, fmt.Sprintf("'%s'!%s", sheet, rowRef))
	}

	if len(titles) == 0 {
		return nil
	}

	fullRef := strings.Join(titles, ",")
	name := "_xlnm.Print_Titles"

	for i, dn := range e.workbook.DefinedNames.Items {
		if dn.Name == name && dn.LocalSheetID != nil && *dn.LocalSheetID == idx {
			e.workbook.DefinedNames.Items[i].Ref = fullRef
			return nil
		}
	}

	e.workbook.DefinedNames.Items = append(e.workbook.DefinedNames.Items, xmlstructs.DefinedName{
		Name:         name,
		LocalSheetID: &idx,
		Ref:          fullRef,
	})
	return nil
}

func (e *sheetProcessor) getSheetIndex(name string) int {
	for i, s := range e.workbook.Sheets {
		if s.Name == name {
			return i
		}
	}
	return -1
}

func countColumnsInRange(ref string) int {
	parts := strings.Split(ref, ":")
	if len(parts) != 2 {
		return 1
	}
	startCol := getColumnFromAxis(strings.ReplaceAll(parts[0], "$", ""))
	endCol := getColumnFromAxis(strings.ReplaceAll(parts[1], "$", ""))

	start := colToNum(startCol)
	end := colToNum(endCol)

	if end < start {
		return 1
	}
	return end - start + 1
}

func colToNum(col string) int {
	num := 0
	col = strings.ToUpper(col)
	for _, r := range col {
		if r >= 'A' && r <= 'Z' {
			num = num*26 + int(r-'A') + 1
		}
	}
	return num
}

func numToCol(num int) string {
	col := ""
	for num > 0 {
		num--
		col = string(rune('A'+(num%26))) + col
		num /= 26
	}
	return col
}

func (e *sheetProcessor) addTable(sheet, ref, name string) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	tableID := len(e.tables) + 1
	tablePath := fmt.Sprintf("xl/tables/table%d.xml", tableID)

	table := &xmlstructs.Table{
		ID:          tableID,
		Name:        name,
		DisplayName: name,
		Ref:         ref,
		AutoFilter:  &xmlstructs.AutoFilter{Ref: ref},
		TableStyleInfo: &xmlstructs.TableStyleInfo{
			Name:           "TableStyleMedium2",
			ShowRowStripes: 1,
		},
	}

	// Column detection and header extraction
	numCols := countColumnsInRange(ref)
	table.TableColumns.Items = make([]xmlstructs.TableColumn, numCols)

	parts := strings.Split(ref, ":")
	startAxis := parts[0]
	startAxisNorm := strings.ReplaceAll(startAxis, "$", "")
	rowStr := strings.TrimLeft(startAxisNorm, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	row, _ := strconv.Atoi(rowStr)
	startCol := getColumnFromAxis(startAxisNorm)
	startNum := colToNum(startCol)

	for i := 0; i < numCols; i++ {
		colName := numToCol(startNum + i)
		axis := fmt.Sprintf("%s%d", colName, row)
		header := fmt.Sprintf("Column%d", i+1)

		// Try to extract header from cell
		if val, err := e.processor().getCellValue(sheet, axis); err == nil && val != "" {
			header = val
		}

		table.TableColumns.Items[i] = xmlstructs.TableColumn{
			ID:   i + 1,
			Name: header,
		}
	}
	table.TableColumns.Count = numCols

	e.tables[tablePath] = table

	if ws.TableParts == nil {
		ws.TableParts = &xmlstructs.TableParts{Items: make([]xmlstructs.TablePart, 0)}
	}

	if e.sheetRels[sheet] == nil {
		e.sheetRels[sheet] = &xmlstructs.Relationships{}
	}
	sRels := e.sheetRels[sheet]

	relPath := fmt.Sprintf("../tables/table%d.xml", tableID)
	rID := sRels.AddRelationship("http://schemas.openxmlformats.org/officeDocument/2006/relationships/table", relPath)

	ws.TableParts.Items = append(ws.TableParts.Items, xmlstructs.TablePart{RID: rID})
	ws.TableParts.Count = len(ws.TableParts.Items)

	return nil
}
