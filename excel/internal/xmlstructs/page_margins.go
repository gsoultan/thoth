package xmlstructs

// PageMargins defines worksheet margins
type PageMargins struct {
	Left   float64 `xml:"left,attr"`
	Right  float64 `xml:"right,attr"`
	Top    float64 `xml:"top,attr"`
	Bottom float64 `xml:"bottom,attr"`
	Header float64 `xml:"header,attr"`
	Footer float64 `xml:"footer,attr"`
}
