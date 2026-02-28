package xmlstructs

import "encoding/xml"

// Table defines the structure of xl/tables/table[n].xml
type Table struct {
	XMLName        xml.Name        `xml:"http://schemas.openxmlformats.org/spreadsheetml/2006/main table"`
	ID             int             `xml:"id,attr"`
	Name           string          `xml:"name,attr"`
	DisplayName    string          `xml:"displayName,attr"`
	Ref            string          `xml:"ref,attr"`
	TotalsRowShown int             `xml:"totalsRowShown,attr,omitempty"`
	AutoFilter     *AutoFilter     `xml:"autoFilter,omitempty"`
	TableColumns   TableColumns    `xml:"tableColumns"`
	TableStyleInfo *TableStyleInfo `xml:"tableStyleInfo,omitempty"`
}

type TableColumns struct {
	Count int           `xml:"count,attr"`
	Items []TableColumn `xml:"tableColumn"`
}

type TableColumn struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type TableStyleInfo struct {
	Name              string `xml:"name,attr,omitempty"`
	ShowFirstColumn   int    `xml:"showFirstColumn,attr,omitempty"`
	ShowLastColumn    int    `xml:"showLastColumn,attr,omitempty"`
	ShowRowStripes    int    `xml:"showRowStripes,attr,omitempty"`
	ShowColumnStripes int    `xml:"showColumnStripes,attr,omitempty"`
}
