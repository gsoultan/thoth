package xmlstructs

// Cell defines a cell in a row
type Cell struct {
	R string `xml:"r,attr"`
	S int    `xml:"s,attr,omitempty"`
	T string `xml:"t,attr,omitempty"`
	V string `xml:"v,omitempty"`
}
