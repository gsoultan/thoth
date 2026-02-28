package pdf

import (
	"fmt"

	"github.com/gsoultan/thoth/document"
)

// tableHandle is a fluent, table-scoped helper implementing document.Table.
type tableHandle struct {
	state *state
	tbl   *contentItem
	err   error
}

func (t *tableHandle) Row(index int) document.Row {
	if t.err != nil {
		return &rowHandle{table: t, err: t.err}
	}
	if t.tbl == nil || index < 0 || index >= t.tbl.rows {
		return &rowHandle{table: t, err: fmt.Errorf("row index out of range")}
	}
	return &rowHandle{table: t, index: index}
}

func (t *tableHandle) MergeCells(row, col, rowSpan, colSpan int) document.Table {
	if t.err != nil {
		return t
	}
	t.err = (&processor{t.state}).mergeTableCells(t.tbl, row, col, rowSpan, colSpan)
	return t
}

func (t *tableHandle) SetColumnWidths(widths ...float64) document.Table {
	if t.err != nil {
		return t
	}
	t.err = (&processor{t.state}).setTableColumnWidths(t.tbl, widths...)
	return t
}

func (t *tableHandle) SetHeaderRows(count int) document.Table {
	if t.err != nil {
		return t
	}
	t.err = (&processor{t.state}).setTableHeaderRows(t.tbl, count)
	return t
}

func (t *tableHandle) SetStyle(style string) document.Table {
	if t.err != nil {
		return t
	}
	t.err = (&processor{t.state}).setTableStyle(t.tbl, style)
	return t
}

func (t *tableHandle) Err() error {
	return t.err
}
