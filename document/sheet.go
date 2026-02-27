package document

// Sheet is a fluent handle bound to a specific worksheet.
// It provides sheet-scoped operations without repeatedly passing the sheet name or context.
type Sheet interface {
	Cell(axis string) Cell
	MergeCells(hRange string) Sheet
	SetColumnWidth(col int, width float64) Sheet
	AutoFilter(ref string) Sheet
	FreezePanes(col, row int) Sheet
	InsertImage(path string, x, y float64) Sheet
	GetCellValue(axis string) (string, error)
	Err() error
}

// Cell is a fluent handle bound to a specific cell in a sheet.
type Cell interface {
	Set(value any) Cell
	Style(style CellStyle) Cell
	Get() (string, error)
	Err() error
}
