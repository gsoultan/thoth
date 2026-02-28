package xmlstructs

import "encoding/xml"

type Settings struct {
	XMLName           xml.Name           `xml:"w:settings"`
	W                 string             `xml:"xmlns:w,attr"`
	EvenAndOddHeaders *EvenAndOddHeaders `xml:"w:evenAndOddHeaders,omitempty"`
}

type EvenAndOddHeaders struct {
	Val string `xml:"w:val,attr,omitempty"`
}

func NewSettings() *Settings {
	return &Settings{
		W: "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
	}
}
