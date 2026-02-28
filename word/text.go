package word

import (
	"fmt"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func (p *processor) AddParagraph(text string, style ...document.CellStyle) error {
	if text == "" {
		// Word allows empty paragraphs
	}
	var pPr *xmlstructs.ParagraphProperties
	var rPr *xmlstructs.RunProperties
	if len(style) > 0 {
		pPr = p.mapParagraphProperties(style[0])
		rPr = p.mapRunProperties(style[0])
	}

	par := &xmlstructs.Paragraph{PPr: pPr}
	if text != "" {
		par.Content = append(par.Content, &xmlstructs.Run{
			RPr: rPr,
			T:   text,
		})
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}

	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) AddRichParagraph(spans []document.TextSpan) error {
	var pPr *xmlstructs.ParagraphProperties
	if len(spans) > 0 {
		pPr = p.mapParagraphProperties(spans[0].Style)
	}

	par := &xmlstructs.Paragraph{PPr: pPr}
	for _, span := range spans {
		if span.Text == "" {
			continue
		}
		rPr := p.mapRunProperties(span.Style)
		run := &xmlstructs.Run{
			RPr: rPr,
			T:   span.Text,
		}
		par.Content = append(par.Content, run)
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) AddHeading(text string, level int, style ...document.CellStyle) error {
	if text == "" {
		return fmt.Errorf("heading text cannot be empty")
	}
	if level < 1 || level > 6 {
		return fmt.Errorf("heading level must be between 1 and 6")
	}
	var s document.CellStyle
	if len(style) > 0 {
		s = style[0]
	} else {
		size := 14
		if level == 1 {
			size = 20
		} else if level == 2 {
			size = 18
		}
		s = document.CellStyle{Size: size, Bold: true}
	}
	return p.AddParagraph(text, s)
}

func (p *processor) AddList(items []string, ordered bool, style ...document.CellStyle) error {
	if len(items) == 0 {
		return fmt.Errorf("list items cannot be empty")
	}
	p.ensureNumbering()
	numID := 1
	if ordered {
		numID = 2
	}

	var s document.CellStyle
	if len(style) > 0 {
		s = style[0]
	}

	// Use Indent as the level index (0-8)
	ilvl := 0
	if s.Indent > 0 && s.Indent < 9 {
		ilvl = int(s.Indent)
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}

	for _, item := range items {
		pPr := &xmlstructs.ParagraphProperties{
			NumPr: &xmlstructs.NumPr{
				ILvl:  &xmlstructs.ValInt{Val: ilvl},
				NumID: &xmlstructs.ValInt{Val: numID},
			},
		}
		if s.Name != "" {
			pPr.PStyle = &xmlstructs.ParagraphStyle{Val: s.Name}
		}

		rPr := p.mapRunProperties(s)

		par := &xmlstructs.Paragraph{
			PPr: pPr,
			Content: []any{
				&xmlstructs.Run{
					RPr: rPr,
					T:   item,
				},
			},
		}
		p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	}
	return nil
}

func (p *processor) AddPageBreak() error {
	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{
				Br: &xmlstructs.Break{Type: "page"},
			},
		},
	}
	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) AddHyperlink(text, url string, style ...document.CellStyle) error {
	if text == "" || url == "" {
		return fmt.Errorf("hyperlink text and url cannot be empty")
	}
	if p.docRels == nil {
		p.docRels = &xmlstructs.Relationships{}
	}
	rID := p.docRels.AddRelationship(
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink",
		url,
	)
	p.docRels.Rels[len(p.docRels.Rels)-1].TargetMode = "External"

	var rPr *xmlstructs.RunProperties
	if len(style) > 0 {
		rPr = p.mapRunProperties(style[0])
	}
	if rPr == nil {
		rPr = &xmlstructs.RunProperties{}
	}
	rPr.Color = &xmlstructs.Color{Val: "0000FF"}
	rPr.U = &xmlstructs.Underline{Val: "single"}

	h := &xmlstructs.Hyperlink{
		ID: rID,
		Runs: []*xmlstructs.Run{
			{
				RPr: rPr,
				T:   text,
			},
		},
	}

	par := &xmlstructs.Paragraph{
		Content: []any{h},
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) AddBookmark(name string) error {
	p.bookmarkCounter++
	id := p.bookmarkCounter

	start := &xmlstructs.BookmarkStart{ID: id, Name: name}
	end := &xmlstructs.BookmarkEnd{ID: id}

	par := &xmlstructs.Paragraph{
		Content: []any{start, end},
	}

	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}

func (p *processor) AddTableOfContents() error {
	// Simple TOC implementation via complex fields
	// Requires MS Word to update on open
	par := &xmlstructs.Paragraph{
		Content: []any{
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "begin"}},
			&xmlstructs.Run{InstrText: &xmlstructs.InstrText{Text: ` TOC \o "1-3" \h \z \u `}},
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "separate"}},
			&xmlstructs.Run{T: "Updating Table of Contents..."},
			&xmlstructs.Run{FldChar: &xmlstructs.FldChar{FldCharType: "end"}},
		},
	}
	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)
	return nil
}
