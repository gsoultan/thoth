package xmlstructs

type Fills struct {
	Count int    `xml:"count,attr"`
	Items []Fill `xml:"fill"`
}
