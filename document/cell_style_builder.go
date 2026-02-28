package document

// CellStyleBuilder is a fluent builder for CellStyle.
type CellStyleBuilder interface {
	Bold() CellStyleBuilder
	Italic() CellStyleBuilder
	Size(size int) CellStyleBuilder
	Color(hex string) CellStyleBuilder
	Background(hex string) CellStyleBuilder
	Border() CellStyleBuilder
	BorderTop() CellStyleBuilder
	BorderBottom() CellStyleBuilder
	BorderLeft() CellStyleBuilder
	BorderRight() CellStyleBuilder
	BorderWidth(width float64) CellStyleBuilder
	BorderColor(hex string) CellStyleBuilder
	Align(h, v string) CellStyleBuilder
	NumberFormat(format string) CellStyleBuilder
	Name(name string) CellStyleBuilder
	Link(url string) CellStyleBuilder
	Font(name string) CellStyleBuilder
	Alt(text string) CellStyleBuilder
	LineSpacing(spacing float64) CellStyleBuilder
	SpacingBefore(spacing float64) CellStyleBuilder
	SpacingAfter(spacing float64) CellStyleBuilder
	Indent(points float64) CellStyleBuilder
	Hanging(points float64) CellStyleBuilder
	Padding(points float64) CellStyleBuilder
	KeepWithNext() CellStyleBuilder
	KeepTogether() CellStyleBuilder
	Superscript() CellStyleBuilder
	Subscript() CellStyleBuilder
	Absolute() CellStyleBuilder
	Pos(x, y float64) CellStyleBuilder
	Opacity(opacity float64) CellStyleBuilder
	DashPattern(pattern []float64) CellStyleBuilder
	Build() CellStyle
}

type cellStyleBuilder struct {
	style CellStyle
}

func NewCellStyleBuilder() CellStyleBuilder {
	return &cellStyleBuilder{}
}

func (b *cellStyleBuilder) Bold() CellStyleBuilder {
	b.style.Bold = true
	return b
}

func (b *cellStyleBuilder) Italic() CellStyleBuilder {
	b.style.Italic = true
	return b
}

func (b *cellStyleBuilder) Size(size int) CellStyleBuilder {
	b.style.Size = size
	return b
}

func (b *cellStyleBuilder) Color(hex string) CellStyleBuilder {
	b.style.Color = hex
	return b
}

func (b *cellStyleBuilder) Background(hex string) CellStyleBuilder {
	b.style.Background = hex
	return b
}

func (b *cellStyleBuilder) Border() CellStyleBuilder {
	b.style.Border = true
	b.style.BorderTop = true
	b.style.BorderBottom = true
	b.style.BorderLeft = true
	b.style.BorderRight = true
	return b
}

func (b *cellStyleBuilder) BorderTop() CellStyleBuilder {
	b.style.BorderTop = true
	return b
}

func (b *cellStyleBuilder) BorderBottom() CellStyleBuilder {
	b.style.BorderBottom = true
	return b
}

func (b *cellStyleBuilder) BorderLeft() CellStyleBuilder {
	b.style.BorderLeft = true
	return b
}

func (b *cellStyleBuilder) BorderRight() CellStyleBuilder {
	b.style.BorderRight = true
	return b
}

func (b *cellStyleBuilder) BorderWidth(width float64) CellStyleBuilder {
	b.style.BorderWidth = width
	return b
}

func (b *cellStyleBuilder) BorderColor(hex string) CellStyleBuilder {
	b.style.BorderColor = hex
	return b
}

func (b *cellStyleBuilder) Align(h, v string) CellStyleBuilder {
	b.style.Horizontal = h
	b.style.Vertical = v
	return b
}

func (b *cellStyleBuilder) NumberFormat(format string) CellStyleBuilder {
	b.style.NumberFormat = format
	return b
}

func (b *cellStyleBuilder) Name(name string) CellStyleBuilder {
	b.style.Name = name
	return b
}

func (b *cellStyleBuilder) Link(url string) CellStyleBuilder {
	b.style.Link = url
	return b
}

func (b *cellStyleBuilder) Font(name string) CellStyleBuilder {
	b.style.Font = name
	return b
}

func (b *cellStyleBuilder) Alt(text string) CellStyleBuilder {
	b.style.Alt = text
	return b
}

func (b *cellStyleBuilder) LineSpacing(spacing float64) CellStyleBuilder {
	b.style.LineSpacing = spacing
	return b
}

func (b *cellStyleBuilder) SpacingBefore(spacing float64) CellStyleBuilder {
	b.style.SpacingBefore = spacing
	return b
}

func (b *cellStyleBuilder) SpacingAfter(spacing float64) CellStyleBuilder {
	b.style.SpacingAfter = spacing
	return b
}

func (b *cellStyleBuilder) Indent(points float64) CellStyleBuilder {
	b.style.Indent = points
	return b
}

func (b *cellStyleBuilder) Hanging(points float64) CellStyleBuilder {
	b.style.Hanging = points
	return b
}

func (b *cellStyleBuilder) Padding(points float64) CellStyleBuilder {
	b.style.Padding = points
	return b
}

func (b *cellStyleBuilder) KeepWithNext() CellStyleBuilder {
	b.style.KeepWithNext = true
	return b
}

func (b *cellStyleBuilder) KeepTogether() CellStyleBuilder {
	b.style.KeepTogether = true
	return b
}

func (b *cellStyleBuilder) Superscript() CellStyleBuilder {
	b.style.Superscript = true
	return b
}

func (b *cellStyleBuilder) Subscript() CellStyleBuilder {
	b.style.Subscript = true
	return b
}

func (b *cellStyleBuilder) Absolute() CellStyleBuilder {
	b.style.Absolute = true
	return b
}

func (b *cellStyleBuilder) Pos(x, y float64) CellStyleBuilder {
	b.style.X = x
	b.style.Y = y
	b.style.Absolute = true
	return b
}

func (b *cellStyleBuilder) Opacity(opacity float64) CellStyleBuilder {
	b.style.Opacity = opacity
	return b
}

func (b *cellStyleBuilder) DashPattern(pattern []float64) CellStyleBuilder {
	b.style.DashPattern = pattern
	return b
}

func (b *cellStyleBuilder) Build() CellStyle {
	return b.style
}
