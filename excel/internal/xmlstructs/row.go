package xmlstructs

// Row defines a row in the sheet data
type Row struct {
	R     int    `xml:"r,attr"`
	Cells []Cell `xml:"c"`
}
