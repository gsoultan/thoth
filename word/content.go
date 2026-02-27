package word

import (
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// content handles reading and searching document content.
type content struct{ *state }

// ReadContent returns the text content of the document.
func (w *content) ReadContent() (string, error) {
	buf := document.GetBuffer()
	defer document.PutBuffer(buf)

	if w.xmlDoc != nil {
		for _, c := range w.xmlDoc.Body.Content {
			if p, ok := c.(xmlstructs.Paragraph); ok {
				for _, r := range p.Runs {
					if r.T != "" {
						buf.WriteString(r.T)
					}
				}
				buf.WriteString("\n")
			}
		}
	}
	return strings.TrimSpace(buf.String()), nil
}

// Search finds keywords in the document.
func (w *content) Search(keywords []string) ([]document.SearchResult, error) {
	content, err := w.ReadContent()
	if err != nil {
		return nil, err
	}

	strategy := document.NewSimpleSearchStrategy()
	return strategy.Execute(w.ctx, content, keywords)
}

// Replace replaces keywords with new values.
func (w *content) Replace(replacements map[string]string) error {
	if w.xmlDoc == nil {
		return nil
	}

	for i, c := range w.xmlDoc.Body.Content {
		if p, ok := c.(xmlstructs.Paragraph); ok {
			for j := range p.Runs {
				for old, new := range replacements {
					p.Runs[j].T = strings.ReplaceAll(p.Runs[j].T, old, new)
				}
			}
			w.xmlDoc.Body.Content[i] = p
		}
	}
	return nil
}
