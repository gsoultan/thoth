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
	s.err = s.processor().mergeCells(s.name, hRange)
	return s
}

func (s *sheetHandle) SetColumnWidth(col int, width float64) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().setColumnWidth(s.name, col, width)
	return s
}

func (s *sheetHandle) SetRowHeight(row int, height float64) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().setRowHeight(s.name, row, height)
	return s
}

func (s *sheetHandle) AutoFilter(ref string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().autoFilter(s.name, ref)
	return s
}

func (s *sheetHandle) FreezePanes(col, row int) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().freezePanes(s.name, col, row)
	return s
}

func (s *sheetHandle) InsertImage(path string, x, y float64) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().insertImage(s.name, path, x, y)
	return s
}

func (s *sheetHandle) SetDataValidation(ref string, options ...string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().sheetProcessor.setDataValidation(s.name, ref, options...)
	return s
}

func (s *sheetHandle) SetConditionalFormatting(ref string, style document.CellStyle) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().setConditionalFormatting(s.name, ref, "cellIs", "greaterThan", "0", style)
	return s
}

func (s *sheetHandle) SetPageSettings(settings document.PageSettings) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().setPageSettings(s.name, settings)
	return s
}

func (s *sheetHandle) Protect(password string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().protect(s.name, password)
	return s
}

func (s *sheetHandle) GroupRows(start, end int, level int) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().groupRows(s.name, start, end, level)
	return s
}

func (s *sheetHandle) GroupCols(start, end int, level int) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().groupCols(s.name, start, end, level)
	return s
}

func (s *sheetHandle) SetHeader(text string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().setHeader(s.name, text)
	return s
}

func (s *sheetHandle) SetFooter(text string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().setFooter(s.name, text)
	return s
}

func (s *sheetHandle) AddTable(ref string, name string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().addTable(s.name, ref, name)
	return s
}

func (s *sheetHandle) SetPrintArea(ref string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().setPrintArea(s.name, ref)
	return s
}

func (s *sheetHandle) SetPrintTitles(rowRef, colRef string) document.Sheet {
	if s.err != nil {
		return s
	}
	s.err = s.processor().setPrintTitles(s.name, rowRef, colRef)
	return s
}

func (s *sheetHandle) GetCellValue(axis string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.processor().getCellValue(s.name, axis)
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
	c.err = c.sheet.processor().setCellValue(c.sheet.name, c.axis, value)
	return c
}

func (c *cellHandle) Formula(formula string) document.Cell {
	if c.err != nil {
		return c
	}
	if c.sheet.err != nil {
		c.err = c.sheet.err
		return c
	}
	c.err = c.sheet.processor().setCellFormula(c.sheet.name, c.axis, formula)
	return c
}

func (c *cellHandle) Hyperlink(url string) document.Cell {
	if c.err != nil {
		return c
	}
	if c.sheet.err != nil {
		c.err = c.sheet.err
		return c
	}
	c.err = c.sheet.processor().setCellHyperlink(c.sheet.name, c.axis, url)
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
	c.err = c.sheet.processor().setCellStyle(c.sheet.name, c.axis, style)
	return c
}

func (c *cellHandle) Comment(text string) document.Cell {
	if c.err != nil {
		return c
	}
	if c.sheet.err != nil {
		c.err = c.sheet.err
		return c
	}
	// TODO: Implement actual comments. For now, we'll store it as Alt text in the cell for metadata purposes,
	// or just a placeholder as VML comments are complex to implement correctly without more infra.
	return c
}

func (c *cellHandle) Get() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if c.sheet.err != nil {
		return "", c.sheet.err
	}
	return c.sheet.processor().getCellValue(c.sheet.name, c.axis)
}

func (c *cellHandle) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.sheet.err
}
