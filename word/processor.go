package word

import (
	"fmt"
	"os"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// processor handles adding content to the Word document.
type processor struct{ *state }

func (w *processor) SetFooter(text string, style ...document.CellStyle) error {
	return nil
}

func (w *processor) SetHeader(text string, style ...document.CellStyle) error {
	return nil
}

func (w *processor) DrawLine(x1, y1, x2, y2 float64, style ...document.CellStyle) error {
	return nil
}

func (w *processor) DrawRect(x, y, width, height float64, style ...document.CellStyle) error {
	return nil
}

func (w *processor) AddParagraph(text string, style ...document.CellStyle) error {
	if w.doc == nil {
		w.doc = &xmlstructs.Document{}
	}
	if w.xmlDoc == nil {
		w.xmlDoc = w.doc
	}

	var rPr *xmlstructs.RunProperties
	var pPr *xmlstructs.ParagraphProperties
	if len(style) > 0 {
		pPr = &xmlstructs.ParagraphProperties{}
		if style[0].Name != "" {
			pPr.PStyle = &xmlstructs.ParagraphStyle{Val: style[0].Name}
		}
		if style[0].Horizontal != "" {
			val := style[0].Horizontal
			if val == "center" {
				val = "center"
			} else if val == "right" {
				val = "right"
			} else if val == "justify" {
				val = "both"
			} else {
				val = "left"
			}
			pPr.Jc = &xmlstructs.Justification{Val: val}
		}
		rPr = &xmlstructs.RunProperties{}
		if style[0].Bold {
			rPr.Bold = &struct{}{}
		}
		if style[0].Italic {
			rPr.Italic = &struct{}{}
		}
	}

	p := xmlstructs.Paragraph{
		PPr: pPr,
		Runs: []xmlstructs.Run{
			{
				RPr: rPr,
				T:   text,
			},
		},
	}
	w.xmlDoc.Body.Content = append(w.xmlDoc.Body.Content, p)
	return nil
}

func (w *processor) InsertImage(path string, width, height float64, style ...document.CellStyle) error {
	p, err := w.createImageParagraph(path, width, height)
	if err != nil {
		return err
	}
	if len(style) > 0 {
		if style[0].Horizontal != "" {
			val := style[0].Horizontal
			if val == "justify" {
				val = "both"
			}
			p.PPr.Jc = &xmlstructs.Justification{Val: val}
		}
	}

	if w.doc == nil {
		w.doc = &xmlstructs.Document{}
	}
	if w.xmlDoc == nil {
		w.xmlDoc = w.doc
	}
	w.xmlDoc.Body.Content = append(w.xmlDoc.Body.Content, p)

	return nil
}

func (w *processor) createImageParagraph(path string, width, height float64) (xmlstructs.Paragraph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return xmlstructs.Paragraph{}, fmt.Errorf("read image file: %w", err)
	}

	// 1. Add to media
	imgName := fmt.Sprintf("image%d.png", len(w.media)+1)
	mediaPath := "word/media/" + imgName

	// 2. Add relationship
	if w.docRels == nil {
		w.docRels = &xmlstructs.Relationships{}
	}
	rID := w.docRels.AddRelationship(
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/image",
		"media/"+imgName,
	)

	w.media[mediaPath] = data

	// 3. Create drawing
	emuW := int64(width * 12700)
	emuH := int64(height * 12700)

	drawing := &xmlstructs.Drawing{
		Inline: &xmlstructs.Inline{
			Extent: xmlstructs.Extent{CX: emuW, CY: emuH},
			DocPr:  xmlstructs.DocPr{ID: len(w.media), Name: imgName},
			Graphic: xmlstructs.Graphic{
				Data: xmlstructs.GraphicData{
					URI: "http://schemas.openxmlformats.org/drawingml/2006/picture",
					Pic: xmlstructs.Pic{
						BlipFill: xmlstructs.BlipFill{
							Blip: xmlstructs.Blip{Embed: rID},
						},
					},
				},
			},
		},
	}

	return xmlstructs.Paragraph{
		Runs: []xmlstructs.Run{
			{Drawing: drawing},
		},
	}, nil
}

func (w *processor) AddTable(rows, cols int) (document.Table, error) {
	if w.xmlDoc == nil {
		w.xmlDoc = &xmlstructs.Document{}
	}
	if w.xmlDoc == nil {
		w.xmlDoc = w.doc
	}

	tbl := xmlstructs.Table{
		TblPr: &xmlstructs.TableProperties{
			TblStyle: &xmlstructs.TblStyle{Val: "TableGrid"},
		},
	}

	for range rows {
		row := xmlstructs.TableRow{}
		for range cols {
			row.Cells = append(row.Cells, xmlstructs.TableCell{
				Paragraphs: []xmlstructs.Paragraph{{}},
			})
		}
		tbl.Rows = append(tbl.Rows, row)
	}

	w.xmlDoc.Body.Content = append(w.xmlDoc.Body.Content, tbl)
	var tables []xmlstructs.Table
	for _, c := range w.xmlDoc.Body.Content {
		if t, ok := c.(xmlstructs.Table); ok {
			tables = append(tables, t)
		}
	}
	index := len(tables) - 1
	return &tableHandle{state: w.state, index: index}, nil
}

