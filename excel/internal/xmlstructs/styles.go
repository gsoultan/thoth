package xmlstructs

import "encoding/xml"

// Styles defines the structure of xl/styles.xml
type Styles struct {
	XMLName      xml.Name      `xml:"http://schemas.openxmlformats.org/spreadsheetml/2006/main styleSheet"`
	NumFmts      *NumFmts      `xml:"numFmts,omitempty"`
	Fonts        Fonts         `xml:"fonts"`
	Fills        Fills         `xml:"fills"`
	Borders      Borders       `xml:"borders"`
	CellStyleXfs *CellStyleXfs `xml:"cellStyleXfs,omitempty"`
	CellXfs      CellXfs       `xml:"cellXfs"`
	Dxfs         *Dxfs         `xml:"dxfs,omitempty"`
}

type Dxfs struct {
	Count int  `xml:"count,attr"`
	Items []Xf `xml:"xf"`
}

type CellStyleXfs struct {
	Count int  `xml:"count,attr"`
	Items []Xf `xml:"xf"`
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
	s.Fonts.Items = append(s.Fonts.Items, font)
	s.Fonts.Count = len(s.Fonts.Items)
	return len(s.Fonts.Items) - 1
}

// AddXf adds a cell format and returns its index
func (s *Styles) AddXf(xf Xf) int {
	if xf.XfID == nil {
		zero := 0
		xf.XfID = &zero
	}
	s.CellXfs.Items = append(s.CellXfs.Items, xf)
	s.CellXfs.Count = len(s.CellXfs.Items)
	return len(s.CellXfs.Items) - 1
}

func (s *Styles) AddNumFmt(formatCode string) int {
	if s.NumFmts == nil {
		s.NumFmts = &NumFmts{Count: 0, Items: make([]NumFmt, 0)}
	}
	// Check for existing custom number formats
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
	s.Fills.Items = append(s.Fills.Items, fill)
	s.Fills.Count = len(s.Fills.Items)
	return len(s.Fills.Items) - 1
}

func (s *Styles) AddBorder(border Border) int {
	s.Borders.Items = append(s.Borders.Items, border)
	s.Borders.Count = len(s.Borders.Items)
	return len(s.Borders.Items) - 1
}

func (s *Styles) AddDxf(xf Xf) int {
	if s.Dxfs == nil {
		s.Dxfs = &Dxfs{Items: make([]Xf, 0)}
	}
	s.Dxfs.Items = append(s.Dxfs.Items, xf)
	s.Dxfs.Count = len(s.Dxfs.Items)
	return len(s.Dxfs.Items) - 1
}

// NewDefaultStyles creates a new instance of Styles with mandatory defaults.
func NewDefaultStyles() *Styles {
	zero := 0
	return &Styles{
		Fonts: Fonts{
			Count: 1,
			Items: []Font{
				{
					Size: &ValInt{Val: 11},
					Name: &ValString{Val: "Calibri"},
				},
			},
		},
		Fills: Fills{
			Count: 2,
			Items: []Fill{
				{PatternFill: &PatternFill{PatternType: "none"}},
				{PatternFill: &PatternFill{PatternType: "gray125"}},
			},
		},
		Borders: Borders{
			Count: 1,
			Items: []Border{{}},
		},
		CellStyleXfs: &CellStyleXfs{
			Count: 1,
			Items: []Xf{
				{
					NumFmtID: 0,
					FontID:   0,
					FillID:   0,
					BorderID: 0,
				},
			},
		},
		CellXfs: CellXfs{
			Count: 1,
			Items: []Xf{
				{
					NumFmtID: 0,
					FontID:   0,
					FillID:   0,
					BorderID: 0,
					XfID:     &zero,
				},
			},
		},
	}
}
