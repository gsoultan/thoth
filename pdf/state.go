package pdf

import (
	"context"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

// state holds the shared internal state for the PDF document.
type state struct {
	ctx          context.Context
	exportFunc   func(doc document.Document, uri string) error
	objects      []objects.Object
	root         objects.Dictionary
	info         objects.Dictionary
	pageSettings document.PageSettings
	contentItems []contentItem
	header       []contentItem
	footer       []contentItem
}
