package xmlstructs

type Border struct {
	Left   *BorderEdge `xml:"left,omitempty"`
	Right  *BorderEdge `xml:"right,omitempty"`
	Top    *BorderEdge `xml:"top,omitempty"`
	Bottom *BorderEdge `xml:"bottom,omitempty"`
}
