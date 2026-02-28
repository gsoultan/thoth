package word

import (
	"fmt"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func (p *processor) AddTable(rows, cols int) (document.Table, error) {
	if rows <= 0 || cols <= 0 {
		return nil, fmt.Errorf("table must have at least one row and one column")
	}
	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}

	tbl := &xmlstructs.Table{
		TblPr: &xmlstructs.TableProperties{
			TblStyle: &xmlstructs.TblStyle{Val: "TableGrid"},
			TblW:     &xmlstructs.TableWidth{W: 0, Type: "auto"},
		},
		TblGrid: &xmlstructs.TableGrid{},
	}

	for range rows {
		row := &xmlstructs.TableRow{}
		for range cols {
			row.Cells = append(row.Cells, &xmlstructs.TableCell{
				Content: []any{&xmlstructs.Paragraph{}},
			})
		}
		tbl.Rows = append(tbl.Rows, row)
	}

	for range cols {
		tbl.TblGrid.Cols = append(tbl.TblGrid.Cols, xmlstructs.TableGridCol{})
	}

	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, tbl)
	return &tableHandle{state: p.state, tbl: tbl}, nil
}

func (p *processor) getTable(index int) (*xmlstructs.Table, error) {
	var tableCount int
	for _, c := range p.xmlDoc.Body.Content {
		if t, ok := c.(*xmlstructs.Table); ok {
			if tableCount == index {
				return t, nil
			}
			tableCount++
		}
	}
	return nil, fmt.Errorf("table index %d not found", index)
}

func (p *processor) addTableCellParagraph(cell *xmlstructs.TableCell, text string, style ...document.CellStyle) error {
	var pPr *xmlstructs.ParagraphProperties
	var rPr *xmlstructs.RunProperties
	if len(style) > 0 {
		pPr = p.mapParagraphProperties(style[0])
		rPr = p.mapRunProperties(style[0])
	}

	par := &xmlstructs.Paragraph{PPr: pPr}
	if text != "" {
		par.Content = append(par.Content, &xmlstructs.Run{
			RPr: rPr,
			T:   text,
		})
	}

	p.addContentToCell(cell, par)
	return nil
}

func (p *processor) addTableCellImage(cell *xmlstructs.TableCell, path string, width, height float64, style ...document.CellStyle) error {
	par, err := p.createImageParagraph(path, width, height)
	if err != nil {
		return err
	}
	if len(style) > 0 {
		if style[0].Horizontal != "" {
			val := style[0].Horizontal
			if val == "justify" {
				val = "both"
			}
			par.PPr.Jc = &xmlstructs.Justification{Val: val}
		}
	}

	p.addContentToCell(cell, &par)
	return nil
}

func (p *processor) addTableCellList(cell *xmlstructs.TableCell, items []string, ordered bool, style ...document.CellStyle) error {
	p.ensureNumbering()
	numID := 1
	if ordered {
		numID = 2
	}

	var s document.CellStyle
	if len(style) > 0 {
		s = style[0]
	}

	// If the cell has only one empty paragraph, remove it
	if len(cell.Content) == 1 {
		if par, ok := cell.Content[0].(*xmlstructs.Paragraph); ok {
			if par.PPr == nil && len(par.Content) == 0 {
				cell.Content = nil
			}
		}
	}

	for _, item := range items {
		pPr := &xmlstructs.ParagraphProperties{
			NumPr: &xmlstructs.NumPr{
				ILvl:  &xmlstructs.ValInt{Val: 0},
				NumID: &xmlstructs.ValInt{Val: numID},
			},
			Ind: &xmlstructs.Ind{Left: 720, Hanging: 360},
		}
		if s.Name != "" {
			pPr.PStyle = &xmlstructs.ParagraphStyle{Val: s.Name}
		}

		rPr := p.mapRunProperties(s)

		par := &xmlstructs.Paragraph{
			PPr: pPr,
			Content: []any{
				&xmlstructs.Run{
					RPr: rPr,
					T:   item,
				},
			},
		}
		cell.Content = append(cell.Content, par)
	}
	return nil
}

func (p *processor) addTableCellRichParagraph(cell *xmlstructs.TableCell, spans []document.TextSpan) error {
	var pPr *xmlstructs.ParagraphProperties
	if len(spans) > 0 {
		pPr = p.mapParagraphProperties(spans[0].Style)
	}

	par := &xmlstructs.Paragraph{PPr: pPr}
	for _, span := range spans {
		if span.Text == "" {
			continue
		}
		rPr := p.mapRunProperties(span.Style)
		par.Content = append(par.Content, &xmlstructs.Run{RPr: rPr, T: span.Text})
	}

	p.addContentToCell(cell, par)
	return nil
}

