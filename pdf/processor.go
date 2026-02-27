package pdf

import (
	"fmt"

	"github.com/gsoultan/thoth/document"
)

// processor handles operations specific to word processing documents.
type processor struct{ *state }

// WordProcessor methods

func (p *processor) SetPageSettings(settings document.PageSettings) error {
	p.pageSettings = settings
	return nil
}

func (p *processor) SetFooter(text string, style ...document.CellStyle) error {
	item := contentItem{isParagraph: true, text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.footer = append(p.footer, item)
	return nil
}

func (p *processor) SetHeader(text string, style ...document.CellStyle) error {
	item := contentItem{isParagraph: true, text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.header = append(p.header, item)
	return nil
}

func (p *processor) DrawLine(x1, y1, x2, y2 float64, style ...document.CellStyle) error {
	item := contentItem{isShape: true, shapeType: "line", x1: x1, y1: y1, x2: x2, y2: y2}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) DrawRect(x, y, w, h float64, style ...document.CellStyle) error {
	item := contentItem{isShape: true, shapeType: "rect", x1: x, y1: y, width: w, height: h}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddParagraph(text string, style ...document.CellStyle) error {
	item := contentItem{isParagraph: true, text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) InsertImage(path string, width, height float64, style ...document.CellStyle) error {
	item := contentItem{isImage: true, path: path, width: width, height: height}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddTable(rows, cols int) (document.Table, error) {
	cells := make([][][]cellItem, rows)
	for i := range rows {
		cells[i] = make([][]cellItem, cols)
	}
	p.contentItems = append(p.contentItems, contentItem{isTable: true, rows: rows, cols: cols, cells: cells})
	index := 0
	for i := range len(p.contentItems) - 1 {
		if p.contentItems[i].isTable {
			index++
		}
	}
	return &tableHandle{state: p.state, index: index}, nil
}

func (p *processor) addTableCellParagraph(tableIdx, row, col int, text string, style ...document.CellStyle) error {
	tbl, err := p.getTable(tableIdx)
	if err != nil {
		return err
	}
	if row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	item := cellItem{text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	tbl.cells[row][col] = append(tbl.cells[row][col], item)
	return nil
}

func (p *processor) addTableCellImage(tableIdx, row, col int, path string, width, height float64, style ...document.CellStyle) error {
	tbl, err := p.getTable(tableIdx)
	if err != nil {
		return err
	}
	if row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	item := cellItem{isImage: true, path: path, width: width, height: height}
	if len(style) > 0 {
		item.style = style[0]
	}
	tbl.cells[row][col] = append(tbl.cells[row][col], item)
	return nil
}

func (p *processor) setTableCellStyle(tableIdx, row, col int, style document.CellStyle) error {
	tbl, err := p.getTable(tableIdx)
	if err != nil {
		return err
	}
	if row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	if len(tbl.cells[row][col]) == 0 {
		tbl.cells[row][col] = append(tbl.cells[row][col], cellItem{style: style})
	} else {
		tbl.cells[row][col][0].style = style
	}
	return nil
}

func (p *processor) setTableStyle(tableIdx int, style string) error {
	tbl, err := p.getTable(tableIdx)
	if err != nil {
		return err
	}
	tbl.text = style
	return nil
}

func (p *processor) mergeTableCells(tableIdx, row, col, rowSpan, colSpan int) error {
	tbl, err := p.getTable(tableIdx)
	if err != nil {
		return err
	}
	if row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}

	// First cell holds the span
	if len(tbl.cells[row][col]) == 0 {
		tbl.cells[row][col] = append(tbl.cells[row][col], cellItem{})
	}
	tbl.cells[row][col][0].rowSpan = rowSpan
	tbl.cells[row][col][0].colSpan = colSpan

	// Mark others as hidden
	for r := row; r < row+rowSpan && r < tbl.rows; r++ {
		for c := col; c < col+colSpan && c < tbl.cols; c++ {
			if r == row && c == col {
				continue
			}
			if len(tbl.cells[r][c]) == 0 {
				tbl.cells[r][c] = append(tbl.cells[r][c], cellItem{})
			}
			tbl.cells[r][c][0].hidden = true
		}
	}
	return nil
}

func (p *processor) AddPageBreak() error {
	p.contentItems = append(p.contentItems, contentItem{isPageBreak: true})
	return nil
}

func (p *processor) AddSection(settings document.PageSettings) error {
	// For PDF, a section break is just a page break with new settings
	p.contentItems = append(p.contentItems, contentItem{isPageBreak: true})
	p.pageSettings = settings
	return nil
}

func (p *processor) Table(index int) (document.Table, error) {
	if _, err := p.getTable(index); err != nil {
		return nil, err
	}
	return &tableHandle{state: p.state, index: index}, nil
}

func (p *processor) getTable(idx int) (*contentItem, error) {
	tableCount := 0
	for i := range p.contentItems {
		if p.contentItems[i].isTable {
			if tableCount == idx {
				return &p.contentItems[i], nil
			}
			tableCount++
		}
	}
	return nil, fmt.Errorf("table index %d not found", idx)
}
