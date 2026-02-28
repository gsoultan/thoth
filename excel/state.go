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
	drawings       map[string]*xmlstructs.WsDr
	tables         map[string]*xmlstructs.Table
	// Optimization caches
	sharedStringsIndex map[string]int
	fontsIndex         map[string]int
	fillsIndex         map[string]int
	bordersIndex       map[string]int
	xfsIndex           map[string]int
	cellCache          map[string]map[string]*xmlstructs.Cell
}

func (e *state) processor() *processor {
	return &processor{
		state:          e,
		sheetProcessor: sheetProcessor{e},
		cellProcessor:  cellProcessor{e},
		styleProcessor: styleProcessor{e},
		mediaProcessor: mediaProcessor{e},
	}
}
