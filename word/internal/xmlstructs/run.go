package xmlstructs

import "encoding/xml"

// Run defines a run of text within a paragraph
type Run struct {
	XMLName xml.Name       `xml:"w:r"`
	RPr     *RunProperties `xml:"w:rPr,omitempty"`
	T       string         `xml:"w:t,omitempty"`
	Drawing *Drawing       `xml:"w:drawing,omitempty"`
	Br      *Break         `xml:"w:br,omitempty"`
}

type Break struct {
	XMLName xml.Name `xml:"w:br"`
	Type    string   `xml:"w:type,attr,omitempty"`
}

// Drawing defines a drawing object (image, etc.)
type Drawing struct {
	XMLName xml.Name `xml:"w:drawing"`
	Inline  *Inline  `xml:"wp:inline"`
}

// Inline defines an inline drawing object
type Inline struct {
	XMLName xml.Name `xml:"wp:inline"`
	Extent  Extent   `xml:"wp:extent"`
	DocPr   DocPr    `xml:"wp:docPr"`
	Graphic Graphic  `xml:"a:graphic"`
}

// Extent defines the size of the drawing
type Extent struct {
	XMLName xml.Name `xml:"wp:extent"`
	CX      int64    `xml:"cx,attr"`
	CY      int64    `xml:"cy,attr"`
}

// DocPr defines drawing properties
type DocPr struct {
	XMLName xml.Name `xml:"wp:docPr"`
	ID      int      `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

// Graphic defines a graphic object
type Graphic struct {
	XMLName xml.Name    `xml:"a:graphic"`
	Data    GraphicData `xml:"a:graphicData"`
}

// GraphicData defines graphic data
type GraphicData struct {
	XMLName xml.Name `xml:"a:graphicData"`
	URI     string   `xml:"uri,attr"`
	Pic     Pic      `xml:"pic:pic"`
}

// Pic defines a picture object
type Pic struct {
	XMLName  xml.Name `xml:"pic:pic"`
	BlipFill BlipFill `xml:"pic:blipFill"`
}

// BlipFill defines blip fill properties
type BlipFill struct {
	XMLName xml.Name `xml:"pic:blipFill"`
	Blip    Blip     `xml:"pic:blip"`
}

// Blip defines a blip (image)
type Blip struct {
	XMLName xml.Name `xml:"a:blip"`
	Embed   string   `xml:"r:embed,attr"`
}

// RunProperties defines formatting for a text run
type RunProperties struct {
	XMLName xml.Name  `xml:"w:rPr"`
	Bold    *struct{} `xml:"w:b,omitempty"`
	Italic  *struct{} `xml:"w:i,omitempty"`
}
