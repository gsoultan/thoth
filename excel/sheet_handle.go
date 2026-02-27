package excel

import (
	"context"

	"github.com/gsoultan/thoth/document"
)

// sheetHandle is a fluent, sheet-scoped helper implementing document.Sheet.
type sheetHandle struct {
	*state
	ctx  context.Context
	name string
	err  error
}

func (s *sheetHandle) Cell(axis string) document.Cell {
	return &cellHandle{sheet: s, axis: axis}
}

func (s *sheetHandle) MergeCells(hRange string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = (&sheetManager{s.state}).mergeCells(s.name, hRange)
	return s
}

func (s *sheetHandle) SetColumnWidth(col int, width float64) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = (&sheetManager{s.state}).setColumnWidth(s.name, col, width)
	return s
}

func (s *sheetHandle) AutoFilter(ref string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = (&sheetManager{s.state}).autoFilter(s.name, ref)
	return s
}

func (s *sheetHandle) FreezePanes(col, row int) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = (&sheetManager{s.state}).freezePanes(s.name, col, row)
	return s
}

func (s *sheetHandle) InsertImage(path string, x, y float64) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = (&imageProcessor{s.state}).insertImage(s.name, path, x, y)
	return s
}

func (s *sheetHandle) GetCellValue(axis string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return (&cellProcessor{s.state}).getCellValue(s.name, axis)
}

func (s *sheetHandle) Err() error {
	return s.err
}

// cellHandle is a fluent, cell-scoped helper implementing document.Cell.
type cellHandle struct {
	sheet *sheetHandle
	axis  string
	err   error
}

func (c *cellHandle) Set(value any) document.Cell {
	if c.err != nil {
		return c
	}
	if c.sheet.err != nil {
		c.err = c.sheet.err
		return c
	}
	c.err = (&cellProcessor{c.sheet.state}).setCellValue(c.sheet.name, c.axis, value)
	return c
}

func (c *cellHandle) Style(style document.CellStyle) document.Cell {
	if c.err != nil {
		return c
	}
	if c.sheet.err != nil {
		c.err = c.sheet.err
		return c
	}
	c.err = (&styleManager{c.sheet.state}).setCellStyle(c.sheet.name, c.axis, style)
	return c
}

func (c *cellHandle) Get() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if c.sheet.err != nil {
		return "", c.sheet.err
	}
	return (&cellProcessor{c.sheet.state}).getCellValue(c.sheet.name, c.axis)
}

func (c *cellHandle) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.sheet.err
}
