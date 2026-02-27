package xmlstructs

// Fonts defines the fonts list in styles
type Fonts struct {
	Count int    `xml:"count,attr"`
	Items []Font `xml:"font"`
}
