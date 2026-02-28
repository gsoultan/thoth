package word

import (
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func (p *processor) AddFootnote(text string) error {
	if p.footnotes == nil {
		p.footnotes = &xmlstructs.Footnotes{
			W: "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
		}
		// Add default footnotes
		p.footnotes.Footnotes = append(p.footnotes.Footnotes, xmlstructs.Footnote{ID: -1, Type: "separator"})
		p.footnotes.Footnotes = append(p.footnotes.Footnotes, xmlstructs.Footnote{ID: 0, Type: "continuationSeparator"})
	}

	p.footnoteCounter++
	id := p.footnoteCounter

	fn := xmlstructs.Footnote{
		ID: id,
		Content: []xmlstructs.Paragraph{
			{
				PPr: &xmlstructs.ParagraphProperties{PStyle: &xmlstructs.ParagraphStyle{Val: "FootnoteText"}},
				Content: []any{
					&xmlstructs.Run{
						RPr:         &xmlstructs.RunProperties{VertAlign: &xmlstructs.ValStr{Val: "superscript"}},
						FootnoteRef: &xmlstructs.FootnoteRef{},
					},
					&xmlstructs.Run{T: " " + text},
				},
			},
		},
	}
	p.footnotes.Footnotes = append(p.footnotes.Footnotes, fn)

	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{
				RPr:               &xmlstructs.RunProperties{VertAlign: &xmlstructs.ValStr{Val: "superscript"}},
				FootnoteReference: &xmlstructs.FootnoteReference{ID: id},
			},
		},
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) createParagraphWithFields(text string, style ...document.CellStyle) *xmlstructs.Paragraph {
	var pPr *xmlstructs.ParagraphProperties
	var rPr *xmlstructs.RunProperties
	if len(style) > 0 {
		pPr = p.mapParagraphProperties(style[0])
		rPr = p.mapRunProperties(style[0])
	}

	par := &xmlstructs.Paragraph{PPr: pPr}

	parts := p.splitTextAndFields(text)
	for _, part := range parts {
		if part.isField {
			par.Content = append(par.Content, &xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "begin"}})
			par.Content = append(par.Content, &xmlstructs.Run{InstrText: &xmlstructs.InstrText{Text: " " + part.content + " "}})
			par.Content = append(par.Content, &xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "separate"}})
			par.Content = append(par.Content, &xmlstructs.Run{RPr: rPr, T: "0"}) // Placeholder
			par.Content = append(par.Content, &xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "end"}})
		} else if part.content != "" {
			par.Content = append(par.Content, &xmlstructs.Run{RPr: rPr, T: part.content})
		}
	}

	return par
}

type textPart struct {
	content string
	isField bool
}

func (p *processor) splitTextAndFields(text string) []textPart {
	var parts []textPart
	curr := ""
	for i := 0; i < len(text); i++ {
		if text[i] == '{' {
			if curr != "" {
				parts = append(parts, textPart{content: curr, isField: false})
				curr = ""
			}
			j := i + 1
			for j < len(text) && text[j] != '}' {
				j++
			}
			if j < len(text) {
				field := text[i+1 : j]
				switch strings.TrimSpace(field) {
				case "n":
					parts = append(parts, textPart{content: "PAGE", isField: true})
				case "nb":
					parts = append(parts, textPart{content: "NUMPAGES", isField: true})
				case "date":
					parts = append(parts, textPart{content: "DATE", isField: true})
				default:
					parts = append(parts, textPart{content: field, isField: true})
				}
				i = j
			}
		} else {
			curr += string(text[i])
		}
	}
	if curr != "" {
		parts = append(parts, textPart{content: curr, isField: false})
	}
	return parts
}
