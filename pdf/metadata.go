package pdf

import (
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

// metadata handles document metadata operations.
type metadata struct{ *state }

// Metadata methods

func (p *metadata) GetMetadata() (document.Metadata, error) {
	if p.info == nil {
		return document.Metadata{}, nil
	}
	meta := document.Metadata{}
	if title, ok := p.info["Title"].(objects.PDFString); ok {
		meta.Title = string(title)
	}
	if author, ok := p.info["Author"].(objects.PDFString); ok {
		meta.Author = string(author)
	}
	if subject, ok := p.info["Subject"].(objects.PDFString); ok {
		meta.Subject = string(subject)
	}
	if keywords, ok := p.info["Keywords"].(objects.PDFString); ok {
		meta.Keywords = strings.Split(string(keywords), ",")
	}
	if desc, ok := p.info["Description"].(objects.PDFString); ok {
		meta.Description = string(desc)
	}
	return meta, nil
}

func (p *metadata) SetMetadata(metadata document.Metadata) error {
	if p.info == nil {
		p.info = make(objects.Dictionary)
	}
	p.info["Title"] = objects.PDFString(metadata.Title)
	p.info["Author"] = objects.PDFString(metadata.Author)
	p.info["Subject"] = objects.PDFString(metadata.Subject)
	p.info["Keywords"] = objects.PDFString(strings.Join(metadata.Keywords, ","))
	p.info["Description"] = objects.PDFString(metadata.Description)
	return nil
}
