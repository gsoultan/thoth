package xmlstructs

// Xf defines a single cell format
type Xf struct {
	NumFmtID          int        `xml:"numFmtId,attr"`
	FontID            int        `xml:"fontId,attr"`
	FillID            int        `xml:"fillId,attr"`
	BorderID          int        `xml:"borderId,attr"`
	XfID              *int       `xml:"xfId,attr,omitempty"`
	ApplyNumberFormat int        `xml:"applyNumberFormat,attr,omitempty"`
	ApplyFont         int        `xml:"applyFont,attr,omitempty"`
	ApplyFill         int        `xml:"applyFill,attr,omitempty"`
	ApplyBorder       int        `xml:"applyBorder,attr,omitempty"`
	ApplyAlignment    int        `xml:"applyAlignment,attr,omitempty"`
	Alignment         *Alignment `xml:"alignment,omitempty"`
	Font              *Font      `xml:"font,omitempty"`   // For DXF
	Fill              *Fill      `xml:"fill,omitempty"`   // For DXF
	Border            *Border    `xml:"border,omitempty"` // For DXF
}

type Alignment struct {
	Horizontal string `xml:"horizontal,attr,omitempty"`
	Vertical   string `xml:"vertical,attr,omitempty"`
	Indent     int    `xml:"indent,attr,omitempty"`
	WrapText   int    `xml:"wrapText,attr,omitempty"`
}
