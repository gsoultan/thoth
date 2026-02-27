package core

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel"
	"github.com/gsoultan/thoth/pdf"
	"github.com/gsoultan/thoth/word"
)

// DocumentFactory is responsible for creating document processors based on file extensions.
// It supports .xlsx, .xls (Excel), .docx, .doc (Word), and .pdf (PDF).
type DocumentFactory struct{}

// NewDocumentFactory creates a new DocumentFactory instance.
func NewDocumentFactory() *DocumentFactory {
	return &DocumentFactory{}
}

// Create returns a new Document processor for the given filename.
// The processor type is determined by the file extension.
func (f *DocumentFactory) Create(filename string) (document.Document, error) {
	ext := strings.ToLower(getExt(filename))

	switch ext {
	case ".xlsx", ".xls":
		return excel.NewDocument(), nil
	case ".docx", ".doc":
		return word.NewDocument(), nil
	case ".pdf":
		return pdf.NewDocument(), nil
	default:
		return nil, fmt.Errorf("%w: %s", document.ErrUnsupportedFormat, ext)
	}
}

func getExt(filename string) string {
	u, err := url.Parse(filename)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "s3") {
		return path.Ext(u.Path)
	}

	if idx := strings.LastIndex(filename, "."); idx != -1 {
		return filename[idx:]
	}
	return ""
}
