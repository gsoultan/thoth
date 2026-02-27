package pdf

import (
	"github.com/gsoultan/thoth/document"
)

// tableHandle is a fluent, table-scoped helper implementing document.Table.
type tableHandle struct {
	state *state
	index int
	err   error
}

func (t *tableHandle) Row(index int) document.Row {
	return &rowHandle{table: t, index: index}
}

func (t *tableHandle) MergeCells(row, col, rowSpan, colSpan int) document.Table {
	if t.err != nil {
		return t
	}
	t.err = (&processor{t.state}).mergeTableCells(t.index, row, col, rowSpan, colSpan)
	return t
}

func (t *tableHandle) SetStyle(style string) document.Table {
	if t.err != nil {
		return t
	}
	t.err = (&processor{t.state}).setTableStyle(t.index, style)
	return t
}

func (t *tableHandle) Err() error {
	return t.err
}

// rowHandle is a fluent, row-scoped helper implementing document.Row.
type rowHandle struct {
	table *tableHandle
	index int
}

func (r *rowHandle) Cell(index int) document.TableCell {
	return &tableCellHandle{row: r, index: index}
}

func (r *rowHandle) Err() error {
	return r.table.err
}

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
	if c.row.table.err != nil {
		c.err = c.row.table.err
		return c
	}
	c.err = (&processor{c.row.table.state}).addTableCellParagraph(c.row.table.index, c.row.index, c.index, text, style...)
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
	c.err = (&processor{c.row.table.state}).addTableCellImage(c.row.table.index, c.row.index, c.index, path, width, height, style...)
	return c
}

func (c *tableCellHandle) Style(style document.CellStyle) document.TableCell {
	if c.err != nil {
		return c
	}
	if c.row.table.err != nil {
		c.err = c.row.table.err
		return c
	}
	c.err = (&processor{c.row.table.state}).setTableCellStyle(c.row.table.index, c.row.index, c.index, style)
	return c
}

func (c *tableCellHandle) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.row.table.err
}
