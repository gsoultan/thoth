package document

// PageSettings represents document layout settings.
type PageSettings struct {
	Orientation Orientation
	PaperType   PaperType
	Margins     Margins
	Columns     int     // Number of columns (default 1)
	ColumnGap   float64 // Gap between columns in points
}
