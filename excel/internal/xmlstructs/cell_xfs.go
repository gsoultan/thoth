package xmlstructs

// CellXfs defines cell format settings
type CellXfs struct {
	Count int  `xml:"count,attr"`
	Items []Xf `xml:"xf"`
}
