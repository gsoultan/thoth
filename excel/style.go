package excel

import (
	"fmt"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

type styleProcessor struct{ *state }

func (e *styleProcessor) setConditionalFormatting(sheet, ref, ruleType, operator, formula string, style document.CellStyle) error {
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	// 1. Create DXF (Differential Formatting)
	dxf := xmlstructs.Xf{}
	if style.Bold {
		dxf.Font = &xmlstructs.Font{Bold: &struct{}{}}
	}
	if style.Color != "" {
		if dxf.Font == nil {
			dxf.Font = &xmlstructs.Font{}
		}
		dxf.Font.Color = &xmlstructs.Color{RGB: style.Color}
	}
	if style.Background != "" {
		dxf.Fill = &xmlstructs.Fill{
			PatternFill: &xmlstructs.PatternFill{
				PatternType: "solid",
				FgColor:     &xmlstructs.Color{RGB: style.Background},
			},
		}
	}

	dxfID := e.styles.AddDxf(dxf)

	// 2. Add Conditional Formatting Rule
	cf := xmlstructs.ConditionalFormatting{
		Sqref: ref,
		CfRule: []xmlstructs.CfRule{
			{
				Type:     ruleType,
				Operator: operator,
				Priority: len(ws.ConditionalFormatting) + 1,
				DxfID:    &dxfID,
				Formula:  []string{formula},
			},
		},
	}

	ws.ConditionalFormatting = append(ws.ConditionalFormatting, cf)
	return nil
}

func (e *styleProcessor) setCellStyle(sheet, axis string, style document.CellStyle) error {
	// 1. Font
	f := xmlstructs.Font{
		Size: &xmlstructs.ValInt{Val: 11},
		Name: &xmlstructs.ValString{Val: "Calibri"},
	}
	if style.Bold {
		f.Bold = new(struct{}{})
	}
	if style.Italic {
		f.Italic = new(struct{}{})
	}
	if style.Size > 0 {
		f.Size = &xmlstructs.ValInt{Val: style.Size}
	}
	if style.Color != "" {
		f.Color = &xmlstructs.Color{RGB: style.Color}
	}
	if style.Font != "" {
		f.Name = &xmlstructs.ValString{Val: style.Font}
	}
	fontID := e.getFontID(f)

	// 2. Fill
	fillID := 0
	if style.Background != "" {
		fill := xmlstructs.Fill{
			PatternFill: &xmlstructs.PatternFill{
				PatternType: "solid",
				FgColor:     &xmlstructs.Color{RGB: style.Background},
			},
		}
		fillID = e.getFillID(fill)
	}

	// 3. Border
	borderID := 0
	if style.Border || style.BorderTop || style.BorderBottom || style.BorderLeft || style.BorderRight {
		border := xmlstructs.Border{}
		borderStyle := "thin"
		if style.BorderWidth > 1 {
			borderStyle = "medium"
		}
		if style.BorderWidth > 2 {
			borderStyle = "thick"
		}

		borderColor := &xmlstructs.Color{RGB: style.BorderColor}
		if style.BorderColor == "" {
			borderColor = nil
		}

		if style.Border || style.BorderLeft {
			border.Left = xmlstructs.BorderEdge{Style: borderStyle, Color: borderColor}
		}
		if style.Border || style.BorderRight {
			border.Right = xmlstructs.BorderEdge{Style: borderStyle, Color: borderColor}
		}
		if style.Border || style.BorderTop {
			border.Top = xmlstructs.BorderEdge{Style: borderStyle, Color: borderColor}
		}
		if style.Border || style.BorderBottom {
			border.Bottom = xmlstructs.BorderEdge{Style: borderStyle, Color: borderColor}
		}
		borderID = e.getBorderID(border)
	}

	// 4. Number Format
	numFmtID := 0
	if style.NumberFormat != "" {
		numFmtID = e.getNumFmtID(style.NumberFormat)
	}

	// 5. XF
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

	if style.Horizontal != "" || style.Vertical != "" || style.Padding > 0 {
		xf.ApplyAlignment = 1
		xf.Alignment = &xmlstructs.Alignment{
			Horizontal: style.Horizontal,
			Vertical:   style.Vertical,
			Indent:     int(style.Indent),
			WrapText:   int(boolToInt(style.WrapText)),
		}
	}

	xfID := e.getXfID(xf)

	cell, err := e.getOrCreateCell(sheet, axis)
	if err != nil {
		return err
	}
	cell.S = xfID

	return nil
}

func (e *styleProcessor) getFontID(f xmlstructs.Font) int {
	key := fmt.Sprintf("b:%v|i:%v|s:%v|c:%v|n:%v",
		f.Bold != nil, f.Italic != nil,
		valOr(f.Size, 0), valOr(f.Color, ""), valOr(f.Name, ""))
	if id, ok := e.fontsIndex[key]; ok {
		return id
	}
	id := e.styles.AddFont(f)
	e.fontsIndex[key] = id
	return id
}

func (e *styleProcessor) getFillID(f xmlstructs.Fill) int {
	key := "none"
	if f.PatternFill != nil {
		key = fmt.Sprintf("t:%s|c:%s", f.PatternFill.PatternType,
			valOr(f.PatternFill.FgColor, ""))
	}
	if id, ok := e.fillsIndex[key]; ok {
		return id
	}
	id := e.styles.AddFill(f)
	e.fillsIndex[key] = id
	return id
}

func (e *styleProcessor) getBorderID(b xmlstructs.Border) int {
	key := fmt.Sprintf("l:%v|r:%v|t:%v|b:%v",
		edgeKey(&b.Left), edgeKey(&b.Right), edgeKey(&b.Top), edgeKey(&b.Bottom))
	if id, ok := e.bordersIndex[key]; ok {
		return id
	}
	id := e.styles.AddBorder(b)
	e.bordersIndex[key] = id
	return id
}

func (e *styleProcessor) getXfID(xf xmlstructs.Xf) int {
	alignKey := ""
	if xf.Alignment != nil {
		alignKey = fmt.Sprintf("|h:%s|v:%s|w:%d",
			xf.Alignment.Horizontal, xf.Alignment.Vertical, xf.Alignment.WrapText)
	}
	key := fmt.Sprintf("n:%d|f:%d|l:%d|b:%d|a:%d%s",
		xf.NumFmtID, xf.FontID, xf.FillID, xf.BorderID, xf.ApplyAlignment, alignKey)
	if id, ok := e.xfsIndex[key]; ok {
		return id
	}
	id := e.styles.AddXf(xf)
	e.xfsIndex[key] = id
	return id
}

func valOr(v any, def any) any {
	if v == nil {
		return def
	}
	switch t := v.(type) {
	case *xmlstructs.ValInt:
		if t == nil {
			return def
		}
		return t.Val
	case *xmlstructs.ValString:
		if t == nil {
			return def
		}
		return t.Val
	case *xmlstructs.Color:
		if t == nil {
			return def
		}
		return t.RGB
	}
	return def
}

func edgeKey(e *xmlstructs.BorderEdge) string {
	if e == nil {
		return "none"
	}
	return fmt.Sprintf("%s:%s", e.Style, valOr(e.Color, "none"))
}

func (e *styleProcessor) getNumFmtID(formatCode string) int {
	if id, ok := standardNumFmts[formatCode]; ok {
		return id
	}
	return e.styles.AddNumFmt(formatCode)
}

var standardNumFmts = map[string]int{
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

func (e *styleProcessor) setDataValidation(sheet string, ref string, options ...string) error {
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
