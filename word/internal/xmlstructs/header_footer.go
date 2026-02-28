package xmlstructs

import "encoding/xml"

type Header struct {
	XMLName xml.Name `xml:"w:hdr"`
	W       string   `xml:"xmlns:w,attr"`
	R       string   `xml:"xmlns:r,attr"`
	O       string   `xml:"xmlns:o,attr,omitempty"`
	V       string   `xml:"xmlns:v,attr,omitempty"`
	Content []any    `xml:",any"`
}

type Footer struct {
	XMLName xml.Name `xml:"w:ftr"`
	W       string   `xml:"xmlns:w,attr"`
	R       string   `xml:"xmlns:r,attr"`
	O       string   `xml:"xmlns:o,attr,omitempty"`
	V       string   `xml:"xmlns:v,attr,omitempty"`
	Content []any    `xml:",any"`
}

type HeaderReference struct {
	XMLName xml.Name `xml:"w:headerReference"`
	Type    string   `xml:"w:type,attr"`
	ID      string   `xml:"r:id,attr"`
}

type FooterReference struct {
	XMLName xml.Name `xml:"w:footerReference"`
	Type    string   `xml:"w:type,attr"`
	ID      string   `xml:"r:id,attr"`
}
