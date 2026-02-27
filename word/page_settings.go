package word

import (
	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// pageSettings handles page layout settings.
type pageSettings struct{ *state }

// SetPageSettings configures the document's layout.
func (w *pageSettings) SetPageSettings(settings document.PageSettings) error {
	if w.xmlDoc == nil {
		w.xmlDoc = &xmlstructs.Document{}
	}
	if w.xmlDoc.Body.SectPr == nil {
		w.xmlDoc.Body.SectPr = &xmlstructs.SectPr{}
	}

	(&processor{w.state}).applyPageSettingsToSect(w.xmlDoc.Body.SectPr, settings)
	return nil
}
