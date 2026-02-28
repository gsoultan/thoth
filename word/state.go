package word

import (
	"archive/zip"
	"context"
	"os"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// state holds the shared internal state for the Word document.
type state struct {
	ctx             context.Context
	exportFunc      func(doc document.Document, uri string) error
	reader          *zip.ReadCloser
	tempFile        *os.File
	doc             *xmlstructs.Document
	xmlDoc          *xmlstructs.Document
	coreProperties  *xmlstructs.CoreProperties
	appProperties   *xmlstructs.AppProperties
	docRels         *xmlstructs.Relationships
	rootRels        *xmlstructs.Relationships
	contentTypes    *xmlstructs.ContentTypes
	headers         map[string]*xmlstructs.Header
	footers         map[string]*xmlstructs.Footer
	headerRels      map[string]*xmlstructs.Relationships
	footerRels      map[string]*xmlstructs.Relationships
	numbering       *xmlstructs.Numbering
	bookmarkCounter int
	footnotes       *xmlstructs.Footnotes
	footnoteCounter int
	styles          *xmlstructs.Styles
	settings        *xmlstructs.Settings
	media           map[string][]byte
}
