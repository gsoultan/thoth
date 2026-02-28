package document

// WordProcessor defines operations specific to word processing documents (Word).
type WordProcessor interface {
	Document
	AddParagraph(text string, style ...CellStyle) error
	AddRichParagraph(spans []TextSpan) error
	AddHeading(text string, level int, style ...CellStyle) error
	InsertImage(path string, width, height float64, style ...CellStyle) error
	AddTable(rows, cols int) (Table, error)

	AddPageBreak() error
	AddSection(settings PageSettings) error

	SetHeader(text string, style ...CellStyle) error
	SetFooter(text string, style ...CellStyle) error

	DrawLine(x1, y1, x2, y2 float64, style ...CellStyle) error
	DrawRect(x, y, w, h float64, style ...CellStyle) error
	DrawEllipse(x, y, w, h float64, style ...CellStyle) error
	RegisterFont(name, path string) error

	AddTextField(name string, x, y, w, h float64) error
	AddCheckbox(name string, x, y float64) error
	AddComboBox(name string, x, y, w, h float64, options ...string) error
	AddRadioButton(name string, x, y float64, options ...string) error
	ImportPage(path string, pageNum int) error
	AddFootnote(text string) error
	AddList(items []string, ordered bool, style ...CellStyle) error
	SetWatermark(text string, style ...CellStyle) error
	AddHyperlink(text, url string, style ...CellStyle) error
	AddBookmark(name string) error
	AddTableOfContents() error
	AttachFile(path, name, description string) error

	// Fluent API: get a table-scoped handle
	Table(index int) (Table, error)
}
