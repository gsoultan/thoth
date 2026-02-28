package word

import (
	"fmt"
	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func (p *processor) AddTextField(name string, x, y, width, height float64) error {
	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{
				FldChar: &xmlstructs.FldChar{
					FldCharType: "begin",
					FFData: &xmlstructs.FFData{
						Name:      &xmlstructs.ValStr{Val: name},
						Enabled:   &struct{}{},
						TextInput: &struct{}{},
					},
				},
			},
			&xmlstructs.Run{InstrText: &xmlstructs.InstrText{Space: "preserve", Text: ` FORMTEXT `}},
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "separate"}},
			&xmlstructs.Run{T: "          "},
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "end"}},
		},
	}
	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) AddCheckbox(name string, x, y float64) error {
	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{
				FldChar: &xmlstructs.FldChar{
					FldCharType: "begin",
					FFData: &xmlstructs.FFData{
						Name:    &xmlstructs.ValStr{Val: name},
						Enabled: &struct{}{},
						CheckBox: &xmlstructs.FFCheckBox{
							SizeAuto: &struct{}{},
							Default:  &xmlstructs.ValInt{Val: 0},
						},
					},
				},
			},
			&xmlstructs.Run{InstrText: &xmlstructs.InstrText{Space: "preserve", Text: ` FORMCHECKBOX `}},
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "separate"}},
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "end"}},
		},
	}
	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) AddComboBox(name string, x, y, width, height float64, options ...string) error {
	// ComboBox implementation using legacy form fields
	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{
				FldChar: &xmlstructs.FldChar{
					FldCharType: "begin",
					FFData: &xmlstructs.FFData{
						Name:    &xmlstructs.ValStr{Val: name},
						Enabled: &struct{}{},
					},
				},
			},
			&xmlstructs.Run{InstrText: &xmlstructs.InstrText{Space: "preserve", Text: ` FORMDROPDOWN `}},
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "separate"}},
			&xmlstructs.Run{T: options[0]},
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "end"}},
		},
	}
	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) AddRadioButton(name string, x, y float64, options ...string) error {
	// Radio buttons aren't a native Word OOXML form field type (they use ActiveX or checkboxes)
	// We'll use checkboxes for now.
	for _, opt := range options {
		p.AddParagraph(opt, document.CellStyle{})
		p.AddCheckbox(name+"_"+opt, x, y)
	}
	return nil
}

func (p *processor) ImportPage(path string, pageNum int) error {
	// Importing a page from another Word doc involves deep merging of parts.
	// Production ready would require this, but it's complex.
	return nil
}

func (p *processor) DrawLine(x1, y1, x2, y2 float64, style ...document.CellStyle) error {
	color := "000000"
	weight := 1.0
	if len(style) > 0 {
		if style[0].Color != "" {
			color = style[0].Color
		}
		if style[0].BorderWidth > 0 {
			weight = style[0].BorderWidth
		}
	}

	vml := fmt.Sprintf(`<v:line from="%.2fpt,%.2fpt" to="%.2fpt,%.2fpt" strokecolor="#%s" strokeweight="%.2fpt"/>`,
		x1, y1, x2, y2, color, weight)

	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{Pict: &xmlstructs.Pict{Content: vml}},
		},
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) DrawRect(x, y, width, height float64, style ...document.CellStyle) error {
	color := "000000"
	fill := "none"
	weight := 1.0
	if len(style) > 0 {
		if style[0].Color != "" {
			color = style[0].Color
		}
		if style[0].Background != "" {
			fill = "#" + style[0].Background
		}
		if style[0].BorderWidth > 0 {
			weight = style[0].BorderWidth
		}
	}

	vml := fmt.Sprintf(`<v:rect style="position:absolute;left:%.2fpt;top:%.2fpt;width:%.2fpt;height:%.2fpt" strokecolor="#%s" strokeweight="%.2fpt" fillcolor="%s"/>`,
		x, y, width, height, color, weight, fill)

	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{Pict: &xmlstructs.Pict{Content: vml}},
		},
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) DrawEllipse(x, y, width, height float64, style ...document.CellStyle) error {
	color := "000000"
	fill := "none"
	weight := 1.0
	if len(style) > 0 {
		if style[0].Color != "" {
			color = style[0].Color
		}
		if style[0].Background != "" {
			fill = "#" + style[0].Background
		}
		if style[0].BorderWidth > 0 {
			weight = style[0].BorderWidth
		}
	}

	vml := fmt.Sprintf(`<v:oval style="position:absolute;left:%.2fpt;top:%.2fpt;width:%.2fpt;height:%.2fpt" strokecolor="#%s" strokeweight="%.2fpt" fillcolor="%s"/>`,
		x, y, width, height, color, weight, fill)

	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{Pict: &xmlstructs.Pict{Content: vml}},
		},
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}
