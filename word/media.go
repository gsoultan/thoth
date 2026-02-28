package word

import (
	"fmt"
	"os"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func (p *processor) InsertImage(path string, width, height float64, style ...document.CellStyle) error {
	if path == "" {
		return fmt.Errorf("image path cannot be empty")
	}
	if width <= 0 || height <= 0 {
		return fmt.Errorf("image width and height must be positive")
	}
	par, err := p.createImageParagraphInternal(path, width, height)
	if err != nil {
		return err
	}
	if len(style) > 0 {
		if style[0].Horizontal != "" {
			val := style[0].Horizontal
			if val == "justify" {
				val = "both"
			}
			par.PPr.Jc = &xmlstructs.Justification{Val: val}
		}
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)

	return nil
}

func (p *processor) createImageParagraphInternal(path string, width, height float64) (xmlstructs.Paragraph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return xmlstructs.Paragraph{}, fmt.Errorf("read image file: %w", err)
	}

	// 1. Add to media
	imgName := fmt.Sprintf("image%d.png", len(p.media)+1)
	mediaPath := "word/media/" + imgName

	// 2. Add relationship
	if p.docRels == nil {
		p.docRels = &xmlstructs.Relationships{}
	}
	rID := p.docRels.AddRelationship(
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/image",
		"media/"+imgName,
	)

	p.media[mediaPath] = data

	// 3. Create drawing
	emuW := int64(width * 12700)
	emuH := int64(height * 12700)

	drawing := &xmlstructs.Drawing{
		Inline: &xmlstructs.Inline{
			Extent:       xmlstructs.Extent{CX: emuW, CY: emuH},
			EffectExtent: &xmlstructs.EffectExtent{L: 0, T: 0, R: 0, B: 0},
			DocPr:        xmlstructs.DocPr{ID: len(p.media), Name: imgName},
			Graphic: xmlstructs.Graphic{
				Data: xmlstructs.GraphicData{
					URI: "http://schemas.openxmlformats.org/drawingml/2006/picture",
					Pic: xmlstructs.Pic{
						NvPicPr: xmlstructs.NvPicPr{
							CNvPr: xmlstructs.CNvPr{
								ID:   len(p.media),
								Name: imgName,
							},
						},
						BlipFill: xmlstructs.BlipFill{
							Blip: xmlstructs.Blip{Embed: rID},
						},
						SpPr: xmlstructs.SpPr{
							Xfrm: xmlstructs.Xfrm{
								Ext: xmlstructs.Extent{CX: emuW, CY: emuH},
							},
							PrstGeom: xmlstructs.PrstGeom{Prst: "rect"},
						},
					},
				},
			},
		},
	}

	return xmlstructs.Paragraph{
		PPr: &xmlstructs.ParagraphProperties{},
		Content: []any{
			xmlstructs.Run{Drawing: drawing},
		},
	}, nil
}

func (p *processor) SetWatermark(text string, style ...document.CellStyle) error {
	if text == "" {
		return fmt.Errorf("watermark text cannot be empty")
	}
	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	if p.xmlDoc.Body.SectPr == nil {
		p.xmlDoc.Body.SectPr = &xmlstructs.SectPr{}
	}

	var headerID string
	if len(p.headers) == 0 {
		p.SetHeader("")
		for id := range p.headers {
			headerID = id
			break
		}
	} else {
		for id := range p.headers {
			headerID = id
			break
		}
	}

	header := p.headers[headerID]
	header.O = "urn:schemas-microsoft-com:office:office"
	header.V = "urn:schemas-microsoft-com:vml"

	color := "silver"
	if len(style) > 0 && style[0].Color != "" {
		color = style[0].Color
	}

	vml := fmt.Sprintf(`<v:shapetype id="_x0000_t136" coordsize="21600,21600" o:spt="136" adj="10800" path="m@7,l@8,m@5,21600l@6,21600e">
		<v:formulas>
			<v:f eqn="if lineDrawn pixelLineWidth 0"/>
			<v:f eqn="sum @0 1 0"/>
			<v:f eqn="sum 0 0 @1"/>
			<v:f eqn="prod @2 1 2"/>
			<v:f eqn="prod @3 21600 pixelWidth"/>
			<v:f eqn="prod @3 21600 pixelHeight"/>
			<v:f eqn="sum @0 0 1"/>
			<v:f eqn="prod @6 1 2"/>
			<v:f eqn="prod @7 1 2"/>
			<v:f eqn="sum @8 21600 0"/>
			<v:f eqn="prod @9 1 2"/>
		</v:formulas>
		<v:path textpathok="t" o:connecttype="custom" o:connectlocs="@10,0;@10,21600" o:connectangles="270,90"/>
		<v:textpath on="t" fitshape="t"/>
		<v:handles>
			<v:h position="#0,bottomRight" xrange="6629,14971"/>
		</v:handles>
		<o:lock v:ext="edit" text="t" shapetype="t"/>
	</v:shapetype>
	<v:shape id="Watermark" o:spid="_x0000_s1025" type="#_x0000_t136" 
		style="position:absolute;margin-left:0;margin-top:0;width:412.4pt;height:137.45pt;z-index:-251651072;mso-wrap-edited:f;mso-width-percent:0;mso-height-percent:0;mso-position-horizontal:center;mso-position-horizontal-relative:margin;mso-position-vertical:center;mso-position-vertical-relative:margin" 
		fillcolor="%s" stroked="f">
		<v:fill opacity=".5"/>
		<v:textpath style="font-family:&quot;Calibri&quot;;font-size:1pt" string="%s"/>
	</v:shape>`, color, text)

	par := xmlstructs.Paragraph{
		Content: []any{
			xmlstructs.Run{
				Pict: &xmlstructs.Pict{
					Content: vml,
				},
			},
		},
	}
	header.Content = append(header.Content, par)

	return nil
}