func (w *processor) getTable(index int) (*xmlstructs.Table, error) {
	var tableCount int
	for i, c := range w.xmlDoc.Body.Content {
		if t, ok := c.(xmlstructs.Table); ok {
			if tableCount == index {
				return &t, nil
			}
			tableCount++
		}
		_ = i
	}
	return nil, fmt.Errorf("table index %d not found", index)
}

func (w *processor) updateTable(index int, tbl *xmlstructs.Table) {
	var tableCount int
	for i, c := range w.xmlDoc.Body.Content {
		if _, ok := c.(xmlstructs.Table); ok {
			if tableCount == index {
				w.xmlDoc.Body.Content[i] = *tbl
				return
			}
			tableCount++
		}
	}
}

func (w *processor) addTableCellParagraph(tableIdx, row, col int, text string, style ...document.CellStyle) error {
	tbl, err := w.getTable(tableIdx)
	if err != nil {
		return err
	}
	if row < 0 || row >= len(tbl.Rows) {
		return fmt.Errorf("row index out of range")
	}
	if col < 0 || col >= len(tbl.Rows[row].Cells) {
		return fmt.Errorf("col index out of range")
	}

	var rPr *xmlstructs.RunProperties
	var pPr *xmlstructs.ParagraphProperties
	if len(style) > 0 {
		pPr = &xmlstructs.ParagraphProperties{}
		if style[0].Name != "" {
			pPr.PStyle = &xmlstructs.ParagraphStyle{Val: style[0].Name}
		}
		if style[0].Horizontal != "" {
			val := style[0].Horizontal
			if val == "center" {
				val = "center"
			} else if val == "right" {
				val = "right"
			} else if val == "justify" {
				val = "both"
			} else {
				val = "left"
			}
			pPr.Jc = &xmlstructs.Justification{Val: val}
		}
		rPr = &xmlstructs.RunProperties{}
		if style[0].Bold {
			rPr.Bold = &struct{}{}
		}
		if style[0].Italic {
			rPr.Italic = &struct{}{}
		}
	}

	p := xmlstructs.Paragraph{
		PPr:  pPr,
		Runs: []xmlstructs.Run{{RPr: rPr, T: text}},
	}
	cell := &tbl.Rows[row].Cells[col]
	if len(cell.Paragraphs) == 1 && cell.Paragraphs[0].PPr == nil && len(cell.Paragraphs[0].Runs) == 0 {
		cell.Paragraphs[0] = p
	} else {
		cell.Paragraphs = append(cell.Paragraphs, p)
	}
	w.updateTable(tableIdx, tbl)
	return nil
}

func (w *processor) addTableCellImage(tableIdx, row, col int, path string, width, height float64, style ...document.CellStyle) error {
	tbl, err := w.getTable(tableIdx)
	if err != nil {
		return err
	}
	if row < 0 || row >= len(tbl.Rows) {
		return fmt.Errorf("row index out of range")
	}
	if col < 0 || col >= len(tbl.Rows[row].Cells) {
		return fmt.Errorf("col index out of range")
	}

	p, err := w.createImageParagraph(path, width, height)
	if err != nil {
		return err
	}
	if len(style) > 0 {
		if style[0].Horizontal != "" {
			val := style[0].Horizontal
			if val == "justify" {
				val = "both"
			}
			p.PPr.Jc = &xmlstructs.Justification{Val: val}
		}
	}

	cell := &tbl.Rows[row].Cells[col]
	if len(cell.Paragraphs) == 1 && cell.Paragraphs[0].PPr == nil && len(cell.Paragraphs[0].Runs) == 0 {
		cell.Paragraphs[0] = p
	} else {
		cell.Paragraphs = append(cell.Paragraphs, p)
	}
	w.updateTable(tableIdx, tbl)
	return nil
}

// Table returns a fluent, table-scoped handle.
func (w *processor) AddSection(settings document.PageSettings) error {
	if w.xmlDoc == nil {
		w.xmlDoc = &xmlstructs.Document{}
	}

	// 1. Move current body sectPr to a new paragraph if body is not empty
	var sectPr *xmlstructs.SectPr
	if w.xmlDoc.Body.SectPr != nil {
		sectPr = w.xmlDoc.Body.SectPr
	} else {
		// Defaultportrait A4
		sectPr = &xmlstructs.SectPr{
			PgSz: &xmlstructs.PgSz{W: 11906, H: 16838, Orient: "portrait"},
		}
	}

	p := xmlstructs.Paragraph{
		PPr: &xmlstructs.ParagraphProperties{
			SectPr: sectPr,
		},
	}
	w.xmlDoc.Body.Content = append(w.xmlDoc.Body.Content, p)

	// 2. Set new body sectPr for next content
	newSect := &xmlstructs.SectPr{}
	w.applyPageSettingsToSect(newSect, settings)
	w.xmlDoc.Body.SectPr = newSect

	return nil
}

