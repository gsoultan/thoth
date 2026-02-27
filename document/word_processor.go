package document

// WordProcessor defines operations specific to word processing documents (Word).
type WordProcessor interface {
	Document
	AddParagraph(text string, style ...CellStyle) error
	InsertImage(path string, width, height float64, style ...CellStyle) error
	AddTable(rows, cols int) (Table, error)

	AddPageBreak() error
	AddSection(settings PageSettings) error

	SetHeader(text string, style ...CellStyle) error
	SetFooter(text string, style ...CellStyle) error

	DrawLine(x1, y1, x2, y2 float64, style ...CellStyle) error
	DrawRect(x, y, w, h float64, style ...CellStyle) error

	// Fluent API: get a table-scoped handle
	Table(index int) (Table, error)
}
