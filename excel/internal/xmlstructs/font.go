package xmlstructs

// Font defines font properties
type Font struct {
	Bold   *struct{} `xml:"b,omitempty"`
	Italic *struct{} `xml:"i,omitempty"`
	Size   *ValInt   `xml:"sz,omitempty"`
	Color  *Color    `xml:"color,omitempty"`
}
