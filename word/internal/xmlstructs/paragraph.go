package xmlstructs

import "encoding/xml"

// Paragraph defines a paragraph in the document body
type Paragraph struct {
	XMLName xml.Name             `xml:"w:p"`
	PPr     *ParagraphProperties `xml:"w:pPr,omitempty"`
	Content []any                `xml:",any"`
}

type Hyperlink struct {
	XMLName xml.Name `xml:"w:hyperlink"`
	ID      string   `xml:"r:id,attr,omitempty"`
	Anchor  string   `xml:"w:anchor,attr,omitempty"`
	Runs    []*Run   `xml:"w:r"`
}

type BookmarkStart struct {
	XMLName xml.Name `xml:"w:bookmarkStart"`
	ID      int      `xml:"w:id,attr"`
	Name    string   `xml:"w:name,attr"`
}

type BookmarkEnd struct {
	XMLName xml.Name `xml:"w:bookmarkEnd"`
	ID      int      `xml:"w:id,attr"`
}

type ParagraphProperties struct {
	XMLName   xml.Name        `xml:"w:pPr"`
	PStyle    *ParagraphStyle `xml:"w:pStyle,omitempty"`
	Jc        *Justification  `xml:"w:jc,omitempty"`
	Ind       *Ind            `xml:"w:ind,omitempty"`
	NumPr     *NumPr          `xml:"w:numPr,omitempty"`
	Spacing   *Spacing        `xml:"w:spacing,omitempty"`
	KeepNext  *struct{}       `xml:"w:keepNext,omitempty"`
	KeepLines *struct{}       `xml:"w:keepLines,omitempty"`
	SectPr    *SectPr         `xml:"w:sectPr,omitempty"`
}

type Spacing struct {
	XMLName  xml.Name `xml:"w:spacing"`
	Before   int      `xml:"w:before,attr,omitempty"`
	After    int      `xml:"w:after,attr,omitempty"`
	Line     int      `xml:"w:line,attr,omitempty"`
	LineRule string   `xml:"w:lineRule,attr,omitempty"` // "auto", "exact", "atLeast"
}

type Ind struct {
	XMLName xml.Name `xml:"w:ind"`
	Left    int      `xml:"w:left,attr,omitempty"`
	Hanging int      `xml:"w:hanging,attr,omitempty"`
}

type ParagraphStyle struct {
	XMLName xml.Name `xml:"w:pStyle"`
	Val     string   `xml:"w:val,attr"`
}

type Justification struct {
	XMLName xml.Name `xml:"w:jc"`
	Val     string   `xml:"w:val,attr"`
}
