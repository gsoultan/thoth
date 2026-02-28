package pdf

import "github.com/gsoultan/thoth/pdf/internal/objects"

type pageInfo struct {
	contentItems []*contentItem
	sb           string
	links        []link
	fields       []*contentItem // Form fields on this page
	footnotes    []string       // Footnotes on this page
	structItems  []objects.Dictionary
	posY         float64
	w, h         int
}
