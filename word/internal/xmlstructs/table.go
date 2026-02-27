package xmlstructs

import "encoding/xml"

// Table defines the structure of a w:tbl element
type Table struct {
	XMLName xml.Name         `xml:"w:tbl"`
	TblPr   *TableProperties `xml:"w:tblPr,omitempty"`
	Rows    []TableRow       `xml:"w:tr"`
}

// TableProperties defines table properties
type TableProperties struct {
	XMLName  xml.Name  `xml:"w:tblPr"`
	TblStyle *TblStyle `xml:"w:tblStyle,omitempty"`
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
	Cells   []TableCell         `xml:"w:tc"`
}

type TableRowProperties struct {
	XMLName xml.Name `xml:"w:trPr"`
}

// TableCell defines a cell within a table row (w:tc)
type TableCell struct {
	XMLName    xml.Name             `xml:"w:tc"`
	TcPr       *TableCellProperties `xml:"w:tcPr,omitempty"`
	Paragraphs []Paragraph          `xml:"w:p"`
}

type TableCellProperties struct {
	XMLName   xml.Name          `xml:"w:tcPr"`
	TcW       *TableCellWidth   `xml:"w:tcW,omitempty"`
	TcBorders *TableCellBorders `xml:"w:tcBorders,omitempty"`
	Shd       *TableCellShading `xml:"w:shd,omitempty"`
	VMerge    *VMerge           `xml:"w:vMerge,omitempty"`
	GridSpan  *GridSpan         `xml:"w:gridSpan,omitempty"`
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
