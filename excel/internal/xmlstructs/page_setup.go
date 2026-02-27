package xmlstructs

// PageSetup defines worksheet page settings
type PageSetup struct {
	Orientation string `xml:"orientation,attr,omitempty"`
	PaperSize   int    `xml:"paperSize,attr,omitempty"`
}
