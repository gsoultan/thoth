package xmlstructs

// Relationship defines a single relationship in a .rels file
type Relationship struct {
	ID     string `xml:"Id,attr"`
	Type   string `xml:"Type,attr"`
	Target string `xml:"Target,attr"`
}
