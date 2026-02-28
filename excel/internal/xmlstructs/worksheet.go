package xmlstructs

import "encoding/xml"

// Worksheet defines the structure of xl/worksheets/sheet[n].xml
type Worksheet struct {
	XMLName               xml.Name                `xml:"http://schemas.openxmlformats.org/spreadsheetml/2006/main worksheet"`
	XMLNS_R               string                  `xml:"xmlns:r,attr"`
	SheetPr               *SheetPr                `xml:"sheetPr,omitempty"`
	Dimension             *Dimension              `xml:"dimension,omitempty"`
	SheetViews            *SheetViews             `xml:"sheetViews,omitempty"`
	SheetFormatPr         *SheetFormatPr          `xml:"sheetFormatPr,omitempty"`
	Cols                  *Cols                   `xml:"cols,omitempty"`
	SheetData             SheetData               `xml:"sheetData"`
	SheetProtection       *SheetProtection        `xml:"sheetProtection,omitempty"`
	AutoFilter            *AutoFilter             `xml:"autoFilter,omitempty"`
	MergeCells            *MergeCells             `xml:"mergeCells,omitempty"`
	ConditionalFormatting []ConditionalFormatting `xml:"conditionalFormatting,omitempty"`
	DataValidations       *DataValidations        `xml:"dataValidations,omitempty"`
	Hyperlinks            *Hyperlinks             `xml:"hyperlinks,omitempty"`
	PageMargins           *PageMargins            `xml:"pageMargins,omitempty"`
	PageSetup             *PageSetup              `xml:"pageSetup,omitempty"`
	HeaderFooter          *HeaderFooter           `xml:"headerFooter,omitempty"`
	Drawing               *WsDrawing              `xml:"drawing,omitempty"`
	TableParts            *TableParts             `xml:"tableParts,omitempty"`
}

type TableParts struct {
	Count int         `xml:"count,attr"`
	Items []TablePart `xml:"tablePart"`
}

type TablePart struct {
	RID string `xml:"r:id,attr"`
}

type HeaderFooter struct {
	OddHeader string `xml:"oddHeader,omitempty"`
	OddFooter string `xml:"oddFooter,omitempty"`
}

type ConditionalFormatting struct {
	Sqref  string   `xml:"sqref,attr"`
	CfRule []CfRule `xml:"cfRule"`
}

type CfRule struct {
	Type     string   `xml:"type,attr"`
	DxfID    *int     `xml:"dxfId,attr,omitempty"`
	Priority int      `xml:"priority,attr"`
	Operator string   `xml:"operator,attr,omitempty"`
	Formula  []string `xml:"formula,omitempty"`
}

type SheetPr struct {
	OutlinePr *OutlinePr `xml:"outlinePr,omitempty"`
}

type OutlinePr struct {
	SummaryBelow int `xml:"summaryBelow,attr"`
	SummaryRight int `xml:"summaryRight,attr"`
}

type SheetProtection struct {
	Sheet          int    `xml:"sheet,attr"`
	Objects        int    `xml:"objects,attr"`
	Scenarios      int    `xml:"scenarios,attr"`
	Password       string `xml:"password,attr,omitempty"`
	SelectLocked   int    `xml:"selectLockedCells,attr"`
	SelectUnlocked int    `xml:"selectUnlockedCells,attr"`
}

type DataValidations struct {
	Count int              `xml:"count,attr"`
	Items []DataValidation `xml:"dataValidation"`
}

type DataValidation struct {
	Type         string `xml:"type,attr,omitempty"`
	AllowBlank   int    `xml:"allowBlank,attr,omitempty"`
	ShowInputMsg int    `xml:"showInputMessage,attr,omitempty"`
	ShowErrorMsg int    `xml:"showErrorMessage,attr,omitempty"`
	Sqref        string `xml:"sqref,attr"`
	Formula1     string `xml:"formula1,omitempty"`
	Formula2     string `xml:"formula2,omitempty"`
}

type Hyperlinks struct {
	Items []Hyperlink `xml:"hyperlink"`
}

type Hyperlink struct {
	Ref      string `xml:"ref,attr"`
	RID      string `xml:"r:id,attr,omitempty"`
	Location string `xml:"location,attr,omitempty"`
	Display  string `xml:"display,attr,omitempty"`
}

type SheetViews struct {
	Items []SheetView `xml:"sheetView"`
}

type SheetView struct {
	TabSelected    int   `xml:"tabSelected,attr"`
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
	Min          int     `xml:"min,attr"`
	Max          int     `xml:"max,attr"`
	Width        float64 `xml:"width,attr"`
	CustomWidth  int     `xml:"customWidth,attr,omitempty"`
	OutlineLevel uint8   `xml:"outlineLevel,attr,omitempty"`
	Collapsed    bool    `xml:"collapsed,attr,omitempty"`
}

type MergeCells struct {
	Count int         `xml:"count,attr"`
	Items []MergeCell `xml:"mergeCell"`
}

type MergeCell struct {
	Ref string `xml:"ref,attr"`
}
