package document

// Sheet is a fluent handle bound to a specific worksheet.
// It provides sheet-scoped operations without repeatedly passing the sheet name or context.
type Sheet interface {
	Cell(axis string) Cell
	MergeCells(hRange string) Sheet
	SetColumnWidth(col int, width float64) Sheet
	SetRowHeight(row int, height float64) Sheet
	AutoFilter(ref string) Sheet
	FreezePanes(col, row int) Sheet
	InsertImage(path string, x, y float64) Sheet
	SetDataValidation(ref string, options ...string) Sheet
	SetConditionalFormatting(ref string, style CellStyle) Sheet
	SetPageSettings(settings PageSettings) Sheet
	Protect(password string) Sheet
	GroupRows(start, end int, level int) Sheet
	GroupCols(start, end int, level int) Sheet
	SetHeader(text string) Sheet
	SetFooter(text string) Sheet
	AddTable(ref string, name string) Sheet
	SetPrintArea(ref string) Sheet
	SetPrintTitles(rowRef, colRef string) Sheet
	GetCellValue(axis string) (string, error)
	Err() error
}

// Cell is a fluent handle bound to a specific cell in a sheet.
type Cell interface {
	Set(value any) Cell
	Formula(formula string) Cell
	Hyperlink(url string) Cell
	Style(style CellStyle) Cell
	Comment(text string) Cell
	Get() (string, error)
	Err() error
}
