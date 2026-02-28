package word

import (
	"bytes"
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
			w.readElementContent(c, buf)
		}
	}
	return strings.TrimSpace(buf.String()), nil
}

func (w *content) readElementContent(c any, buf *bytes.Buffer) {
	switch v := c.(type) {
	case *xmlstructs.Paragraph:
		for _, item := range v.Content {
			w.readRunContent(item, buf)
		}
		buf.WriteString("\n")
	case *xmlstructs.Table:
		for _, row := range v.Rows {
			for i, cell := range row.Cells {
				if i > 0 {
					buf.WriteString("\t")
				}
				for _, cellContent := range cell.Content {
					w.readElementContent(cellContent, buf)
				}
			}
			buf.WriteString("\n")
		}
	}
}

func (w *content) readRunContent(item any, buf *bytes.Buffer) {
	switch v := item.(type) {
	case *xmlstructs.Run:
		if v.T != "" {
			buf.WriteString(v.T)
		}
	case xmlstructs.Run:
		if v.T != "" {
			buf.WriteString(v.T)
		}
	case *xmlstructs.Hyperlink:
		for _, r := range v.Runs {
			if r.T != "" {
				buf.WriteString(r.T)
			}
		}
	case xmlstructs.Hyperlink:
		for _, r := range v.Runs {
			if r.T != "" {
				buf.WriteString(r.T)
			}
		}
	}
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

	for _, c := range w.xmlDoc.Body.Content {
		if p, ok := c.(*xmlstructs.Paragraph); ok {
			for j, item := range p.Content {
				switch v := item.(type) {
				case *xmlstructs.Run:
					for old, new := range replacements {
						v.T = strings.ReplaceAll(v.T, old, new)
					}
				case xmlstructs.Run:
					for old, new := range replacements {
						v.T = strings.ReplaceAll(v.T, old, new)
					}
					p.Content[j] = v
				case *xmlstructs.Hyperlink:
					for k := range v.Runs {
						for old, new := range replacements {
							v.Runs[k].T = strings.ReplaceAll(v.Runs[k].T, old, new)
						}
					}
				case xmlstructs.Hyperlink:
					for k := range v.Runs {
						for old, new := range replacements {
							v.Runs[k].T = strings.ReplaceAll(v.Runs[k].T, old, new)
						}
					}
					p.Content[j] = v
				}
			}
		}
	}
	return nil
}
