package word

import (
	"fmt"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func (p *processor) SetHeader(text string, style ...document.CellStyle) error {
	headerType := "default"
	if len(style) > 0 && style[0].Name != "" {
		switch strings.ToLower(style[0].Name) {
		case "first", "even", "default":
			headerType = strings.ToLower(style[0].Name)
		}
	}

	if p.headers == nil {
		p.headers = make(map[string]*xmlstructs.Header)
	}
	if p.headerRels == nil {
		p.headerRels = make(map[string]*xmlstructs.Relationships)
	}

	header := &xmlstructs.Header{
		W: "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
		R: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
	}

	if text != "" {
		par := p.createParagraphWithFields(text, style...)
		header.Content = append(header.Content, par)
	}

	headerID := fmt.Sprintf("header%d.xml", len(p.headers)+1)
	p.headers[headerID] = header

	sect := p.ensureSectPr()
	if headerType == "first" {
		sect.TitlePg = &xmlstructs.TitlePg{Val: "1"}
	} else if headerType == "even" {
		if p.settings == nil {
			p.settings = xmlstructs.NewSettings()
		}
		p.settings.EvenAndOddHeaders = &xmlstructs.EvenAndOddHeaders{Val: "1"}
	}

	if p.docRels == nil {
		p.docRels = &xmlstructs.Relationships{}
	}
	rID := p.docRels.AddRelationship("http://schemas.openxmlformats.org/officeDocument/2006/relationships/header", headerID)

	// Replace existing header of same type or add new
	found := false
	for i, ref := range sect.HeaderRefs {
		if ref.Type == headerType {
			sect.HeaderRefs[i].ID = rID
			found = true
			break
		}
	}
	if !found {
		sect.HeaderRefs = append(sect.HeaderRefs, xmlstructs.HeaderReference{
			Type: headerType,
			ID:   rID,
		})
	}

	return nil
}

func (p *processor) SetFooter(text string, style ...document.CellStyle) error {
	footerType := "default"
	if len(style) > 0 && style[0].Name != "" {
		switch strings.ToLower(style[0].Name) {
		case "first", "even", "default":
			footerType = strings.ToLower(style[0].Name)
		}
	}

	if p.footers == nil {
		p.footers = make(map[string]*xmlstructs.Footer)
	}
	if p.footerRels == nil {
		p.footerRels = make(map[string]*xmlstructs.Relationships)
	}

	footer := &xmlstructs.Footer{
		W: "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
		R: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
	}

	if text != "" {
		par := p.createParagraphWithFields(text, style...)
		footer.Content = append(footer.Content, par)
	}

	footerID := fmt.Sprintf("footer%d.xml", len(p.footers)+1)
	p.footers[footerID] = footer

	sect := p.ensureSectPr()
	if footerType == "first" {
		sect.TitlePg = &xmlstructs.TitlePg{Val: "1"}
	} else if footerType == "even" {
		if p.settings == nil {
			p.settings = xmlstructs.NewSettings()
		}
		p.settings.EvenAndOddHeaders = &xmlstructs.EvenAndOddHeaders{Val: "1"}
	}

	if p.docRels == nil {
		p.docRels = &xmlstructs.Relationships{}
	}
	rID := p.docRels.AddRelationship("http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer", footerID)

	// Replace existing footer of same type or add new
	found := false
	for i, ref := range sect.FooterRefs {
		if ref.Type == footerType {
			sect.FooterRefs[i].ID = rID
			found = true
			break
		}
	}
	if !found {
		sect.FooterRefs = append(sect.FooterRefs, xmlstructs.FooterReference{
			Type: footerType,
			ID:   rID,
		})
	}

	return nil
}

func (p *processor) AddSection(settings document.PageSettings) error {
	sectPr := p.ensureSectPr()

	par := xmlstructs.Paragraph{
		PPr: &xmlstructs.ParagraphProperties{
			SectPr: sectPr,
		},
	}
	p.xmlDoc.Body.Content = append(p.xmlDoc.Body.Content, par)

	newSectPr := &xmlstructs.SectPr{}
	p.applyPageSettingsToSect(newSectPr, settings)
	p.xmlDoc.Body.SectPr = newSectPr

	return nil
}

func (p *processor) ensureSectPr() *xmlstructs.SectPr {
	if p.xmlDoc == nil {
		p.xmlDoc = p.doc
	}
	if p.xmlDoc.Body.SectPr == nil {
		p.xmlDoc.Body.SectPr = &xmlstructs.SectPr{}
	}
	return p.xmlDoc.Body.SectPr
}

func (p *processor) applyPageSettingsToSect(sect *xmlstructs.SectPr, settings document.PageSettings) {
	width, height := 11906, 16838 // A4
	if settings.PaperType == document.PaperLetter {
		width, height = 12240, 15840
	}
	orient := "portrait"
	if settings.Orientation == document.OrientationLandscape {
		orient = "landscape"
		width, height = height, width
	}

	sect.PgSz = &xmlstructs.PgSz{
		W:      width,
		H:      height,
		Orient: orient,
	}

	sect.PgMar = &xmlstructs.PgMar{
		Top:    int(settings.Margins.Top * 20),
		Bottom: int(settings.Margins.Bottom * 20),
		Left:   int(settings.Margins.Left * 20),
		Right:  int(settings.Margins.Right * 20),
	}
}
