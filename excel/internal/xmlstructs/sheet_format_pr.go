package xmlstructs

// SheetFormatPr defines the structure of xl/worksheets/sheet[n].xml sheetFormatPr element.
type SheetFormatPr struct {
	DefaultRowHeight float64 `xml:"defaultRowHeight,attr,omitempty"`
	OutlineLevelRow  uint8   `xml:"outlineLevelRow,attr,omitempty"`
	OutlineLevelCol  uint8   `xml:"outlineLevelCol,attr,omitempty"`
}
