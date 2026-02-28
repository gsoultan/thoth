package pdf

import (
	"strings"

	"github.com/gsoultan/thoth/pdf/internal/objects"
)

type renderingContext struct {
	mgr              *objects.ObjectManager
	imageRefs        map[string]objects.Reference
	imageNames       map[string]string
	fontRefs         map[string]objects.Reference
	fontNames        map[string]string
	importRefs       map[string]objects.Reference
	pageRefs         []objects.Reference
	pages            []pageInfo
	currentSb        strings.Builder
	currentLinks     []link
	currentFields    []*contentItem
	currentFootnotes []string
	currentStructs   []objects.Dictionary
	bookmarks        []bookmark
	currentColumn    int
	mcidCounter      int
	posY             float64
	w, h             int
	allFields        []objects.Reference
	fileID           []byte
	customWidths     map[string]map[rune]uint16
	unitsPerEm       map[string]uint16
	smaskRefs        map[string]objects.Reference
	extGStates       map[float64]objects.Reference
}
