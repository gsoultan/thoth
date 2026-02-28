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
	processor
	metadata
	content
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

func (d *Document) SetPassword(password string) error {
	if d.workbook == nil {
		return fmt.Errorf("workbook not initialized")
	}
	d.workbook.WorkbookProtection = &xmlstructs.WorkbookProtection{
		WorkbookPassword: excelPasswordHash(password),
		LockStructure:    1,
	}
	return nil
}

func (d *Document) SetNamedRange(name, ref string) error {
	return d.setNamedRange(name, ref)
}

// NewDocument creates a new instance of an Excel document processor.
func NewDocument() document.Document {
	state := &state{
		sheets:    make(map[string]*xmlstructs.Worksheet),
		media:     make(map[string][]byte),
		sheetRels: make(map[string]*xmlstructs.Relationships),
		drawings:  make(map[string]*xmlstructs.WsDr),
		tables:    make(map[string]*xmlstructs.Table),
		workbook: &xmlstructs.Workbook{
			XMLNS_R: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
			WorkbookPr: &xmlstructs.WorkbookPr{
				Date1904: 0,
			},
			CalcPr: &xmlstructs.CalcPr{
				FullCalcOnLoad: 1,
			},
			WorkbookViews: &xmlstructs.WorkbookViews{
				Items: []xmlstructs.WorkbookView{
					{ActiveTab: 0},
				},
			},
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
		contentTypes:       xmlstructs.NewContentTypes(),
		styles:             xmlstructs.NewDefaultStyles(),
		wbRelsPath:         "xl/_rels/workbook.xml.rels",
		rootRelsPath:       "_rels/.rels",
		sharedStringsIndex: make(map[string]int),
		fontsIndex:         make(map[string]int),
		fillsIndex:         make(map[string]int),
		bordersIndex:       make(map[string]int),
		xfsIndex:           make(map[string]int),
		cellCache:          make(map[string]map[string]*xmlstructs.Cell),
	}
	return &Document{
		state:     state,
		lifecycle: lifecycle{state},
		processor: processor{
			state:          state,
			sheetProcessor: sheetProcessor{state},
			cellProcessor:  cellProcessor{state},
			styleProcessor: styleProcessor{state},
			mediaProcessor: mediaProcessor{state},
		},
		metadata: metadata{state},
		content:  content{state},
	}
}
