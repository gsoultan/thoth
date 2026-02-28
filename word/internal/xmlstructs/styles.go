package xmlstructs

import "encoding/xml"

type Styles struct {
	XMLName xml.Name `xml:"w:styles"`
	W       string   `xml:"xmlns:w,attr"`
	Content []any    `xml:",any"`
}

type Style struct {
	Type    string               `xml:"w:type,attr,omitempty"`
	StyleID string               `xml:"w:styleId,attr,omitempty"`
	Default string               `xml:"w:default,attr,omitempty"`
	Name    *ValStr              `xml:"w:name,omitempty"`
	Next    *ValStr              `xml:"w:next,omitempty"`
	BasedOn *ValStr              `xml:"w:basedOn,omitempty"`
	PPr     *ParagraphProperties `xml:"w:pPr,omitempty"`
	RPr     *RunProperties       `xml:"w:rPr,omitempty"`
	TblPr   *TableProperties     `xml:"w:tblPr,omitempty"`
}

func NewStyles() *Styles {
	return &Styles{
		W: "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
		Content: []any{
			&Style{
				Type:    "paragraph",
				StyleID: "Normal",
				Default: "1",
				Name:    &ValStr{Val: "Normal"},
			},
			&Style{
				Type:    "character",
				StyleID: "DefaultParagraphFont",
				Default: "1",
				Name:    &ValStr{Val: "Default Paragraph Font"},
			},
			&Style{
				Type:    "table",
				StyleID: "TableNormal",
				Default: "1",
				Name:    &ValStr{Val: "Normal Table"},
				TblPr: &TableProperties{
					TblBorders: &TableBorders{
						Top:     &BorderLine{Val: "none", Sz: 0, Color: "auto"},
						Left:    &BorderLine{Val: "none", Sz: 0, Color: "auto"},
						Bottom:  &BorderLine{Val: "none", Sz: 0, Color: "auto"},
						Right:   &BorderLine{Val: "none", Sz: 0, Color: "auto"},
						InsideH: &BorderLine{Val: "none", Sz: 0, Color: "auto"},
						InsideV: &BorderLine{Val: "none", Sz: 0, Color: "auto"},
					},
				},
			},
			&Style{
				Type:    "paragraph",
				StyleID: "Heading1",
				Name:    &ValStr{Val: "heading 1"},
				BasedOn: &ValStr{Val: "Normal"},
				Next:    &ValStr{Val: "Normal"},
				RPr: &RunProperties{
					Bold: &struct{}{},
					Sz:   &ValInt{Val: 32}, // 16pt
					SzCs: &ValInt{Val: 32},
				},
			},
			&Style{
				Type:    "paragraph",
				StyleID: "Heading2",
				Name:    &ValStr{Val: "heading 2"},
				BasedOn: &ValStr{Val: "Normal"},
				Next:    &ValStr{Val: "Normal"},
				RPr: &RunProperties{
					Bold: &struct{}{},
					Sz:   &ValInt{Val: 28}, // 14pt
					SzCs: &ValInt{Val: 28},
				},
			},
		},
	}
}
