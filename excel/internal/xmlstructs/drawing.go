package xmlstructs

import "encoding/xml"

// WsDr defines the worksheet drawing
type WsDr struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing wsDr"`
	XMLNS_R string   `xml:"xmlns:r,attr"`
	Anchors []Anchor `xml:",any"`
}

type Anchor struct {
	TwoCellAnchor *TwoCellAnchor `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing twoCellAnchor,omitempty"`
	OneCellAnchor *OneCellAnchor `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing oneCellAnchor,omitempty"`
}

type TwoCellAnchor struct {
	EditAs     string `xml:"editAs,attr,omitempty"`
	From       Marker `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing from"`
	To         Marker `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing to"`
	Pic        *Pic   `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing pic,omitempty"`
	Sp         *Any   `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing sp,omitempty"`
	ClientData *Any   `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing clientData"`
}

type OneCellAnchor struct {
	From       Marker `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing from"`
	Ext        Extent `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing ext"`
	Pic        *Pic   `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing pic,omitempty"`
	Sp         *Any   `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing sp,omitempty"`
	ClientData *Any   `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing clientData"`
}

type Marker struct {
	Col    int   `xml:"col"`
	ColOff int64 `xml:"colOff"`
	Row    int   `xml:"row"`
	RowOff int64 `xml:"rowOff"`
}

type Extent struct {
	Cx int64 `xml:"cx,attr"`
	Cy int64 `xml:"cy,attr"`
}

type Pic struct {
	NvPicPr  NvPicPr  `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing nvPicPr"`
	BlipFill BlipFill `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing blipFill"`
	SpPr     SpPr     `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing spPr"`
}

type NvPicPr struct {
	CNvPr    CNvPr `xml:"cNvPr"`
	CNvPicPr Any   `xml:"cNvPicPr"`
}

type CNvPr struct {
	ID    int    `xml:"id,attr"`
	Name  string `xml:"name,attr"`
	Descr string `xml:"descr,attr,omitempty"`
}

type BlipFill struct {
	Blip    Blip    `xml:"blip"`
	Stretch Stretch `xml:"stretch"`
}

type Blip struct {
	Embed string `xml:"r:embed,attr"`
}

type Stretch struct {
	FillRect Any `xml:"fillRect"`
}

type SpPr struct {
	Xfrm     Xfrm     `xml:"xfrm"`
	PrstGeom PrstGeom `xml:"prstGeom"`
}

type Xfrm struct {
	Off Point  `xml:"off"`
	Ext Extent `xml:"ext"`
}

type Point struct {
	X int64 `xml:"x,attr"`
	Y int64 `xml:"y,attr"`
}

type PrstGeom struct {
	Prst  string `xml:"prst,attr"`
	AvLst Any    `xml:"avLst"`
}

type Any struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
}
