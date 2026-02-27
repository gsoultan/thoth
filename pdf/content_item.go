package pdf

import "github.com/gsoultan/thoth/document"

type contentItem struct {
	isParagraph bool
	text        string
	style       document.CellStyle
	isImage     bool
	path        string
	width       float64
	height      float64
	isTable     bool
	rows        int
	cols        int
	cells       [][][]cellItem
	isPageBreak bool
	isShape     bool
	shapeType   string
	x1, y1      float64
	x2, y2      float64
}
