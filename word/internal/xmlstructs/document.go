package xmlstructs

import "encoding/xml"

// Document defines the structure of word/document.xml
type Document struct {
	XMLName xml.Name `xml:"w:document"`
	W       string   `xml:"xmlns:w,attr"`
	R       string   `xml:"xmlns:r,attr"`
	WP      string   `xml:"xmlns:wp,attr"`
	A       string   `xml:"xmlns:a,attr"`
	Pic     string   `xml:"xmlns:pic,attr"`
	O       string   `xml:"xmlns:o,attr,omitempty"`
	V       string   `xml:"xmlns:v,attr,omitempty"`
	W10     string   `xml:"xmlns:w10,attr,omitempty"`
	Body    Body     `xml:"w:body"`
}