func (p *processor) addTableCellTable(cell *xmlstructs.TableCell, rows, cols int) (document.Table, error) {
	if rows <= 0 || cols <= 0 {
		return nil, fmt.Errorf("table must have at least one row and one column")
	}
	innerTbl := &xmlstructs.Table{
		TblPr: &xmlstructs.TableProperties{
			TblStyle: &xmlstructs.TblStyle{Val: "TableGrid"},
			TblW:     &xmlstructs.TableWidth{W: 0, Type: "auto"},
		},
		TblGrid: &xmlstructs.TableGrid{},
	}

	for range rows {
		tr := &xmlstructs.TableRow{}
		for range cols {
			tr.Cells = append(tr.Cells, &xmlstructs.TableCell{
				Content: []any{&xmlstructs.Paragraph{}},
			})
		}
		innerTbl.Rows = append(innerTbl.Rows, tr)
	}

	for range cols {
		innerTbl.TblGrid.Cols = append(innerTbl.TblGrid.Cols, xmlstructs.TableGridCol{})
	}

	p.addContentToCell(cell, innerTbl)
	return &tableHandle{state: p.state, tbl: innerTbl}, nil
}

func (p *processor) addContentToCell(cell *xmlstructs.TableCell, content any) {
	if len(cell.Content) == 1 {
		if par, ok := cell.Content[0].(*xmlstructs.Paragraph); ok {
			if par.PPr == nil && len(par.Content) == 0 {
				cell.Content[0] = content
				return
			}
		}
	}
	cell.Content = append(cell.Content, content)
}

func (p *processor) Table(index int) (document.Table, error) {
	tbl, err := p.getTable(index)
	if err != nil {
		return nil, err
	}
	return &tableHandle{state: p.state, tbl: tbl}, nil
}

func (p *processor) setTableStyle(tbl *xmlstructs.Table, style string) error {
	if tbl.TblPr == nil {
		tbl.TblPr = &xmlstructs.TableProperties{}
	}
	tbl.TblPr.TblStyle = &xmlstructs.TblStyle{Val: style}

	if style == "TableGrid" {
		tbl.TblPr.TblBorders = &xmlstructs.TableBorders{
			Top:     &xmlstructs.BorderLine{Val: "single", Sz: 4, Color: "auto"},
			Left:    &xmlstructs.BorderLine{Val: "single", Sz: 4, Color: "auto"},
			Bottom:  &xmlstructs.BorderLine{Val: "single", Sz: 4, Color: "auto"},
			Right:   &xmlstructs.BorderLine{Val: "single", Sz: 4, Color: "auto"},
			InsideH: &xmlstructs.BorderLine{Val: "single", Sz: 4, Color: "auto"},
			InsideV: &xmlstructs.BorderLine{Val: "single", Sz: 4, Color: "auto"},
		}
	}
	return nil
}

func (p *processor) setTableHeaderRows(tbl *xmlstructs.Table, count int) error {
	for i := range min(count, len(tbl.Rows)) {
		if tbl.Rows[i].TrPr == nil {
			tbl.Rows[i].TrPr = &xmlstructs.TableRowProperties{}
		}
		tbl.Rows[i].TrPr.TblHeader = &struct{}{}
	}
	return nil
}

func (p *processor) setTableColumnWidths(tbl *xmlstructs.Table, widths ...float64) error {
	if tbl.TblGrid == nil {
		tbl.TblGrid = &xmlstructs.TableGrid{}
	}
	tbl.TblGrid.Cols = nil
	totalWidth := 0.0
	for _, width := range widths {
		twips := int(width * 20)
		tbl.TblGrid.Cols = append(tbl.TblGrid.Cols, xmlstructs.TableGridCol{W: twips})
		totalWidth += width
	}

	if tbl.TblPr == nil {
		tbl.TblPr = &xmlstructs.TableProperties{}
	}
	tbl.TblPr.TblW = &xmlstructs.TableWidth{W: int(totalWidth * 20), Type: "dxa"}

	for _, row := range tbl.Rows {
		for i, cell := range row.Cells {
			if i < len(widths) {
				if cell.TcPr == nil {
					cell.TcPr = &xmlstructs.TableCellProperties{}
				}
				cell.TcPr.TcW = &xmlstructs.TableCellWidth{W: int(widths[i] * 20), Type: "dxa"}
			}
		}
	}
	return nil
}

func (p *processor) mergeTableCells(tbl *xmlstructs.Table, row, col, rowSpan, colSpan int) error {
	if row < 0 || row >= len(tbl.Rows) || col < 0 || col >= len(tbl.Rows[row].Cells) {
		return fmt.Errorf("origin cell out of range")
	}

	if colSpan > 1 {
		cell := tbl.Rows[row].Cells[col]
		if cell.TcPr == nil {
			cell.TcPr = &xmlstructs.TableCellProperties{}
		}
		cell.TcPr.GridSpan = &xmlstructs.GridSpan{Val: colSpan}
	}

	if rowSpan > 1 {
		for r := range rowSpan {
			currRow := row + r
			if currRow >= len(tbl.Rows) {
				break
			}
			cell := tbl.Rows[currRow].Cells[col]
			if cell.TcPr == nil {
				cell.TcPr = &xmlstructs.TableCellProperties{}
			}
			if r == 0 {
				cell.TcPr.VMerge = &xmlstructs.VMerge{Val: "restart"}
			} else {
				cell.TcPr.VMerge = &xmlstructs.VMerge{}
			}
		}
	}
	return nil
}

func (p *processor) setTableCellStyle(cell *xmlstructs.TableCell, style document.CellStyle) error {
	cell.TcPr = p.mapTableCellProperties(style)
	return nil
}

func (p *processor) createImageParagraph(path string, width, height float64) (xmlstructs.Paragraph, error) {
	return p.createImageParagraphInternal(path, width, height)
}
