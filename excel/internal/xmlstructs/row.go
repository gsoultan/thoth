package xmlstructs

// Row defines a row in the sheet data
type Row struct {
	R            int     `xml:"r,attr"`
	Cells        []Cell  `xml:"c"`
	Ht           float64 `xml:"ht,attr,omitempty"`
	CustomHeight int     `xml:"customHeight,attr,omitempty"`
	OutlineLevel uint8   `xml:"outlineLevel,attr,omitempty"`
	Collapsed    bool    `xml:"collapsed,attr,omitempty"`
}
