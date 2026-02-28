package excel

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

type cellProcessor struct{ *state }

func (e *cellProcessor) setCellValue(sheet, axis string, value any) error {
	if sheet == "" || axis == "" {
		return fmt.Errorf("sheet and axis cannot be empty")
	}
	targetCell, err := e.getOrCreateCell(sheet, axis)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case string:
		if e.sharedStrings == nil {
			e.sharedStrings = &xmlstructs.SharedStrings{SI: make([]xmlstructs.SI, 0)}
		}
		var idx int
		if cachedIdx, exists := e.sharedStringsIndex[v]; exists {
			idx = cachedIdx
		} else {
			idx = len(e.sharedStrings.SI)
			e.sharedStrings.SI = append(e.sharedStrings.SI, xmlstructs.SI{T: v})
			e.sharedStrings.Count++
			e.sharedStrings.Unique++
			e.sharedStringsIndex[v] = idx
		}
		targetCell.T = "s"
		targetCell.V = strconv.Itoa(idx)
		targetCell.F = nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float64, float32:
		targetCell.T = "n"
		targetCell.V = fmt.Sprintf("%v", v)
		targetCell.F = nil
	case bool:
		targetCell.T = "b"
		if v {
			targetCell.V = "1"
		} else {
			targetCell.V = "0"
		}
		targetCell.F = nil
	case time.Time:
		excelBaseDate := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		days := v.Sub(excelBaseDate).Hours() / 24
		targetCell.T = "n" // Dates are numbers in Excel
		targetCell.V = fmt.Sprintf("%v", days)
		targetCell.F = nil
	case []document.TextSpan:
		targetCell.T = "inlineStr"
		targetCell.IS = &xmlstructs.Rst{
			R: make([]xmlstructs.Run, 0, len(v)),
		}
		for _, span := range v {
			run := xmlstructs.Run{T: span.Text}
			if span.Style.Bold || span.Style.Italic || span.Style.Size > 0 || span.Style.Color != "" || span.Style.Font != "" {
				run.RPr = &xmlstructs.RPr{}
				if span.Style.Bold {
					run.RPr.Bold = &struct{}{}
				}
				if span.Style.Italic {
					run.RPr.Italic = &struct{}{}
				}
				if span.Style.Size > 0 {
					run.RPr.Size = &xmlstructs.ValInt{Val: span.Style.Size}
				}
				if span.Style.Color != "" {
					run.RPr.Color = &xmlstructs.Color{RGB: span.Style.Color}
				}
				if span.Style.Font != "" {
					run.RPr.RFont = &xmlstructs.ValString{Val: span.Style.Font}
				}
			}
			targetCell.IS.R = append(targetCell.IS.R, run)
		}
		targetCell.V = ""
		targetCell.F = nil
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}

	return nil
}

func (e *cellProcessor) setCellFormula(sheet, axis string, formula string) error {
	if sheet == "" || axis == "" {
		return fmt.Errorf("sheet and axis cannot be empty")
	}
	targetCell, err := e.getOrCreateCell(sheet, axis)
	if err != nil {
		return err
	}

	if formula != "" && formula[0] == '=' {
		formula = formula[1:]
	}
	targetCell.F = &formula
	targetCell.T = ""
	targetCell.V = "0"
	return nil
}

func (e *cellProcessor) setCellHyperlink(sheet, axis string, url string) error {
	if sheet == "" || axis == "" || url == "" {
		return fmt.Errorf("parameters cannot be empty")
	}
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if ws.Hyperlinks == nil {
		ws.Hyperlinks = &xmlstructs.Hyperlinks{}
	}

	if e.sheetRels == nil {
		e.sheetRels = make(map[string]*xmlstructs.Relationships)
	}
	rels, ok := e.sheetRels[sheet]
	if !ok {
		rels = &xmlstructs.Relationships{}
		e.sheetRels[sheet] = rels
	}

	relID := rels.AddRelationshipMode("http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink", url, "External")

	ws.Hyperlinks.Items = append(ws.Hyperlinks.Items, xmlstructs.Hyperlink{
		Ref: axis,
		RID: relID,
	})

	return nil
}

func (e *cellProcessor) getCellValue(sheet, axis string) (string, error) {
	ws, ok := e.sheets[sheet]
	if !ok {
		return "", fmt.Errorf("%w: %s", document.ErrSheetNotFound, sheet)
	}

	if e.cellCache[sheet] != nil {
		if cell, ok := e.cellCache[sheet][axis]; ok {
			return e.resolveValue(*cell), nil
		}
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
