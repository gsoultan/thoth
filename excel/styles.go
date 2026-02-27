package excel

import (
	"fmt"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// styleManager handles operations related to cell styles.
type styleManager struct{ *state }

func (e *styleManager) setCellStyle(sheet, axis string, style document.CellStyle) error {
	if e.styles == nil {
		e.styles = &xmlstructs.Styles{
			Fonts:   xmlstructs.Fonts{Count: 1, Items: []xmlstructs.Font{{}}},
			Fills:   xmlstructs.Fills{Count: 2, Items: []xmlstructs.Fill{{PatternFill: &xmlstructs.PatternFill{PatternType: "none"}}, {PatternFill: &xmlstructs.PatternFill{PatternType: "gray125"}}}},
			Borders: xmlstructs.Borders{Count: 1, Items: []xmlstructs.Border{{}}},
			CellXfs: xmlstructs.CellXfs{Count: 1, Items: []xmlstructs.Xf{{}}},
		}
	}

	f := xmlstructs.Font{}
	if style.Bold {
		f.Bold = &struct{}{}
	}
	if style.Italic {
		f.Italic = &struct{}{}
	}
	if style.Size > 0 {
		f.Size = &xmlstructs.ValInt{Val: style.Size}
	}
	if style.Color != "" {
		f.Color = &xmlstructs.Color{RGB: style.Color}
	}

	fontID := e.styles.AddFont(f)

	fillID := 0
	if style.Background != "" {
		fill := xmlstructs.Fill{
			PatternFill: &xmlstructs.PatternFill{
				PatternType: "solid",
				FgColor:     &xmlstructs.Color{RGB: style.Background},
			},
		}
		fillID = e.styles.AddFill(fill)
	}

	borderID := 0
	if style.Border {
		border := xmlstructs.Border{
			Left:   &xmlstructs.BorderEdge{Style: "thin"},
			Right:  &xmlstructs.BorderEdge{Style: "thin"},
			Top:    &xmlstructs.BorderEdge{Style: "thin"},
			Bottom: &xmlstructs.BorderEdge{Style: "thin"},
		}
		borderID = e.styles.AddBorder(border)
	}

	numFmtID := 0
	if style.NumberFormat != "" {
		numFmtID = e.getNumFmtID(style.NumberFormat)
	}

	xf := xmlstructs.Xf{
		NumFmtID: numFmtID,
		FontID:   fontID,
		FillID:   fillID,
		BorderID: borderID,
	}
	if numFmtID > 0 {
		xf.ApplyNumberFormat = 1
	}
	if fontID > 0 {
		xf.ApplyFont = 1
	}
	if fillID > 0 {
		xf.ApplyFill = 1
	}
	if borderID > 0 {
		xf.ApplyBorder = 1
	}

	if style.Horizontal != "" || style.Vertical != "" {
		xf.ApplyAlignment = 1
		xf.Alignment = &xmlstructs.Alignment{
			Horizontal: style.Horizontal,
			Vertical:   style.Vertical,
		}
	}

	xfID := e.styles.AddXf(xf)

	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("%w: %s", document.ErrSheetNotFound, sheet)
	}

	cell, err := e.getOrCreateCell(ws, axis)
	if err != nil {
		return err
	}
	cell.S = xfID

	return nil
}

func (e *styleManager) getNumFmtID(formatCode string) int {
	// Standard IDs
	standards := map[string]int{
		"General":       0,
		"0":             1,
		"0.00":          2,
		"#,##0":         3,
		"#,##0.00":      4,
		"0%":            9,
		"0.00%":         10,
		"mm-dd-yy":      14,
		"d-mmm-yy":      15,
		"d-mmm":         16,
		"mmm-yy":        17,
		"h:mm AM/PM":    18,
		"h:mm:ss AM/PM": 19,
		"h:mm":          20,
		"h:mm:ss":       21,
		"m/d/yy h:mm":   22,
		"@":             49,
	}
	if id, ok := standards[formatCode]; ok {
		return id
	}
	return e.styles.AddNumFmt(formatCode)
}