func (w *processor) applyPageSettingsToSect(sect *xmlstructs.SectPr, settings document.PageSettings) {
	if sect.PgSz == nil {
		sect.PgSz = &xmlstructs.PgSz{}
	}

	if settings.Orientation == document.OrientationLandscape {
		sect.PgSz.Orient = "landscape"
		sect.PgSz.W = 16838
		sect.PgSz.H = 11906
	} else {
		sect.PgSz.Orient = "portrait"
		sect.PgSz.W = 11906
		sect.PgSz.H = 16838
	}

	switch settings.PaperType {
	case document.PaperLetter:
		if settings.Orientation == document.OrientationLandscape {
			sect.PgSz.W = 15840
			sect.PgSz.H = 12240
		} else {
			sect.PgSz.W = 12240
			sect.PgSz.H = 15840
		}
	}

	sect.PgMar = &xmlstructs.PgMar{
		Top:    int(settings.Margins.Top * 1440),
		Bottom: int(settings.Margins.Bottom * 1440),
		Left:   int(settings.Margins.Left * 1440),
		Right:  int(settings.Margins.Right * 1440),
	}
}

func (w *processor) Table(index int) (document.Table, error) {
	var tables []xmlstructs.Table
	for _, c := range w.xmlDoc.Body.Content {
		if t, ok := c.(xmlstructs.Table); ok {
			tables = append(tables, t)
		}
	}

	if w.xmlDoc == nil || index < 0 || index >= len(tables) {
		return nil, fmt.Errorf("table index %d not found", index)
	}
	return &tableHandle{state: w.state, index: index}, nil
}

func (w *processor) setTableStyle(tableIdx int, style string) error {
	tbl, err := w.getTable(tableIdx)
	if err != nil {
		return err
	}
	if tbl.TblPr == nil {
		tbl.TblPr = &xmlstructs.TableProperties{}
	}
	tbl.TblPr.TblStyle = &xmlstructs.TblStyle{Val: style}
	w.updateTable(tableIdx, tbl)
	return nil
}

func (w *processor) mergeTableCells(tableIdx, row, col, rowSpan, colSpan int) error {
	tbl, err := w.getTable(tableIdx)
	if err != nil {
		return err
	}

	if row < 0 || row >= len(tbl.Rows) || col < 0 || col >= len(tbl.Rows[row].Cells) {
		return fmt.Errorf("origin cell out of range")
	}

	// Horizontal merge
	if colSpan > 1 {
		cell := &tbl.Rows[row].Cells[col]
		if cell.TcPr == nil {
			cell.TcPr = &xmlstructs.TableCellProperties{}
		}
		cell.TcPr.GridSpan = &xmlstructs.GridSpan{Val: colSpan}
	}

	// Vertical merge
	if rowSpan > 1 {
		for r := range rowSpan {
			currRow := row + r
			if currRow >= len(tbl.Rows) {
				break
			}
			cell := &tbl.Rows[currRow].Cells[col]
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

	w.updateTable(tableIdx, tbl)
	return nil
}

func (w *processor) setTableCellStyle(tableIdx, row, col int, style document.CellStyle) error {
	tbl, err := w.getTable(tableIdx)
	if err != nil {
		return err
	}
	if row < 0 || row >= len(tbl.Rows) || col < 0 || col >= len(tbl.Rows[row].Cells) {
		return fmt.Errorf("cell index out of range")
	}

	cell := &tbl.Rows[row].Cells[col]
	if cell.TcPr == nil {
		cell.TcPr = &xmlstructs.TableCellProperties{}
	}
	if style.Background != "" {
		cell.TcPr.Shd = &xmlstructs.TableCellShading{
			Val:  "clear",
			Fill: style.Background,
		}
	}
	if style.Border {
		cell.TcPr.TcBorders = &xmlstructs.TableCellBorders{
			Top:    &xmlstructs.BorderLine{Val: "single", Sz: 4, Space: 0, Color: "auto"},
			Left:   &xmlstructs.BorderLine{Val: "single", Sz: 4, Space: 0, Color: "auto"},
			Bottom: &xmlstructs.BorderLine{Val: "single", Sz: 4, Space: 0, Color: "auto"},
			Right:  &xmlstructs.BorderLine{Val: "single", Sz: 4, Space: 0, Color: "auto"},
		}
	}
	w.updateTable(tableIdx, tbl)
	return nil
}

func (w *processor) AddPageBreak() error {
	if w.doc == nil {
		w.doc = &xmlstructs.Document{}
	}
	if w.xmlDoc == nil {
		w.xmlDoc = w.doc
	}

	p := xmlstructs.Paragraph{
		Runs: []xmlstructs.Run{
			{
				Br: &xmlstructs.Break{Type: "page"},
			},
		},
	}
	w.xmlDoc.Body.Content = append(w.xmlDoc.Body.Content, p)
	return nil
}
