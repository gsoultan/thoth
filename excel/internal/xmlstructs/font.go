package xmlstructs

// Font defines font properties
type Font struct {
	Bold   *struct{}  `xml:"b,omitempty"`
	Italic *struct{}  `xml:"i,omitempty"`
	Size   *ValInt    `xml:"sz"`
	Color  *Color     `xml:"color,omitempty"`
	Name   *ValString `xml:"name"`
}
