package xmlstructs

import "encoding/xml"

// Body contains the paragraphs and other elements of the document
type Body struct {
	XMLName xml.Name `xml:"w:body"`
	Content []any    `xml:",any"`
	SectPr  *SectPr  `xml:"w:sectPr,omitempty"`
}

// SectPr defines section properties
type SectPr struct {
	XMLName    xml.Name          `xml:"w:sectPr"`
	HeaderRefs []HeaderReference `xml:"w:headerReference,omitempty"`
	FooterRefs []FooterReference `xml:"w:footerReference,omitempty"`
	PgSz       *PgSz             `xml:"w:pgSz,omitempty"`
	PgMar      *PgMar            `xml:"w:pgMar,omitempty"`
	PgNumType  *PgNumType        `xml:"w:pgNumType,omitempty"`
	Cols       *Columns          `xml:"w:cols,omitempty"`
	DocGrid    *DocGrid          `xml:"w:docGrid,omitempty"`
	TitlePg    *TitlePg          `xml:"w:titlePg,omitempty"`
}

type TitlePg struct {
	Val string `xml:"w:val,attr,omitempty"`
}

type PgNumType struct {
	XMLName xml.Name `xml:"w:pgNumType"`
	Start   int      `xml:"w:start,attr,omitempty"`
	Fmt     string   `xml:"w:fmt,attr,omitempty"`
}

type Columns struct {
	XMLName xml.Name `xml:"w:cols"`
	Num     int      `xml:"w:num,attr,omitempty"`
	Space   int      `xml:"w:space,attr,omitempty"`
}

type DocGrid struct {
	XMLName   xml.Name `xml:"w:docGrid"`
	LinePitch int      `xml:"w:linePitch,attr,omitempty"`
}

// PgSz defines page size and orientation
type PgSz struct {
	XMLName xml.Name `xml:"w:pgSz"`
	W       int      `xml:"w:w,attr,omitempty"`
	H       int      `xml:"w:h,attr,omitempty"`
	Orient  string   `xml:"w:orient,attr,omitempty"`
}

// PgMar defines page margins
type PgMar struct {
	XMLName xml.Name `xml:"w:pgMar"`
	Top     int      `xml:"w:top,attr,omitempty"`
	Bottom  int      `xml:"w:bottom,attr,omitempty"`
	Left    int      `xml:"w:left,attr,omitempty"`
	Right   int      `xml:"w:right,attr,omitempty"`
}
