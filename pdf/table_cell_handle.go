package pdf

import (
	"github.com/gsoultan/thoth/document"
)

// tableCellHandle is a fluent, cell-scoped helper implementing document.TableCell.
type tableCellHandle struct {
	row   *rowHandle
	index int
	err   error
}

func (c *tableCellHandle) AddParagraph(text string, style ...document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.err != nil {
		c.err = c.row.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellParagraph(c.row.table.tbl, c.row.index, c.index, text, style...)
	return c
}

func (c *tableCellHandle) AddRichParagraph(spans []document.TextSpan) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.err != nil {
		c.err = c.row.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellRichParagraph(c.row.table.tbl, c.row.index, c.index, spans)
	return c
}

func (c *tableCellHandle) AddImage(path string, width, height float64, style ...document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.err != nil {
		c.err = c.row.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellImage(c.row.table.tbl, c.row.index, c.index, path, width, height, style...)
	return c
}

func (c *tableCellHandle) AddList(items []string, ordered bool, style ...document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.err != nil {
		c.err = c.row.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellList(c.row.table.tbl, c.row.index, c.index, items, ordered, style...)
	return c
}

func (c *tableCellHandle) AddTable(rows, cols int) document.Table {
	if c.err != nil {
		return &tableHandle{err: c.err}
	}
	if c.row.err != nil {
		return &tableHandle{err: c.row.err}
	}
	table, err := (&processor{c.row.table.state}).addTableCellTable(c.row.table.tbl, c.row.index, c.index, rows, cols)
	if err != nil {
		return &tableHandle{err: err}
	}
	return table
}

func (c *tableCellHandle) Style(style document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.err != nil {
		c.err = c.row.err
		return c
	}
	c.err = (&processor{c.row.table.state}).setTableCellStyle(c.row.table.tbl, c.row.index, c.index, style)
	return c
}

func (c *tableCellHandle) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.row.table.err
}
