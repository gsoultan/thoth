package pdf

import (
	"fmt"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

type tableRenderer struct {
	*state
}

func (p *tableRenderer) calculateRowHeight(ctx *renderingContext, item *contentItem, row int, colWidths []float64) float64 {
	maxH := 20.0
	for c := range item.cols {
		if len(item.cells[row][c]) > 0 && item.cells[row][c][0].hidden {
			continue
		}
		cW := colWidths[c]
		if len(item.cells[row][c]) > 0 && item.cells[row][c][0].colSpan > 1 {
			for i := 1; i < item.cells[row][c][0].colSpan && c+i < item.cols; i++ {
				cW += colWidths[c+i]
			}
		}
		cellH := p.calculateCellHeight(ctx, item.cells[row][c], cW)
		if cellH > maxH {
			maxH = cellH
		}
	}
	return maxH
}

func (p *tableRenderer) calculateCellHeight(ctx *renderingContext, cells []cellItem, width float64) float64 {
	cellH := 0.0
	maxPadding := 2.0
	for _, ci := range cells {
		if ci.style.Padding > maxPadding {
			maxPadding = ci.style.Padding
		}

		if ci.isImage {
			cellH += ci.height
		} else if ci.isTable {
			tableH := 0.0
			subColWidths := make([]float64, ci.table.cols)
			for c := range ci.table.cols {
				subColWidths[c] = width / float64(ci.table.cols)
			}
			for r := range ci.table.rows {
				tableH += p.calculateRowHeight(ctx, ci.table, r, subColWidths)
			}
			cellH += tableH
		} else if ci.isRich {
			h, _ := p.renderRichParagraph(ctx, nil, ci.spans, ci.style, 0, 0, width, 10, true)
			cellH += h
		} else if ci.isList {
			size := 10.0
			if ci.style.Size > 0 {
				size = float64(ci.style.Size)
			}
			fontName := p.getFontName(ci.style)
			customWidths, unitsPerEm := p.getCustomWidths(ci.style, ctx)

			listStyle := ci.style
			if listStyle.Hanging == 0 {
				listStyle.Hanging = 12
			}
			padding := 2.0
			if listStyle.Padding > 0 {
				padding = listStyle.Padding
			}

			for i, listItem := range ci.listItems {
				prefix := "• "
				if ci.ordered {
					prefix = fmt.Sprintf("%d. ", i+1)
				}
				firstLineMaxWidth := width - (padding * 2) - listStyle.Indent
				otherLinesMaxWidth := width - (padding * 2) - listStyle.Indent - listStyle.Hanging
				if firstLineMaxWidth < 10 {
					firstLineMaxWidth = 10
				}
				if otherLinesMaxWidth < 10 {
					otherLinesMaxWidth = 10
				}
				lines := wrapText(prefix+listItem, []float64{firstLineMaxWidth, otherLinesMaxWidth}, size, fontName, customWidths, unitsPerEm)
				lineSpacing := 1.2
				if ci.style.LineSpacing > 0 {
					lineSpacing = ci.style.LineSpacing
				}
				cellH += float64(len(lines)) * size * lineSpacing
				cellH += ci.style.SpacingBefore + ci.style.SpacingAfter
			}
		} else if ci.text != "" {
			size := 10.0
			if ci.style.Size > 0 {
				size = float64(ci.style.Size)
			}
			fontName := p.getFontName(ci.style)
			customWidths, unitsPerEm := p.getCustomWidths(ci.style, ctx)

			padding := 2.0
			if ci.style.Padding > 0 {
				padding = ci.style.Padding
			}

			firstLineMaxWidth := width - (padding * 2) - ci.style.Indent
			otherLinesMaxWidth := width - (padding * 2) - ci.style.Indent - ci.style.Hanging
			if firstLineMaxWidth < 10 {
				firstLineMaxWidth = 10
			}
			if otherLinesMaxWidth < 10 {
				otherLinesMaxWidth = 10
			}
			lines := wrapText(ci.text, []float64{firstLineMaxWidth, otherLinesMaxWidth}, size, fontName, customWidths, unitsPerEm)
			lineSpacing := 1.2
			if ci.style.LineSpacing > 0 {
				lineSpacing = ci.style.LineSpacing
			}
			cellH += float64(len(lines)) * size * lineSpacing
			cellH += ci.style.SpacingBefore + ci.style.SpacingAfter
		}
		cellH += 2.0
	}
	return cellH + (maxPadding * 2)
}

func (p *tableRenderer) calculateCellPreferredWidth(ctx *renderingContext, cells []cellItem) float64 {
	maxW := 0.0
	for _, ci := range cells {
		padding := 2.0
		if ci.style.Padding > 0 {
			padding = ci.style.Padding
		}

		if ci.isImage {
			if ci.width+(padding*2) > maxW {
				maxW = ci.width + (padding * 2)
			}
		} else if ci.isRich {
			for _, span := range ci.spans {
				words := strings.Fields(span.Text)
				for _, w := range words {
					sz := 10.0
					if span.Style.Size > 0 {
						sz = float64(span.Style.Size)
					}
					cw, upe := p.getCustomWidths(span.Style, ctx)
					wW := getTextWidth(w, sz, p.getFontName(span.Style), cw, upe)
					if wW+(padding*2) > maxW {
						maxW = wW + (padding * 2)
					}
				}
			}
		} else if ci.isList {
			for _, item := range ci.listItems {
				words := strings.Fields(item)
				for _, w := range words {
					sz := 10.0
					if ci.style.Size > 0 {
						sz = float64(ci.style.Size)
					}
					cw, upe := p.getCustomWidths(ci.style, ctx)
					wW := getTextWidth(w, sz, p.getFontName(ci.style), cw, upe)
					if wW+(padding*2)+ci.style.Indent+ci.style.Hanging > maxW {
						maxW = wW + (padding * 2) + ci.style.Indent + ci.style.Hanging
					}
				}
			}
		} else if ci.text != "" {
			words := strings.Fields(ci.text)
			for _, w := range words {
				sz := 10.0
				if ci.style.Size > 0 {
					sz = float64(ci.style.Size)
				}
				cw, upe := p.getCustomWidths(ci.style, ctx)
				wW := getTextWidth(w, sz, p.getFontName(ci.style), cw, upe)
				if wW+(padding*2)+ci.style.Indent+ci.style.Hanging > maxW {
					maxW = wW + (padding * 2) + ci.style.Indent + ci.style.Hanging
				}
			}
		} else if ci.isTable {
			// Nested table preferred width: sum of min widths of columns
			tw := 0.0
			for c := range ci.table.cols {
				maxColMinW := 0.0
				for r := range ci.table.rows {
					w := p.calculateCellPreferredWidth(ctx, ci.table.cells[r][c])
					if w > maxColMinW {
						maxColMinW = w
					}
				}
				tw += maxColMinW
			}
			if tw+(padding*2) > maxW {
				maxW = tw + (padding * 2)
			}
		}
	}
	return maxW
}

func (p *tableRenderer) calculateCellMaxWidth(ctx *renderingContext, cells []cellItem) float64 {
	maxW := 0.0
	for _, ci := range cells {
		padding := 2.0
		if ci.style.Padding > 0 {
			padding = ci.style.Padding
		}

		if ci.isImage {
			if ci.width+(padding*2) > maxW {
				maxW = ci.width + (padding * 2)
			}
		} else if ci.isRich {
			lineW := 0.0
			for _, span := range ci.spans {
				sz := 10.0
				if span.Style.Size > 0 {
					sz = float64(span.Style.Size)
				}
				cw, upe := p.getCustomWidths(span.Style, ctx)
				lineW += getTextWidth(span.Text, sz, p.getFontName(span.Style), cw, upe)
			}
			if lineW+(padding*2) > maxW {
				maxW = lineW + (padding * 2)
			}
		} else if ci.isList {
			for _, item := range ci.listItems {
				sz := 10.0
				if ci.style.Size > 0 {
					size := float64(ci.style.Size)
					sz = size
				}
				cw, upe := p.getCustomWidths(ci.style, ctx)
				lineW := getTextWidth(item, sz, p.getFontName(ci.style), cw, upe)
				if lineW+(padding*2)+ci.style.Indent+ci.style.Hanging+15 > maxW {
					maxW = lineW + (padding * 2) + ci.style.Indent + ci.style.Hanging + 15
				}
			}
		} else if ci.text != "" {
			lines := strings.Split(ci.text, "\n")
			for _, line := range lines {
				sz := 10.0
				if ci.style.Size > 0 {
					sz = float64(ci.style.Size)
				}
				cw, upe := p.getCustomWidths(ci.style, ctx)
				lineW := getTextWidth(line, sz, p.getFontName(ci.style), cw, upe)
				if lineW+(padding*2)+ci.style.Indent+ci.style.Hanging > maxW {
					maxW = lineW + (padding * 2) + ci.style.Indent + ci.style.Hanging
				}
			}
		} else if ci.isTable {
			// Nested table max width: sum of max widths of columns
			tw := 0.0
			for c := range ci.table.cols {
				maxColMaxW := 0.0
				for r := range ci.table.rows {
					w := p.calculateCellMaxWidth(ctx, ci.table.cells[r][c])
					if w > maxColMaxW {
						maxColMaxW = w
					}
				}
				tw += maxColMaxW
			}
			if tw+(padding*2) > maxW {
				maxW = tw + (padding * 2)
			}
		}
	}
	return maxW
}

func (p *tableRenderer) renderTable(ctx *renderingContext, item *contentItem, x, width float64) {
	colWidths := make([]float64, item.cols)
	if len(item.colWidths) == item.cols {
		copy(colWidths, item.colWidths)
		// Scale if necessary? For now assume widths are correct
	} else {
		// Auto-size columns relative to 'width'
		totalPrefWidth := 0.0
		prefWidths := make([]float64, item.cols)
		for c := range item.cols {
			maxCellW := 0.0
			for r := range item.rows {
				cellW := p.calculateCellPreferredWidth(ctx, item.cells[r][c])
				if cellW > maxCellW {
					maxCellW = cellW
				}
			}
			prefWidths[c] = maxCellW
			totalPrefWidth += maxCellW
		}

		if totalPrefWidth > 0 {
			for c := range item.cols {
				colWidths[c] = (prefWidths[c] / totalPrefWidth) * width
			}
		} else {
			for c := range item.cols {
				colWidths[c] = width / float64(item.cols)
			}
		}
	}

	rowHeights := make([]float64, item.rows)
	for r := range item.rows {
		rowHeights[r] = p.calculateRowHeight(ctx, item, r, colWidths)
	}

	// Adjust row heights for rowSpan
	for r := range item.rows {
		for c := range item.cols {
			if len(item.cells[r][c]) > 0 && item.cells[r][c][0].hidden {
				continue
			}
			if len(item.cells[r][c]) > 0 && item.cells[r][c][0].rowSpan > 1 {
				rs := item.cells[r][c][0].rowSpan
				cW := colWidths[c]
				if item.cells[r][c][0].colSpan > 1 {
					for i := 1; i < item.cells[r][c][0].colSpan && c+i < item.cols; i++ {
						cW += colWidths[c+i]
					}
				}
				cellH := p.calculateCellHeight(ctx, item.cells[r][c], cW)

				currentTotalH := 0.0
				for i := range min(rs, item.rows-r) {
					currentTotalH += rowHeights[r+i]
				}

				if cellH > currentTotalH {
					extra := cellH - currentTotalH
					lastRow := min(r+rs-1, item.rows-1)
					rowHeights[lastRow] += extra
				}
			}
		}
	}

	m := p.getMargins()
	// Accessibility: Tagged content for Table
	tableMcid := ctx.mcidCounter
	ctx.mcidCounter++
	ctx.currentSb.WriteString(fmt.Sprintf("/Table << /MCID %d >> BDC\n", tableMcid))
	ctx.currentStructs = append(ctx.currentStructs, objects.Dictionary{
		"Type": objects.Name("StructElem"),
		"S":    objects.Name("Table"),
		"P":    objects.Integer(0),
		"Pg":   objects.Integer(0),
		"K":    objects.Integer(tableMcid),
	})

	for r := range item.rows {
		rh := rowHeights[r]
		if ctx.posY-rh < m.Bottom {
			p.finishPage(ctx)
			// Repeat headers if any
			if item.HeaderRows > 0 && r >= item.HeaderRows {
				for h := range min(item.HeaderRows, item.rows) {
					p.renderTableRow(ctx, item, h, ctx.posY, x, colWidths, rowHeights)
					ctx.posY -= rowHeights[h]
				}
			}
		}
		p.renderTableRow(ctx, item, r, ctx.posY, x, colWidths, rowHeights)
		ctx.posY -= rh
	}
	ctx.currentSb.WriteString("EMC\n")
}

func (p *tableRenderer) renderTableRow(ctx *renderingContext, item *contentItem, r int, y, startX float64, colWidths []float64, rowHeights []float64) {
	rowHeight := rowHeights[r]

	// Accessibility: Tagged content for Table Row
	rowMcid := ctx.mcidCounter
	ctx.mcidCounter++
	ctx.currentSb.WriteString(fmt.Sprintf("/TR << /MCID %d >> BDC\n", rowMcid))
	ctx.currentStructs = append(ctx.currentStructs, objects.Dictionary{
		"Type": objects.Name("StructElem"),
		"S":    objects.Name("TR"),
		"P":    objects.Integer(0),
		"Pg":   objects.Integer(0),
		"K":    objects.Integer(rowMcid),
	})

	for c := range item.cols {
		if len(item.cells[r][c]) > 0 && item.cells[r][c][0].hidden {
			continue
		}
		cellX := startX
		for i := range c {
			cellX += colWidths[i]
		}

		cW := colWidths[c]
		cH := rowHeight
		if len(item.cells[r][c]) > 0 {
			if item.cells[r][c][0].colSpan > 1 {
				for i := 1; i < item.cells[r][c][0].colSpan && c+i < item.cols; i++ {
					cW += colWidths[c+i]
				}
			}
			if item.cells[r][c][0].rowSpan > 1 {
				cH = 0
				for i := range min(item.cells[r][c][0].rowSpan, item.rows-r) {
					cH += rowHeights[r+i]
				}
			}
		}

		// Draw cell background
		bg := ""
		if len(item.cells[r][c]) > 0 {
			bg = item.cells[r][c][0].style.Background
		}

		if bg == "" && strings.HasPrefix(item.text, "zebra:") {
			if r%2 == 1 {
				bg = strings.TrimPrefix(item.text, "zebra:")
			}
		}

		if bg != "" {
			rb, gb, bb := hexToRGB(bg)
			opacity := 1.0
			if len(item.cells[r][c]) > 0 {
				opacity = item.cells[r][c][0].style.Opacity
			}
			ctx.currentSb.WriteString("q ")
			if gs := p.getExtGState(ctx, opacity); gs != "" {
				ctx.currentSb.WriteString(fmt.Sprintf("/%s gs ", gs))
			}
			ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f rg %.2f %.2f %.2f %.2f re f Q\n", rb, gb, bb, cellX, y-cH, cW, cH))
		}

		// Vertical Alignment
		cellContentH := p.calculateCellHeight(ctx, item.cells[r][c], cW)
		offsetY := 0.0
		padding := 2.0
		if len(item.cells[r][c]) > 0 {
			if item.cells[r][c][0].style.Padding > 0 {
				padding = item.cells[r][c][0].style.Padding
			}
			vAlign := item.cells[r][c][0].style.Vertical
			if vAlign == "center" {
				offsetY = (cH - cellContentH) / 2.0
			} else if vAlign == "bottom" {
				offsetY = cH - cellContentH
			}
		}
		if offsetY < 0 {
			offsetY = 0
		}

		// Accessibility: Tagged content for Table Cell
		mcid := ctx.mcidCounter
		ctx.mcidCounter++
		ctx.currentSb.WriteString(fmt.Sprintf("/TD << /MCID %d >> BDC\n", mcid))

		ctx.currentStructs = append(ctx.currentStructs, objects.Dictionary{
			"Type": objects.Name("StructElem"),
			"S":    objects.Name("TD"),
			"P":    objects.Integer(0), // Parent (placeholder)
			"Pg":   objects.Integer(0), // Page (placeholder)
			"K":    objects.Integer(mcid),
		})

		currY := y - padding - offsetY
		for _, ci := range item.cells[r][c] {
			if ci.isImage {
				ix := cellX + padding
				if ci.style.Horizontal == "center" {
					ix = cellX + (cW-ci.width)/2.0
				} else if ci.style.Horizontal == "right" {
					ix = cellX + cW - ci.width - padding
				}
				imgName := ctx.imageNames[ci.path]
				if imgName != "" {
					ctx.currentSb.WriteString(fmt.Sprintf("q %.2f 0 0 %.2f %.2f %.2f cm /%s Do Q\n", ci.width, ci.height, ix, currY-ci.height, imgName))
				}
				currY -= ci.height + 2.0
			} else if ci.isTable {
				p.renderTable(ctx, ci.table, cellX+padding, cW-(padding*2))
				currY = ctx.posY // Table renderer updates posY
				currY -= 2.0
			} else if ci.isRich {
				size := 10.0
				if ci.style.Size > 0 {
					size = float64(ci.style.Size)
				}
				h, _ := p.renderRichParagraph(ctx, &ctx.currentSb, ci.spans, ci.style, cellX+padding, currY, cW-(padding*2), size, false)
				currY -= h
				currY -= 2.0
			} else if ci.isList {
				size := 10.0
				if ci.style.Size > 0 {
					size = float64(ci.style.Size)
				}
				listStyle := ci.style
				if listStyle.Hanging == 0 {
					listStyle.Hanging = 12
				}
				currY -= p.renderList(ctx, &ctx.currentSb, ci.listItems, ci.ordered, listStyle, cellX+padding, currY, cW-(padding*2), size)
				currY -= 2.0
			} else if ci.text != "" {
				size := 10.0
				if ci.style.Size > 0 {
					size = float64(ci.style.Size)
				}
				currY -= p.renderParagraph(ctx, &ctx.currentSb, ci.text, ci.style, cellX+padding, currY, cW-(padding*2), size)
				currY -= 2.0
			}
		}
		ctx.currentSb.WriteString("EMC\n")

		// Draw borders after content to ensure they are on top
		if len(item.cells[r][c]) > 0 {
			style := item.cells[r][c][0].style
			if style.Border || style.BorderTop || style.BorderBottom || style.BorderLeft || style.BorderRight {
				bWidth := 0.5
				if style.BorderWidth > 0 {
					bWidth = style.BorderWidth
				}
				bColor := "000000"
				if style.BorderColor != "" {
					bColor = style.BorderColor
				}
				rb, gb, bb := hexToRGB(bColor)

				ctx.currentSb.WriteString("q ")
				if len(style.DashPattern) > 0 {
					ctx.currentSb.WriteString("[")
					for i, d := range style.DashPattern {
						if i > 0 {
							ctx.currentSb.WriteString(" ")
						}
						ctx.currentSb.WriteString(fmt.Sprintf("%.2f", d))
					}
					ctx.currentSb.WriteString("] 0 d ")
				}
				ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f RG %.2f w ", rb, gb, bb, bWidth))
				if style.Border || style.BorderTop {
					ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f m %.2f %.2f l S ", cellX, y, cellX+cW, y))
				}
				if style.Border || style.BorderBottom {
					ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f m %.2f %.2f l S ", cellX, y-cH, cellX+cW, y-cH))
				}
				if style.Border || style.BorderLeft {
					ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f m %.2f %.2f l S ", cellX, y, cellX, y-cH))
				}
				if style.Border || style.BorderRight {
					ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f m %.2f %.2f l S ", cellX+cW, y, cellX+cW, y-cH))
				}
				ctx.currentSb.WriteString("Q\n")
			}
		}
	}
	ctx.currentSb.WriteString("EMC\n")
}

// Delegation helpers

func (p *tableRenderer) getFontName(style document.CellStyle) string {
	return (&renderer{p.state}).getFontName(style)
}

func (p *tableRenderer) getCustomWidths(style document.CellStyle, ctx *renderingContext) (map[rune]uint16, uint16) {
	return (&renderer{p.state}).getCustomWidths(style, ctx)
}

func (p *tableRenderer) renderParagraph(ctx *renderingContext, sb *strings.Builder, text string, style document.CellStyle, x, y, maxWidth, fontSize float64) float64 {
	return (&textRenderer{p.state}).renderParagraph(ctx, sb, text, style, x, y, maxWidth, fontSize)
}

func (p *tableRenderer) renderRichParagraph(ctx *renderingContext, sb *strings.Builder, spans []document.TextSpan, style document.CellStyle, x, y, maxWidth, fontSize float64, dryRun bool) (float64, float64) {
	return (&textRenderer{p.state}).renderRichParagraph(ctx, sb, spans, style, x, y, maxWidth, fontSize, dryRun)
}

func (p *tableRenderer) renderList(ctx *renderingContext, sb *strings.Builder, items []string, ordered bool, style document.CellStyle, x, y, maxWidth, fontSize float64) float64 {
	return (&textRenderer{p.state}).renderList(ctx, sb, items, ordered, style, x, y, maxWidth, fontSize)
}

func (p *tableRenderer) finishPage(ctx *renderingContext) {
	(&pageRenderer{p.state}).finishPage(ctx)
}

func (p *tableRenderer) getExtGState(ctx *renderingContext, opacity float64) string {
	return (&renderer{p.state}).getExtGState(ctx, opacity)
}
