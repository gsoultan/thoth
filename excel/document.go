package excel

import (
	"context"
	"fmt"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// Document implements the document.Document and document.Spreadsheet interfaces.
// It uses composition to delegate logic to specialized smaller units.
type Document struct {
	*state
	lifecycle
	sheetManager
	cellProcessor
	styleManager
	metadata
	imageProcessor
	content
	pageSettings
}

// Fluent API: Sheet returns a sheet-scoped handle, auto-creating the sheet if needed.
func (d *Document) Sheet(name string) (document.Sheet, error) {
	if _, ok := d.sheets[name]; !ok {
		if err := d.addSheet(name); err != nil {
			return nil, err
		}
	}
	return &sheetHandle{state: d.state, ctx: d.ctx, name: name}, nil
}

func (d *Document) Export(uri string) error {
	if d.exportFunc == nil {
		return fmt.Errorf("export function not configured")
	}
	return d.exportFunc(d, uri)
}

func (d *Document) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *Document) SetExportFunc(fn func(doc document.Document, uri string) error) {
	d.exportFunc = fn
}

// NewDocument creates a new instance of an Excel document processor.
func NewDocument() document.Document {
	state := &state{
		sheets:    make(map[string]*xmlstructs.Worksheet),
		media:     make(map[string][]byte),
		sheetRels: make(map[string]*xmlstructs.Relationships),
		workbook: &xmlstructs.Workbook{
			XMLNS_R: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
		},
		workbookRels: &xmlstructs.Relationships{},
		rootRels: &xmlstructs.Relationships{
			Rels: []xmlstructs.Relationship{
				{
					ID:     "rId1",
					Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument",
					Target: "xl/workbook.xml",
				},
			},
		},
		contentTypes: xmlstructs.NewContentTypes(),
		wbRelsPath:   "xl/_rels/workbook.xml.rels",
		rootRelsPath: "_rels/.rels",
	}
	return &Document{
		state:          state,
		lifecycle:      lifecycle{state},
		sheetManager:   sheetManager{state},
		cellProcessor:  cellProcessor{state},
		styleManager:   styleManager{state},
		metadata:       metadata{state},
		imageProcessor: imageProcessor{state},
		content:        content{state},
		pageSettings:   pageSettings{state},
	}
}
