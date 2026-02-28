package word

import (
	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// metadata handles document metadata operations.
type metadata struct{ *state }

func (w *metadata) GetMetadata() (document.Metadata, error) {
	if w.coreProperties == nil {
		return document.Metadata{}, nil
	}
	return document.Metadata{
		Title:       w.coreProperties.Title,
		Author:      w.coreProperties.Creator,
		Subject:     w.coreProperties.Subject,
		Description: w.coreProperties.Description,
	}, nil
}

func (w *metadata) SetMetadata(metadata document.Metadata) error {
	if w.coreProperties == nil {
		w.coreProperties = &xmlstructs.CoreProperties{}
	}
	w.coreProperties.Title = metadata.Title
	w.coreProperties.Creator = metadata.Author
	w.coreProperties.Subject = metadata.Subject
	w.coreProperties.Description = metadata.Description
	if metadata.Company != "" {
		w.appProperties.Company = metadata.Company
	}
	return nil
}
