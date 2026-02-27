package document

// CellStyle represents the visual formatting of a spreadsheet cell.
type CellStyle struct {
	Bold         bool
	Italic       bool
	Size         int
	Color        string // Hexadecimal color code, e.g., "FF0000"
	Background   string // Hexadecimal color code
	Border       bool
	Horizontal   string // "left", "center", "right"
	Vertical     string // "top", "center", "bottom"
	NumberFormat string // e.g., "0.00", "mm-dd-yy", "#,##0.00"
	Name         string // Named style, e.g., "Heading1"
	Link         string // Hyperlink URL
}
