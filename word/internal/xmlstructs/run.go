package xmlstructs

import "encoding/xml"

// Run defines a run of text within a paragraph
type Run struct {
	XMLName           xml.Name           `xml:"w:r"`
	RPr               *RunProperties     `xml:"w:rPr,omitempty"`
	T                 string             `xml:"w:t,omitempty"`
	Drawing           *Drawing           `xml:"w:drawing,omitempty"`
	Br                *Break             `xml:"w:br,omitempty"`
	FootnoteRef       *FootnoteRef       `xml:"w:footnoteRef,omitempty"`
	FootnoteReference *FootnoteReference `xml:"w:footnoteReference,omitempty"`
	Pict              *Pict              `xml:"w:pict,omitempty"`
	FldChar           *FldChar           `xml:"w:fldChar,omitempty"`
	InstrText         *InstrText         `xml:"w:instrText,omitempty"`
}

type FldChar struct {
	XMLName     xml.Name `xml:"w:fldChar"`
	FldCharType string   `xml:"w:fldCharType,attr"`
	FFData      *FFData  `xml:"w:ffData,omitempty"`
}

type InstrText struct {
	Space string `xml:"xml:space,attr,omitempty"`
	Text  string `xml:",chardata"`
}

type FFData struct {
	XMLName   xml.Name    `xml:"w:ffData"`
	Name      *ValStr     `xml:"w:name,omitempty"`
	Enabled   *struct{}   `xml:"w:enabled,omitempty"`
	TextInput *struct{}   `xml:"w:textInput,omitempty"`
	CheckBox  *FFCheckBox `xml:"w:checkBox,omitempty"`
}

type FFCheckBox struct {
	XMLName  xml.Name  `xml:"w:checkBox"`
	SizeAuto *struct{} `xml:"w:sizeAuto,omitempty"`
	Default  *ValInt   `xml:"w:default,omitempty"`
}

type Pict struct {
	XMLName xml.Name `xml:"w:pict"`
	Content string   `xml:",innerxml"`
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
	XMLName      xml.Name      `xml:"wp:inline"`
	DistT        int           `xml:"distT,attr"`
	DistB        int           `xml:"distB,attr"`
	DistL        int           `xml:"distL,attr"`
	DistR        int           `xml:"distR,attr"`
	Extent       Extent        `xml:"wp:extent"`
	EffectExtent *EffectExtent `xml:"wp:effectExtent,omitempty"`
	DocPr        DocPr         `xml:"wp:docPr"`
	Graphic      Graphic       `xml:"a:graphic"`
}

type EffectExtent struct {
	XMLName xml.Name `xml:"wp:effectExtent"`
	L       int64    `xml:"l,attr"`
	T       int64    `xml:"t,attr"`
	R       int64    `xml:"r,attr"`
	B       int64    `xml:"b,attr"`
}

// Extent defines the size of the drawing
type Extent struct {
	CX int64 `xml:"cx,attr"`
	CY int64 `xml:"cy,attr"`
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
	NvPicPr  NvPicPr  `xml:"pic:nvPicPr"`
	BlipFill BlipFill `xml:"pic:blipFill"`
	SpPr     SpPr     `xml:"pic:spPr"`
}

type NvPicPr struct {
	XMLName  xml.Name `xml:"pic:nvPicPr"`
	CNvPr    CNvPr    `xml:"pic:cNvPr"`
	CNvPicPr struct{} `xml:"pic:cNvPicPr"`
}

type CNvPr struct {
	XMLName xml.Name `xml:"pic:cNvPr"`
	ID      int      `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type SpPr struct {
	XMLName  xml.Name `xml:"pic:spPr"`
	Xfrm     Xfrm     `xml:"a:xfrm"`
	PrstGeom PrstGeom `xml:"a:prstGeom"`
}

type Xfrm struct {
	XMLName xml.Name `xml:"a:xfrm"`
	Off     Off      `xml:"a:off"`
	Ext     Extent   `xml:"a:ext"`
}

type Off struct {
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
}

type PrstGeom struct {
	XMLName xml.Name `xml:"a:prstGeom"`
	Prst    string   `xml:"prst,attr"`
	AvLst   struct{} `xml:"a:avLst"`
}

// BlipFill defines blip fill properties
type BlipFill struct {
	XMLName xml.Name `xml:"pic:blipFill"`
	Blip    Blip     `xml:"a:blip"`
	Stretch Stretch  `xml:"a:stretch"`
}

type Stretch struct {
	XMLName  xml.Name `xml:"a:stretch"`
	FillRect struct{} `xml:"a:fillRect"`
}

// Blip defines a blip (image)
type Blip struct {
	XMLName xml.Name `xml:"a:blip"`
	Embed   string   `xml:"r:embed,attr"`
}

// RunProperties defines formatting for a text run
type RunProperties struct {
	XMLName   xml.Name          `xml:"w:rPr"`
	RStyle    *RStyle           `xml:"w:rStyle,omitempty"`
	Bold      *struct{}         `xml:"w:b,omitempty"`
	Italic    *struct{}         `xml:"w:i,omitempty"`
	Sz        *ValInt           `xml:"w:sz,omitempty"`
	SzCs      *ValInt           `xml:"w:szCs,omitempty"`
	Color     *Color            `xml:"w:color,omitempty"`
	VertAlign *ValStr           `xml:"w:vertAlign,omitempty"`
	U         *Underline        `xml:"w:u,omitempty"`
	Shd       *TableCellShading `xml:"w:shd,omitempty"`
}

type Underline struct {
	XMLName xml.Name `xml:"w:u"`
	Val     string   `xml:"w:val,attr"`
	Color   string   `xml:"w:color,attr,omitempty"`
}

type Color struct {
	Val string `xml:"w:val,attr"`
}

type RStyle struct {
	XMLName xml.Name `xml:"w:rStyle"`
	Val     string   `xml:"w:val,attr"`
}
