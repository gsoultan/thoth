package word

import (
	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// tableCellHandle is a fluent, cell-scoped helper implementing document.TableCell.
type tableCellHandle struct {
	row  *rowHandle
	cell *xmlstructs.TableCell
	err  error
}

func (c *tableCellHandle) AddParagraph(text string, style ...document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.table.err != nil {
		c.err = c.row.table.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellParagraph(c.cell, text, style...)
	return c
}

func (c *tableCellHandle) AddRichParagraph(spans []document.TextSpan) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.table.err != nil {
		c.err = c.row.table.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellRichParagraph(c.cell, spans)
	return c
}

func (c *tableCellHandle) AddImage(path string, width, height float64, style ...document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.table.err != nil {
		c.err = c.row.table.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellImage(c.cell, path, width, height, style...)
	return c
}

func (c *tableCellHandle) AddList(items []string, ordered bool, style ...document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.table.err != nil {
		c.err = c.row.table.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellList(c.cell, items, ordered, style...)
	return c
}

func (c *tableCellHandle) AddTable(rows, cols int) document.Table {
	if c.err != nil {
		return &tableHandle{err: c.err}
	}
	if c.row.table.err != nil {
		return &tableHandle{err: c.row.table.err}
	}
	table, err := (&processor{c.row.table.state}).addTableCellTable(c.cell, rows, cols)
	if err != nil {
		return &tableHandle{err: err}
	}
	return table
}

func (c *tableCellHandle) Style(style document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.table.err != nil {
		c.err = c.row.table.err
		return c
	}
	c.err = (&processor{c.row.table.state}).setTableCellStyle(c.cell, style)
	return c
}

func (c *tableCellHandle) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.row.table.err
}
