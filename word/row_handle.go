package word

import (
	"fmt"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// rowHandle is a fluent, row-scoped helper implementing document.Row.
type rowHandle struct {
	table *tableHandle
	row   *xmlstructs.TableRow
	err   error
}

func (r *rowHandle) Cell(index int) document.TableCell {
	if r.err != nil {
		return &tableCellHandle{row: r, err: r.err}
	}
	if r.row == nil || index < 0 || index >= len(r.row.Cells) {
		return &tableCellHandle{row: r, err: fmt.Errorf("cell index %d out of range", index)}
	}
	return &tableCellHandle{row: r, cell: r.row.Cells[index]}
}

func (r *rowHandle) Err() error {
	if r.err != nil {
		return r.err
	}
	return r.table.err
}
