package xmlstructs

type Borders struct {
	Count int      `xml:"count,attr"`
	Items []Border `xml:"border"`
}
