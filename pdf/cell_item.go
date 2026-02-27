package pdf

import "github.com/gsoultan/thoth/document"

type cellItem struct {
	text    string
	style   document.CellStyle
	isImage bool
	path    string
	width   float64
	height  float64
	rowSpan int
	colSpan int
	hidden  bool
}
