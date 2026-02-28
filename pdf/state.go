package pdf

import (
	"cmp"
	"context"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

// state holds the shared internal state for the PDF document.
type bookmark struct {
	title string
	level int
	page  int
	posY  float64
}

type state struct {
	ctx          context.Context
	exportFunc   func(doc document.Document, uri string) error
	objects      []objects.Object
	root         objects.Dictionary
	info         objects.Dictionary
	meta         document.Metadata
	pageSettings document.PageSettings
	contentItems []*contentItem
	header       []*contentItem
	footer       []*contentItem
	watermark    *contentItem
	bookmarks    []bookmark
	fonts        map[string]string // Name -> path
	password     string
	attachments  []attachment
}

type attachment struct {
	path        string
	name        string
	description string
}

func (s *state) getMargins() document.Margins {
	m := s.pageSettings.Margins
	m.Top = cmp.Or(m.Top, 50.0)
	m.Bottom = cmp.Or(m.Bottom, 50.0)
	m.Left = cmp.Or(m.Left, 50.0)
	m.Right = cmp.Or(m.Right, 50.0)
	return m
}
