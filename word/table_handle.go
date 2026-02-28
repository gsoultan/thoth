package word

import (
	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// tableHandle is a fluent, table-scoped helper implementing document.Table.
type tableHandle struct {
	state *state
	tbl   *xmlstructs.Table
	err   error
}

func (t *tableHandle) Row(index int) document.Row {
	if t.err != nil {
		return &rowHandle{table: t, err: t.err}
	}
	if t.tbl == nil || index < 0 || index >= len(t.tbl.Rows) {
		return &rowHandle{table: t, err: t.err}
	}
	return &rowHandle{table: t, row: t.tbl.Rows[index]}
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
