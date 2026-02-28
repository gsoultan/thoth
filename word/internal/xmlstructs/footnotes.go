package xmlstructs

import "encoding/xml"

type Footnotes struct {
	XMLName   xml.Name   `xml:"w:footnotes"`
	W         string     `xml:"xmlns:w,attr"`
	Footnotes []Footnote `xml:"w:footnote"`
}

type Footnote struct {
	XMLName xml.Name    `xml:"w:footnote"`
	ID      int         `xml:"w:id,attr"`
	Type    string      `xml:"w:type,attr,omitempty"`
	Content []Paragraph `xml:"w:p"`
}

type FootnoteReference struct {
	XMLName xml.Name `xml:"w:footnoteReference"`
	ID      int      `xml:"w:id,attr"`
}

type FootnoteRef struct {
	XMLName xml.Name `xml:"w:footnoteRef"`
}
