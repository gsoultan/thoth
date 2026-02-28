package document

// Table is a fluent handle bound to a specific Word/PDF table instance.
// It provides table-scoped operations without passing the table index each time.
type Table interface {
	Row(index int) Row
	MergeCells(row, col, rowSpan, colSpan int) Table
	SetColumnWidths(widths ...float64) Table
	SetHeaderRows(count int) Table
	SetStyle(style string) Table
	Err() error
}

// Row is a fluent handle bound to a specific row in a table.
type Row interface {
	Cell(index int) TableCell
	Err() error
}

// TableCell is a fluent handle bound to a specific cell in a table row.
type TableCell interface {
	AddParagraph(text string, style ...CellStyle) TableCell
	AddRichParagraph(spans []TextSpan) TableCell
	AddImage(path string, width, height float64, style ...CellStyle) TableCell
	AddList(items []string, ordered bool, style ...CellStyle) TableCell
	AddTable(rows, cols int) Table
	Style(style CellStyle) TableCell
	Err() error
}
