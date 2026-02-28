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
	if text == "" {
		return fmt.Errorf("footer text cannot be empty")
	}
	item := &contentItem{isParagraph: true, text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.footer = append(p.footer, item)
	return nil
}

func (p *processor) SetHeader(text string, style ...document.CellStyle) error {
	if text == "" {
		return fmt.Errorf("header text cannot be empty")
	}
	item := &contentItem{isParagraph: true, text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.header = append(p.header, item)
	return nil
}

func (p *processor) DrawLine(x1, y1, x2, y2 float64, style ...document.CellStyle) error {
	item := &contentItem{isShape: true, shapeType: "line", x1: x1, y1: y1, x2: x2, y2: y2}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) DrawRect(x, y, w, h float64, style ...document.CellStyle) error {
	if w <= 0 || h <= 0 {
		return fmt.Errorf("rectangle width and height must be positive")
	}
	item := &contentItem{isShape: true, shapeType: "rect", x1: x, y1: y, width: w, height: h}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) DrawEllipse(x, y, w, h float64, style ...document.CellStyle) error {
	if w <= 0 || h <= 0 {
		return fmt.Errorf("ellipse width and height must be positive")
	}
	item := &contentItem{isShape: true, shapeType: "ellipse", x1: x, y1: y, width: w, height: h}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) RegisterFont(name, path string) error {
	if p.fonts == nil {
		p.fonts = make(map[string]string)
	}
	p.fonts[name] = path
	return nil
}

func (p *processor) AddTextField(name string, x, y, w, h float64) error {
	item := &contentItem{isFormField: true, fieldType: "text", fieldName: name, x1: x, y1: y, width: w, height: h}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddCheckbox(name string, x, y float64) error {
	item := &contentItem{isFormField: true, fieldType: "checkbox", fieldName: name, x1: x, y1: y, width: 12, height: 12}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddComboBox(name string, x, y, w, h float64, options ...string) error {
	item := &contentItem{isFormField: true, fieldType: "combobox", fieldName: name, x1: x, y1: y, width: w, height: h, options: options}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddRadioButton(name string, x, y float64, options ...string) error {
	item := &contentItem{isFormField: true, fieldType: "radio", fieldName: name, x1: x, y1: y, width: 12, height: 12, options: options}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) ImportPage(path string, pageNum int) error {
	item := &contentItem{isImported: true, importPath: path, importPage: pageNum}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddFootnote(text string) error {
	item := &contentItem{isFootnote: true, text: text}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddHyperlink(text, url string, style ...document.CellStyle) error {
	if text == "" || url == "" {
		return fmt.Errorf("hyperlink text and url cannot be empty")
	}
	var s document.CellStyle
	if len(style) > 0 {
		s = style[0]
	}
	s.Link = url
	return p.AddParagraph(text, s)
}

func (p *processor) AddBookmark(name string) error {
	item := &contentItem{isBookmark: true, bookmarkName: name}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddParagraph(text string, style ...document.CellStyle) error {
	item := &contentItem{isParagraph: true, text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddRichParagraph(spans []document.TextSpan) error {
	if len(spans) == 0 {
		return fmt.Errorf("rich paragraph spans cannot be empty")
	}
	item := &contentItem{isRich: true, spans: spans}
	if len(spans) > 0 {
		item.style = spans[0].Style
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddHeading(text string, level int, style ...document.CellStyle) error {
	if text == "" {
		return fmt.Errorf("heading text cannot be empty")
	}
	if level < 1 || level > 6 {
		return fmt.Errorf("heading level must be between 1 and 6")
	}
	item := &contentItem{isHeading: true, level: level, text: text}
	if len(style) > 0 {
		item.style = style[0]
	} else {
		size := 14
		bold := true
		if level == 1 {
			size = 20
		} else if level == 2 {
			size = 18
		} else if level == 3 {
			size = 16
		}
		item.style = document.CellStyle{Size: size, Bold: bold}
	}
	if item.style.Name == "" {
		if level == 1 {
			item.style.Name = "Title"
		} else {
			item.style.Name = fmt.Sprintf("Heading%d", level)
		}
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) InsertImage(path string, width, height float64, style ...document.CellStyle) error {
	if path == "" {
		return fmt.Errorf("image path cannot be empty")
	}
	if width <= 0 || height <= 0 {
		return fmt.Errorf("image width and height must be positive")
	}
	item := &contentItem{isImage: true, path: path, width: width, height: height}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddTable(rows, cols int) (document.Table, error) {
	if rows <= 0 || cols <= 0 {
		return nil, fmt.Errorf("table must have at least one row and one column")
	}
	cells := make([][][]cellItem, rows)
	for i := range rows {
		cells[i] = make([][]cellItem, cols)
	}
	item := &contentItem{isTable: true, rows: rows, cols: cols, cells: cells}
	p.contentItems = append(p.contentItems, item)
	return &tableHandle{state: p.state, tbl: item}, nil
}

func (p *processor) addTableCellParagraph(tbl *contentItem, row, col int, text string, style ...document.CellStyle) error {
	if tbl == nil || row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	item := cellItem{text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	tbl.cells[row][col] = append(tbl.cells[row][col], item)
	return nil
}

func (p *processor) addTableCellImage(tbl *contentItem, row, col int, path string, width, height float64, style ...document.CellStyle) error {
	if tbl == nil || row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	item := cellItem{isImage: true, path: path, width: width, height: height}
	if len(style) > 0 {
		item.style = style[0]
	}
	tbl.cells[row][col] = append(tbl.cells[row][col], item)
	return nil
}

func (p *processor) addTableCellList(tbl *contentItem, row, col int, items []string, ordered bool, style ...document.CellStyle) error {
	if tbl == nil || row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	item := cellItem{isList: true, listItems: items, ordered: ordered}
	if len(style) > 0 {
		item.style = style[0]
	}
	tbl.cells[row][col] = append(tbl.cells[row][col], item)
	return nil
}

func (p *processor) addTableCellRichParagraph(tbl *contentItem, row, col int, spans []document.TextSpan) error {
	if tbl == nil || row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	item := cellItem{isRich: true, spans: spans}
	if len(spans) > 0 {
		item.style = spans[0].Style
	}
	tbl.cells[row][col] = append(tbl.cells[row][col], item)
	return nil
}

func (p *processor) addTableCellTable(tbl *contentItem, row, col int, rows, cols int) (document.Table, error) {
	if tbl == nil || row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return nil, fmt.Errorf("cell index out of range")
	}
	cells := make([][][]cellItem, rows)
	for i := range rows {
		cells[i] = make([][]cellItem, cols)
	}
	innerTbl := &contentItem{isTable: true, rows: rows, cols: cols, cells: cells}
	tbl.cells[row][col] = append(tbl.cells[row][col], cellItem{isTable: true, table: innerTbl})
	return &tableHandle{state: p.state, tbl: innerTbl}, nil
}

func (p *processor) setTableCellStyle(tbl *contentItem, row, col int, style document.CellStyle) error {
	if tbl == nil || row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	if len(tbl.cells[row][col]) == 0 {
		tbl.cells[row][col] = append(tbl.cells[row][col], cellItem{style: style})
	} else {
		tbl.cells[row][col][0].style = style
	}
	return nil
}

func (p *processor) setTableStyle(tbl *contentItem, style string) error {
	if tbl == nil {
		return fmt.Errorf("table is nil")
	}
	tbl.text = style
	return nil
}

func (p *processor) setTableColumnWidths(tbl *contentItem, widths ...float64) error {
	if tbl == nil {
		return fmt.Errorf("table is nil")
	}
	tbl.colWidths = widths
	return nil
}

func (p *processor) setTableHeaderRows(tbl *contentItem, count int) error {
	if tbl == nil {
		return fmt.Errorf("table is nil")
	}
	tbl.HeaderRows = count
	return nil
}

func (p *processor) mergeTableCells(tbl *contentItem, row, col, rowSpan, colSpan int) error {
	if tbl == nil || row < 0 || row >= tbl.rows || col < 0 || col >= tbl.cols {
		return fmt.Errorf("cell index out of range")
	}
	if len(tbl.cells[row][col]) == 0 {
		tbl.cells[row][col] = append(tbl.cells[row][col], cellItem{})
	}
	tbl.cells[row][col][0].rowSpan = rowSpan
	tbl.cells[row][col][0].colSpan = colSpan
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

func (p *processor) SetWatermark(text string, style ...document.CellStyle) error {
	item := &contentItem{isParagraph: true, text: text}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.watermark = item
	return nil
}

func (p *processor) AddList(items []string, ordered bool, style ...document.CellStyle) error {
	item := &contentItem{isList: true, listItems: items, ordered: ordered}
	if len(style) > 0 {
		item.style = style[0]
	}
	p.contentItems = append(p.contentItems, item)
	return nil
}

func (p *processor) AddPageBreak() error {
	p.contentItems = append(p.contentItems, &contentItem{isPageBreak: true})
	return nil
}

func (p *processor) AddSection(settings document.PageSettings) error {
	p.contentItems = append(p.contentItems, &contentItem{isPageBreak: true})
	p.pageSettings = settings
	return nil
}

func (p *processor) Table(index int) (document.Table, error) {
	tbl, err := p.getTableByIdx(index)
	if err != nil {
		return nil, err
	}
	return &tableHandle{state: p.state, tbl: tbl}, nil
}

func (p *processor) getTableByIdx(idx int) (*contentItem, error) {
	tableCount := 0
	for i := range p.contentItems {
		if p.contentItems[i].isTable {
			if tableCount == idx {
				return p.contentItems[i], nil
			}
			tableCount++
		}
	}
	return nil, fmt.Errorf("table index %d not found", idx)
}

func (p *processor) AddTableOfContents() error {
	p.contentItems = append(p.contentItems, &contentItem{isTOC: true})
	return nil
}

func (p *processor) AttachFile(path, name, description string) error {
	p.attachments = append(p.attachments, attachment{
		path:        path,
		name:        name,
		description: description,
	})
	return nil
}
