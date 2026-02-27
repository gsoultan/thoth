package xmlstructs

import "encoding/xml"

// Paragraph defines a paragraph in the document body
type Paragraph struct {
	XMLName xml.Name             `xml:"w:p"`
	PPr     *ParagraphProperties `xml:"w:pPr,omitempty"`
	Runs    []Run                `xml:"w:r"`
}

type ParagraphProperties struct {
	XMLName xml.Name        `xml:"w:pPr"`
	PStyle  *ParagraphStyle `xml:"w:pStyle,omitempty"`
	Jc      *Justification  `xml:"w:jc,omitempty"`
	SectPr  *SectPr         `xml:"w:sectPr,omitempty"`
}

type ParagraphStyle struct {
	XMLName xml.Name `xml:"w:pStyle"`
	Val     string   `xml:"w:val,attr"`
}

type Justification struct {
	XMLName xml.Name `xml:"w:jc"`
	Val     string   `xml:"w:val,attr"`
}
