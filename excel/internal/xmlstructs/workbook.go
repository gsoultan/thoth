package xmlstructs

import "encoding/xml"

// Workbook defines the structure of xl/workbook.xml
type Workbook struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/spreadsheetml/2006/main workbook"`
	XMLNS_R string   `xml:"xmlns:r,attr"`
	Sheets  []Sheet  `xml:"sheets>sheet"`
}
