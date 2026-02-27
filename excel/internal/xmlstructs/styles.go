package xmlstructs

import "encoding/xml"

// Styles defines the structure of xl/styles.xml
type Styles struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/spreadsheetml/2006/main styleSheet"`
	NumFmts *NumFmts `xml:"numFmts,omitempty"`
	Fonts   Fonts    `xml:"fonts"`
	Fills   Fills    `xml:"fills"`
	Borders Borders  `xml:"borders"`
	CellXfs CellXfs  `xml:"cellXfs"`
}

type NumFmts struct {
	Count int      `xml:"count,attr"`
	Items []NumFmt `xml:"numFmt"`
}

type NumFmt struct {
	NumFmtID   int    `xml:"numFmtId,attr"`
	FormatCode string `xml:"formatCode,attr"`
}

// AddFont adds a font and returns its index
func (s *Styles) AddFont(font Font) int {
	for i, f := range s.Fonts.Items {
		if f.equals(font) {
			return i
		}
	}
	s.Fonts.Items = append(s.Fonts.Items, font)
	s.Fonts.Count++
	return len(s.Fonts.Items) - 1
}

// AddXf adds a cell format and returns its index
func (s *Styles) AddXf(xf Xf) int {
	for i, x := range s.CellXfs.Items {
		if x.equals(xf) {
			return i
		}
	}
	s.CellXfs.Items = append(s.CellXfs.Items, xf)
	s.CellXfs.Count++
	return len(s.CellXfs.Items) - 1
}

func (s *Styles) AddNumFmt(formatCode string) int {
	if s.NumFmts == nil {
		s.NumFmts = &NumFmts{Count: 0, Items: make([]NumFmt, 0)}
	}
	for _, nf := range s.NumFmts.Items {
		if nf.FormatCode == formatCode {
			return nf.NumFmtID
		}
	}
	// Custom IDs start from 164
	newID := 164
	for _, nf := range s.NumFmts.Items {
		if nf.NumFmtID >= newID {
			newID = nf.NumFmtID + 1
		}
	}
	s.NumFmts.Items = append(s.NumFmts.Items, NumFmt{NumFmtID: newID, FormatCode: formatCode})
	s.NumFmts.Count++
	return newID
}

func (s *Styles) AddFill(fill Fill) int {
	for i, f := range s.Fills.Items {
		if f.equals(fill) {
			return i
		}
	}
	s.Fills.Items = append(s.Fills.Items, fill)
	s.Fills.Count++
	return len(s.Fills.Items) - 1
}

func (s *Styles) AddBorder(border Border) int {
	for i, b := range s.Borders.Items {
		if b.equals(border) {
			return i
		}
	}
	s.Borders.Items = append(s.Borders.Items, border)
	s.Borders.Count++
	return len(s.Borders.Items) - 1
}

func (f Font) equals(other Font) bool {
	if (f.Bold == nil) != (other.Bold == nil) {
		return false
	}
	if (f.Italic == nil) != (other.Italic == nil) {
		return false
	}
	if (f.Size == nil) != (other.Size == nil) {
		return false
	}
	if f.Size != nil && f.Size.Val != other.Size.Val {
		return false
	}
	if (f.Color == nil) != (other.Color == nil) {
		return false
	}
	if f.Color != nil && f.Color.RGB != other.Color.RGB {
		return false
	}
	return true
}

func (f Fill) equals(other Fill) bool {
	if (f.PatternFill == nil) != (other.PatternFill == nil) {
		return false
	}
	if f.PatternFill != nil {
		if f.PatternFill.PatternType != other.PatternFill.PatternType {
			return false
		}
		if (f.PatternFill.FgColor == nil) != (other.PatternFill.FgColor == nil) {
			return false
		}
		if f.PatternFill.FgColor != nil && f.PatternFill.FgColor.RGB != other.PatternFill.FgColor.RGB {
			return false
		}
	}
	return true
}

func (b Border) equals(other Border) bool {
	checkEdge := func(e1, e2 *BorderEdge) bool {
		if (e1 == nil) != (e2 == nil) {
			return false
		}
		if e1 != nil && e1.Style != e2.Style {
			return false
		}
		return true
	}
	return checkEdge(b.Left, other.Left) &&
		checkEdge(b.Right, other.Right) &&
		checkEdge(b.Top, other.Top) &&
		checkEdge(b.Bottom, other.Bottom)
}

func (x Xf) equals(other Xf) bool {
	if x.NumFmtID != other.NumFmtID ||
		x.FontID != other.FontID ||
		x.FillID != other.FillID ||
		x.BorderID != other.BorderID ||
		x.ApplyNumberFormat != other.ApplyNumberFormat ||
		x.ApplyFont != other.ApplyFont ||
		x.ApplyFill != other.ApplyFill ||
		x.ApplyBorder != other.ApplyBorder ||
		x.ApplyAlignment != other.ApplyAlignment {
		return false
	}
	if (x.Alignment == nil) != (other.Alignment == nil) {
		return false
	}
	if x.Alignment != nil {
		if x.Alignment.Horizontal != other.Alignment.Horizontal ||
			x.Alignment.Vertical != other.Alignment.Vertical {
			return false
		}
	}
	return true
}
