package xmlstructs

import "encoding/xml"

// Table defines the structure of a w:tbl element
type Table struct {
	XMLName xml.Name         `xml:"w:tbl"`
	TblPr   *TableProperties `xml:"w:tblPr,omitempty"`
	TblGrid *TableGrid       `xml:"w:tblGrid,omitempty"`
	Rows    []*TableRow      `xml:"w:tr"`
}

type TableGrid struct {
	XMLName xml.Name       `xml:"w:tblGrid"`
	Cols    []TableGridCol `xml:"w:gridCol"`
}

type TableGridCol struct {
	XMLName xml.Name `xml:"w:gridCol"`
	W       int      `xml:"w:w,attr"`
}

// TableProperties defines table properties
type TableProperties struct {
	XMLName    xml.Name       `xml:"w:tblPr"`
	TblStyle   *TblStyle      `xml:"w:tblStyle,omitempty"`
	TblW       *TableWidth    `xml:"w:tblW,omitempty"`
	TblInd     *TableIndent   `xml:"w:tblInd,omitempty"`
	TblBorders *TableBorders  `xml:"w:tblBorders,omitempty"`
	TblLayout  *TableLayout   `xml:"w:tblLayout,omitempty"`
	Jc         *Justification `xml:"w:jc,omitempty"`
}

type TableWidth struct {
	XMLName xml.Name `xml:"w:tblW"`
	W       int      `xml:"w:w,attr"`
	Type    string   `xml:"w:type,attr"`
}

type TableIndent struct {
	XMLName xml.Name `xml:"w:tblInd"`
	W       int      `xml:"w:w,attr"`
	Type    string   `xml:"w:type,attr"`
}

type TableBorders struct {
	XMLName xml.Name    `xml:"w:tblBorders"`
	Top     *BorderLine `xml:"w:top,omitempty"`
	Left    *BorderLine `xml:"w:left,omitempty"`
	Bottom  *BorderLine `xml:"w:bottom,omitempty"`
	Right   *BorderLine `xml:"w:right,omitempty"`
	InsideH *BorderLine `xml:"w:insideH,omitempty"`
	InsideV *BorderLine `xml:"w:insideV,omitempty"`
}

type TableLayout struct {
	XMLName xml.Name `xml:"w:tblLayout"`
	Type    string   `xml:"w:type,attr"`
}

// TblStyle defines the table style
type TblStyle struct {
	XMLName xml.Name `xml:"w:tblStyle"`
	Val     string   `xml:"w:val,attr"`
}

// TableRow defines a row within a table (w:tr)
type TableRow struct {
	XMLName xml.Name            `xml:"w:tr"`
	TrPr    *TableRowProperties `xml:"w:trPr,omitempty"`
	Cells   []*TableCell        `xml:"w:tc"`
}

type TableRowProperties struct {
	XMLName   xml.Name       `xml:"w:trPr"`
	TrHeight  *TrHeight      `xml:"w:trHeight,omitempty"`
	TblHeader *struct{}      `xml:"w:tblHeader,omitempty"`
	Jc        *Justification `xml:"w:jc,omitempty"`
}

type TrHeight struct {
	XMLName xml.Name `xml:"w:trHeight"`
	Val     int      `xml:"w:val,attr"`
	HRule   string   `xml:"w:hRule,attr,omitempty"` // "atLeast", "exact"
}

// TableCell defines a cell within a table row (w:tc)
type TableCell struct {
	XMLName xml.Name             `xml:"w:tc"`
	TcPr    *TableCellProperties `xml:"w:tcPr,omitempty"`
	Content []any                `xml:",any"`
}

type TableCellProperties struct {
	XMLName   xml.Name          `xml:"w:tcPr"`
	TcW       *TableCellWidth   `xml:"w:tcW,omitempty"`
	TcBorders *TableCellBorders `xml:"w:tcBorders,omitempty"`
	Shd       *TableCellShading `xml:"w:shd,omitempty"`
	NoWrap    *struct{}         `xml:"w:noWrap,omitempty"`
	VMerge    *VMerge           `xml:"w:vMerge,omitempty"`
	GridSpan  *GridSpan         `xml:"w:gridSpan,omitempty"`
	VAlign    *VAlign           `xml:"w:vAlign,omitempty"`
	TcMar     *TableCellMargins `xml:"w:tcMar,omitempty"`
}

type VAlign struct {
	XMLName xml.Name `xml:"w:vAlign"`
	Val     string   `xml:"w:val,attr"` // "top", "center", "both", "bottom"
}

type TableCellWidth struct {
	XMLName xml.Name `xml:"w:tcW"`
	W       int      `xml:"w:w,attr"`
	Type    string   `xml:"w:type,attr"`
}

type TableCellBorders struct {
	XMLName xml.Name    `xml:"w:tcBorders"`
	Top     *BorderLine `xml:"w:top,omitempty"`
	Left    *BorderLine `xml:"w:left,omitempty"`
	Bottom  *BorderLine `xml:"w:bottom,omitempty"`
	Right   *BorderLine `xml:"w:right,omitempty"`
}

type BorderLine struct {
	Val   string `xml:"w:val,attr"`
	Sz    int    `xml:"w:sz,attr,omitempty"`
	Space int    `xml:"w:space,attr,omitempty"`
	Color string `xml:"w:color,attr,omitempty"`
}

type TableCellShading struct {
	XMLName xml.Name `xml:"w:shd"`
	Val     string   `xml:"w:val,attr"`
	Color   string   `xml:"w:color,attr,omitempty"`
	Fill    string   `xml:"w:fill,attr,omitempty"`
}

type VMerge struct {
	XMLName xml.Name `xml:"w:vMerge"`
	Val     string   `xml:"w:val,attr,omitempty"` // "restart" or empty
}

type GridSpan struct {
	XMLName xml.Name `xml:"w:gridSpan"`
	Val     int      `xml:"w:val,attr"`
}

type TableCellMargins struct {
	XMLName xml.Name    `xml:"w:tcMar"`
	Top     *TableCellW `xml:"w:top,omitempty"`
	Left    *TableCellW `xml:"w:left,omitempty"`
	Bottom  *TableCellW `xml:"w:bottom,omitempty"`
	Right   *TableCellW `xml:"w:right,omitempty"`
}

type TableCellW struct {
	W    int    `xml:"w:w,attr"`
	Type string `xml:"w:type,attr"`
}
