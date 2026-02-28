package pdf

import (
	"fmt"

	"github.com/gsoultan/thoth/document"
)

// rowHandle is a fluent, row-scoped helper implementing document.Row.
type rowHandle struct {
	table *tableHandle
	index int
	err   error
}

func (r *rowHandle) Cell(index int) document.TableCell {
	if r.err != nil {
		return &tableCellHandle{err: r.err}
	}
	if r.table.err != nil {
		return &tableCellHandle{err: r.table.err}
	}
	if r.table.tbl == nil || index < 0 || index >= r.table.tbl.cols {
		return &tableCellHandle{err: fmt.Errorf("cell index out of range")}
	}
	return &tableCellHandle{row: r, index: index}
}

func (r *rowHandle) Err() error {
	if r.err != nil {
		return r.err
	}
	return r.table.err
}
