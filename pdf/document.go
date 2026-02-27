package pdf

import (
	"context"
	"fmt"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

// Document implements the document.Document and document.WordProcessor interfaces.
// It uses composition to delegate logic to specialized smaller units.
type Document struct {
	*state
	lifecycle
	content
	metadata
	processor
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

// NewDocument creates a new instance of a PDF document processor.
func NewDocument() document.Document {
	state := &state{
		objects:      make([]objects.Object, 0),
		contentItems: make([]contentItem, 0),
	}
	return &Document{
		state:     state,
		lifecycle: lifecycle{state},
		content:   content{state},
		metadata:  metadata{state},
		processor: processor{state},
	}
}
