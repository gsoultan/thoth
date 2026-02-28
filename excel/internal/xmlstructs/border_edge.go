package xmlstructs

type BorderEdge struct {
	Style string `xml:"style,attr,omitempty"`
	Color *Color `xml:"color,omitempty"`
}
