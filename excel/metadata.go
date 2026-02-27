package excel

import (
	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// metadata handles document metadata operations.
type metadata struct{ *state }

func (e *metadata) GetMetadata() (document.Metadata, error) {
	if e.coreProperties == nil {
		return document.Metadata{}, nil
	}
	return document.Metadata{
		Title:       e.coreProperties.Title,
		Author:      e.coreProperties.Creator,
		Subject:     e.coreProperties.Subject,
		Description: e.coreProperties.Description,
	}, nil
}

func (e *metadata) SetMetadata(metadata document.Metadata) error {
	if e.coreProperties == nil {
		e.coreProperties = &xmlstructs.CoreProperties{}
	}
	e.coreProperties.Title = metadata.Title
	e.coreProperties.Creator = metadata.Author
	e.coreProperties.Subject = metadata.Subject
	e.coreProperties.Description = metadata.Description
	return nil
}
