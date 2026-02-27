package xmlstructs

import "encoding/xml"

// Worksheet defines the structure of xl/worksheets/sheet[n].xml
type Worksheet struct {
	XMLName     xml.Name     `xml:"http://schemas.openxmlformats.org/spreadsheetml/2006/main worksheet"`
	SheetViews  *SheetViews  `xml:"sheetViews,omitempty"`
	Cols        *Cols        `xml:"cols,omitempty"`
	SheetData   SheetData    `xml:"sheetData"`
	AutoFilter  *AutoFilter  `xml:"autoFilter,omitempty"`
	MergeCells  *MergeCells  `xml:"mergeCells,omitempty"`
	Drawing     *WsDrawing   `xml:"drawing,omitempty"`
	PageMargins *PageMargins `xml:"pageMargins,omitempty"`
	PageSetup   *PageSetup   `xml:"pageSetup,omitempty"`
}

type SheetViews struct {
	Items []SheetView `xml:"sheetView"`
}

type SheetView struct {
	TabSelected    int   `xml:"tabSelected,attr,omitempty"`
	WorkbookViewID int   `xml:"workbookViewId,attr"`
	Pane           *Pane `xml:"pane,omitempty"`
}

type Pane struct {
	XSplit      int    `xml:"xSplit,attr,omitempty"`
	YSplit      int    `xml:"ySplit,attr,omitempty"`
	TopLeftCell string `xml:"topLeftCell,attr,omitempty"`
	ActivePane  string `xml:"activePane,attr,omitempty"`
	State       string `xml:"state,attr,omitempty"` // "frozen"
}

type AutoFilter struct {
	Ref string `xml:"ref,attr"`
}

type Cols struct {
	Items []Col `xml:"col"`
}

type Col struct {
	Min         int     `xml:"min,attr"`
	Max         int     `xml:"max,attr"`
	Width       float64 `xml:"width,attr"`
	CustomWidth int     `xml:"customWidth,attr,omitempty"`
}

type MergeCells struct {
	Count int         `xml:"count,attr"`
	Items []MergeCell `xml:"mergeCell"`
}

type MergeCell struct {
	Ref string `xml:"ref,attr"`
}
