package excel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// cellProcessor handles operations related to cells and values.
type cellProcessor struct{ *state }

func (e *cellProcessor) setCellValue(sheet, axis string, value any) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	targetCell, err := e.getOrCreateCell(ws, axis)
	if err != nil {
		return err
	}

	// 4. Set value
	switch v := value.(type) {
	case string:
		if e.sharedStrings == nil {
			e.sharedStrings = &xmlstructs.SharedStrings{SI: make([]xmlstructs.SI, 0)}
		}
		idx := e.sharedStrings.AddString(v)
		targetCell.T = "s"
		targetCell.V = strconv.Itoa(idx)
	case int, int64, float64:
		targetCell.T = ""
		targetCell.V = fmt.Sprintf("%v", v)
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}

	return nil
}

func (e *state) getOrCreateCell(ws *xmlstructs.Worksheet, axis string) (*xmlstructs.Cell, error) {
	// 1. Resolve row index from axis
	rowIdx, err := getRowFromAxis(axis)
	if err != nil {
		return nil, err
	}

	// 2. Find or create row (ensuring they are sorted)
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

	// 3. Find or create cell (ensuring they are sorted by column too if possible)
	// For cells, it's also better to keep them sorted, but let's at least ensure unique ones.
	for i := range targetRow.Cells {
		if targetRow.Cells[i].R == axis {
			return &targetRow.Cells[i], nil
		}
	}
	targetRow.Cells = append(targetRow.Cells, xmlstructs.Cell{R: axis})
	return &targetRow.Cells[len(targetRow.Cells)-1], nil
}

func getRowFromAxis(axis string) (int, error) {
	// Proper axis parsing (e.g., "A1" -> 1, "B12" -> 12)
	idx := strings.IndexFunc(axis, func(r rune) bool {
		return r >= '0' && r <= '9'
	})
	if idx == -1 || idx == 0 {
		return 0, fmt.Errorf("invalid axis: %s", axis)
	}

	// Validate prefix is all letters
	prefix := axis[:idx]
	for _, r := range prefix {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return 0, fmt.Errorf("invalid axis prefix: %s", axis)
		}
	}

	return strconv.Atoi(axis[idx:])
}

func (e *cellProcessor) getCellValue(sheet, axis string) (string, error) {
	ws, ok := e.sheets[sheet]
	if !ok {
		return "", fmt.Errorf("%w: %s", document.ErrSheetNotFound, sheet)
	}

	for _, row := range ws.SheetData.Rows {
		for _, cell := range row.Cells {
			if cell.R == axis {
				return e.resolveValue(cell), nil
			}
		}
	}

	return "", nil
}
