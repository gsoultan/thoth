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
	ctx            context.Context
	exportFunc     func(doc document.Document, uri string) error
	reader         *zip.ReadCloser
	tempFile       *os.File
	doc            *xmlstructs.Document
	xmlDoc         *xmlstructs.Document
	coreProperties *xmlstructs.CoreProperties
	appProperties  *xmlstructs.AppProperties
	docRels        *xmlstructs.Relationships
	rootRels       *xmlstructs.Relationships
	contentTypes   *xmlstructs.ContentTypes
	media          map[string][]byte
}
