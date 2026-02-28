package xmlstructs

// Cell defines a cell in a row
type Cell struct {
	R  string  `xml:"r,attr"`
	S  int     `xml:"s,attr,omitempty"`
	T  string  `xml:"t,attr,omitempty"`
	F  *string `xml:"f,omitempty"`
	V  string  `xml:"v,omitempty"`
	IS *Rst    `xml:"is,omitempty"` // Inline string/Rich text
}

// Rst represents a rich text run or inline string.
type Rst struct {
	T string `xml:"t,omitempty"`
	R []Run  `xml:"r,omitempty"`
}

// Run represents a styled text run.
type Run struct {
	RPr *RPr   `xml:"rPr,omitempty"`
	T   string `xml:"t"`
}

// RPr represents text run properties.
type RPr struct {
	Bold      *struct{}  `xml:"b,omitempty"`
	Italic    *struct{}  `xml:"i,omitempty"`
	Size      *ValInt    `xml:"sz,omitempty"`
	Color     *Color     `xml:"color,omitempty"`
	RFont     *ValString `xml:"rFont,omitempty"`
	Underline *ValString `xml:"u,omitempty"`
}
