package excel

import (
	"archive/zip"
	"context"
	"os"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// state holds the shared internal state for the Excel document.
type state struct {
	ctx            context.Context
	exportFunc     func(doc document.Document, uri string) error
	reader         *zip.ReadCloser
	tempFile       *os.File
	workbook       *xmlstructs.Workbook
	sharedStrings  *xmlstructs.SharedStrings
	sheets         map[string]*xmlstructs.Worksheet
	coreProperties *xmlstructs.CoreProperties
	workbookRels   *xmlstructs.Relationships
	styles         *xmlstructs.Styles
	wbRelsPath     string
	rootRelsPath   string
	contentTypes   *xmlstructs.ContentTypes
	rootRels       *xmlstructs.Relationships
	media          map[string][]byte
	sheetRels      map[string]*xmlstructs.Relationships
}
