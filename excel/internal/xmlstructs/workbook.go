package xmlstructs

import "encoding/xml"

// Workbook defines the structure of xl/workbook.xml
type Workbook struct {
	XMLName            xml.Name            `xml:"http://schemas.openxmlformats.org/spreadsheetml/2006/main workbook"`
	XMLNS_R            string              `xml:"xmlns:r,attr"`
	WorkbookPr         *WorkbookPr         `xml:"workbookPr,omitempty"`
	WorkbookProtection *WorkbookProtection `xml:"workbookProtection,omitempty"`
	WorkbookViews      *WorkbookViews      `xml:"bookViews,omitempty"`
	Sheets             []Sheet             `xml:"sheets>sheet"`
	DefinedNames       *DefinedNames       `xml:"definedNames,omitempty"`
	CalcPr             *CalcPr             `xml:"calcPr,omitempty"`
}

type WorkbookPr struct {
	Date1904 int `xml:"date1904,attr"`
}

type WorkbookProtection struct {
	WorkbookPassword string `xml:"workbookPassword,attr,omitempty"`
	LockStructure    int    `xml:"lockStructure,attr,omitempty"`
}

type WorkbookViews struct {
	Items []WorkbookView `xml:"workbookView"`
}

type WorkbookView struct {
	XWindow      int `xml:"xWindow,attr,omitempty"`
	YWindow      int `xml:"yWindow,attr,omitempty"`
	WindowWidth  int `xml:"windowWidth,attr,omitempty"`
	WindowHeight int `xml:"windowHeight,attr,omitempty"`
	ActiveTab    int `xml:"activeTab,attr,omitempty"`
}

type DefinedNames struct {
	Items []DefinedName `xml:"definedName"`
}

type DefinedName struct {
	Name         string `xml:"name,attr"`
	LocalSheetID *int   `xml:"localSheetId,attr,omitempty"`
	Ref          string `xml:",chardata"`
}

type CalcPr struct {
	CalcID         string `xml:"calcId,attr,omitempty"`
	FullCalcOnLoad int    `xml:"fullCalcOnLoad,attr,omitempty"`
}
