package xmlstructs

import "encoding/xml"

type Numbering struct {
	XMLName      xml.Name      `xml:"w:numbering"`
	W            string        `xml:"xmlns:w,attr"`
	AbstractNums []AbstractNum `xml:"w:abstractNum"`
	Nums         []Num         `xml:"w:num"`
}

type AbstractNum struct {
	XMLName       xml.Name `xml:"w:abstractNum"`
	AbstractNumID int      `xml:"w:abstractNumId,attr"`
	Levels        []Level  `xml:"w:lvl"`
}

type Level struct {
	XMLName xml.Name             `xml:"w:lvl"`
	ILvl    int                  `xml:"w:ilvl,attr"`
	Start   *ValInt              `xml:"w:start"`
	NumFmt  *ValStr              `xml:"w:numFmt"`
	LvlText *ValStr              `xml:"w:lvlText"`
	LvlJc   *ValStr              `xml:"w:lvlJc"`
	PPr     *ParagraphProperties `xml:"w:pPr,omitempty"`
	RPr     *RunProperties       `xml:"w:rPr,omitempty"`
}

type Num struct {
	XMLName       xml.Name `xml:"w:num"`
	NumID         int      `xml:"w:numId,attr"`
	AbstractNumID *ValInt  `xml:"w:abstractNumId"`
}

type ValStr struct {
	Val string `xml:"w:val,attr"`
}

type NumPr struct {
	XMLName xml.Name `xml:"w:numPr"`
	ILvl    *ValInt  `xml:"w:ilvl"`
	NumID   *ValInt  `xml:"w:numId"`
}
