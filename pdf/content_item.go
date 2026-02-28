package pdf

import "github.com/gsoultan/thoth/document"

type contentItem struct {
	isParagraph  bool
	isHeading    bool
	level        int
	text         string
	style        document.CellStyle
	isImage      bool
	path         string
	width        float64
	height       float64
	isTable      bool
	rows         int
	cols         int
	HeaderRows   int
	colWidths    []float64
	cells        [][][]cellItem
	isPageBreak  bool
	isShape      bool
	shapeType    string
	x1, y1       float64
	x2, y2       float64
	isFormField  bool
	fieldType    string // "text", "checkbox"
	fieldName    string
	options      []string
	isImported   bool
	importPath   string
	importPage   int
	isFootnote   bool
	isTOC        bool
	isList       bool
	listItems    []string
	ordered      bool
	isBookmark   bool
	bookmarkName string
	isRich       bool
	spans        []document.TextSpan
}
