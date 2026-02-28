package xmlstructs

type Border struct {
	Left     BorderEdge `xml:"left"`
	Right    BorderEdge `xml:"right"`
	Top      BorderEdge `xml:"top"`
	Bottom   BorderEdge `xml:"bottom"`
	Diagonal BorderEdge `xml:"diagonal"`
}
