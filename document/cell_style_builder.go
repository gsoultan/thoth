package document

// CellStyleBuilder is a fluent builder for CellStyle.
type CellStyleBuilder interface {
	Bold() CellStyleBuilder
	Italic() CellStyleBuilder
	Size(size int) CellStyleBuilder
	Color(hex string) CellStyleBuilder
	Background(hex string) CellStyleBuilder
	Border() CellStyleBuilder
	Align(h, v string) CellStyleBuilder
	NumberFormat(format string) CellStyleBuilder
	Name(name string) CellStyleBuilder
	Link(url string) CellStyleBuilder
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

func (b *cellStyleBuilder) Build() CellStyle {
	return b.style
}
