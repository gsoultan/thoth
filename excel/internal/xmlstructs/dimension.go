package xmlstructs

// Dimension defines the structure of xl/worksheets/sheet[n].xml dimension element.
type Dimension struct {
	Ref string `xml:"ref,attr"`
}
