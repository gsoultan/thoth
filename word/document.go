package word

import (
	"context"
	"fmt"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// Document implements the document.Document and document.WordProcessor interfaces.
// It uses composition to delegate logic to specialized smaller units.
type Document struct {
	*state
	lifecycle
	processor
	metadata
	content
	pageSettings
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

// NewDocument creates a new instance of a Word document processor.
func NewDocument() document.Document {
	state := &state{
		media: make(map[string][]byte),
		doc: &xmlstructs.Document{
			W:   "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
			R:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
			WP:  "http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing",
			A:   "http://schemas.openxmlformats.org/drawingml/2006/main",
			Pic: "http://schemas.openxmlformats.org/drawingml/2006/picture",
			Body: xmlstructs.Body{
				Content: make([]any, 0),
				SectPr: &xmlstructs.SectPr{
					PgSz: &xmlstructs.PgSz{
						W:      11906,
						H:      16838,
						Orient: "portrait",
					},
					PgMar: &xmlstructs.PgMar{
						Top:    1440,
						Bottom: 1440,
						Left:   1440,
						Right:  1440,
					},
				},
			},
		},
		coreProperties: &xmlstructs.CoreProperties{
			Creator: "Thoth Go Library",
		},
		appProperties: xmlstructs.NewAppProperties(),
		docRels:       &xmlstructs.Relationships{},
		rootRels: &xmlstructs.Relationships{
			Rels: []xmlstructs.Relationship{
				{
					ID:     "rId1",
					Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument",
					Target: "word/document.xml",
				},
				{
					ID:     "rId2",
					Type:   "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties",
					Target: "docProps/core.xml",
				},
				{
					ID:     "rId3",
					Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties",
					Target: "docProps/app.xml",
				},
			},
		},
		contentTypes: &xmlstructs.ContentTypes{
			Defaults: []xmlstructs.Default{
				{Extension: "rels", ContentType: "application/vnd.openxmlformats-package.relationships+xml"},
				{Extension: "xml", ContentType: "application/xml"},
			},
			Override: []xmlstructs.Override{
				{PartName: "/word/document.xml", ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"},
				{PartName: "/docProps/core.xml", ContentType: "application/vnd.openxmlformats-package.core-properties+xml"},
				{PartName: "/docProps/app.xml", ContentType: "application/vnd.openxmlformats-officedocument.extended-properties+xml"},
			},
		},
	}
	state.xmlDoc = state.doc
	return &Document{
		state:        state,
		lifecycle:    lifecycle{state},
		processor:    processor{state},
		metadata:     metadata{state},
		content:      content{state},
		pageSettings: pageSettings{state},
	}
}
