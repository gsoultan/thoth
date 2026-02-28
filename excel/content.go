package excel

import (
	"strings"

	"github.com/gsoultan/thoth/document"
)

// content handles reading and searching document content.
type content struct{ *state }

// ReadContent returns the text content of the document.
func (e *content) ReadContent() (string, error) {
	buf := document.GetBuffer()
	defer document.PutBuffer(buf)

	sheets, _ := e.GetSheets()
	for _, sheet := range sheets {
		ws := e.sheets[sheet]
		for _, row := range ws.SheetData.Rows {
			for _, cell := range row.Cells {
				val := e.resolveValue(cell)
				if val != "" {
					buf.WriteString(val)
					buf.WriteString(" ")
				}
			}
		}
	}
	return strings.TrimSpace(buf.String()), nil
}

// Search finds keywords in the document.
func (e *content) Search(keywords []string) ([]document.SearchResult, error) {
	content, err := e.ReadContent()
	if err != nil {
		return nil, err
	}

	strategy := document.NewSimpleSearchStrategy()
	return strategy.Execute(e.ctx, content, keywords)
}

// Replace replaces keywords with new values.
func (e *content) Replace(replacements map[string]string) error {
	if e.sharedStrings == nil {
		return nil
	}

	for i := range e.sharedStrings.SI {
		for old, new := range replacements {
			e.sharedStrings.SI[i].T = strings.ReplaceAll(e.sharedStrings.SI[i].T, old, new)
		}
	}

	return nil
}
