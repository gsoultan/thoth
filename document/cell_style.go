package document

// CellStyle represents the visual formatting of a spreadsheet cell.
type CellStyle struct {
	Bold          bool
	Italic        bool
	Size          int
	Color         string // Hexadecimal color code, e.g., "FF0000"
	Background    string // Hexadecimal color code
	Border        bool
	BorderTop     bool
	BorderBottom  bool
	BorderLeft    bool
	BorderRight   bool
	BorderWidth   float64
	BorderColor   string
	Horizontal    string // "left", "center", "right"
	Vertical      string // "top", "center", "bottom"
	NumberFormat  string // e.g., "0.00", "mm-dd-yy", "#,##0.00"
	Name          string // Named style, e.g., "Heading1"
	Link          string // Hyperlink URL
	Font          string // Custom font name
	Alt           string // Alternative text for accessibility
	LineSpacing   float64
	SpacingBefore float64
	SpacingAfter  float64
	Indent        float64   // Left indent in points
	Hanging       float64   // Hanging indent in points
	Padding       float64   // Cell padding in points
	KeepWithNext  bool      // Ensure paragraph stays on same page as next item
	KeepTogether  bool      // Ensure paragraph doesn't break across pages
	WrapText      bool      // Wrap text within a cell (Excel/Word)
	Superscript   bool      // Render as superscript
	Subscript     bool      // Render as subscript
	DashPattern   []float64 // PDF dash pattern, e.g., [3, 3] for dashed
	Opacity       float64   // Opacity from 0.0 to 1.0 (PDF)
	Absolute      bool      // If true, the item is positioned relative to the current flow (PDF)
	X             float64   // Absolute X position (PDF)
	Y             float64   // Absolute Y position (PDF)
}

// TextSpan represents a styled span within a rich text paragraph.
type TextSpan struct {
	Text  string
	Style CellStyle
}
