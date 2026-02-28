package pdf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

type textRenderer struct {
	*state
}

func (p *textRenderer) renderParagraph(ctx *renderingContext, sb *strings.Builder, text string, style document.CellStyle, x, y, maxWidth, fontSize float64) float64 {
	if text == "" && style.SpacingBefore == 0 && style.SpacingAfter == 0 {
		return 0
	}

	y -= style.SpacingBefore

	// Accessibility: Tagged content
	mcid := ctx.mcidCounter
	ctx.mcidCounter++
	tag := "P"
	if strings.HasPrefix(style.Name, "Heading") {
		tag = "H" + strings.TrimPrefix(style.Name, "Heading")
	} else if style.Name == "Title" {
		tag = "H1"
	}
	sb.WriteString(fmt.Sprintf("/%s << /MCID %d >> BDC\n", tag, mcid))

	// Add to structure tree
	ctx.currentStructs = append(ctx.currentStructs, objects.Dictionary{
		"Type": objects.Name("StructElem"),
		"S":    objects.Name(tag),
		"P":    objects.Integer(0), // Parent (placeholder)
		"Pg":   objects.Integer(0), // Page (placeholder)
		"K":    objects.Integer(mcid),
	})

	// Add bookmark for headings
	if style.Name != "" && (strings.HasPrefix(style.Name, "Heading") || style.Name == "Title") {
		// Only add bookmarks for contentItems, not headers/footers
		// Check if we are rendering into currentSb (main content)
		if sb == &ctx.currentSb {
			lvl := 1
			if strings.HasPrefix(style.Name, "Heading") {
				if n := strings.TrimPrefix(style.Name, "Heading"); n != "" {
					if v, err := strconv.Atoi(n); err == nil {
						lvl = v
					}
				}
			}
			ctx.bookmarks = append(ctx.bookmarks, bookmark{
				title: text,
				level: lvl,
				page:  len(ctx.pages),
				posY:  y,
			})
		}
	}

	fontName := p.getFontName(style)
	customWidths, unitsPerEm := p.getCustomWidths(style, ctx)
	size := fontSize
	if style.Size > 0 {
		size = float64(style.Size)
	}

	// Map common characters to WinAnsiEncoding equivalents
	text = strings.ReplaceAll(text, "•", "\x95")

	lineSpacing := 1.2
	if style.LineSpacing > 0 {
		lineSpacing = style.LineSpacing
	}

	totalHeight := 0.0
	paragraphs := strings.Split(text, "\n")
	for _, para := range paragraphs {
		firstLineMaxWidth := maxWidth - style.Indent
		otherLinesMaxWidth := maxWidth - style.Indent - style.Hanging
		if firstLineMaxWidth < 20 {
			firstLineMaxWidth = 20
		}
		if otherLinesMaxWidth < 20 {
			otherLinesMaxWidth = 20
		}

		lines := wrapText(para, []float64{firstLineMaxWidth, otherLinesMaxWidth}, size, fontName, customWidths, unitsPerEm)
		for i, line := range lines {
			currIndent := style.Indent
			currentMaxWidth := firstLineMaxWidth
			if i > 0 {
				currIndent += style.Hanging
				currentMaxWidth = otherLinesMaxWidth
			}

			// Page break check (only if sb is main content and not KeepTogether)
			m := p.getMargins()
			if !style.KeepTogether && sb == &ctx.currentSb && y < m.Bottom {
				p.finishPage(ctx)
				sb = &ctx.currentSb
				y = float64(ctx.h) - m.Top
				x = m.Left // Reset to left margin
			}

			r, g, b := 0.0, 0.0, 0.0
			if style.Color != "" {
				r, g, b = hexToRGB(style.Color)
			}
			offsetX := 0.0
			if style.Horizontal == "center" {
				offsetX = (currentMaxWidth - getTextWidth(line, size, fontName, customWidths, unitsPerEm)) / 2.0
			} else if style.Horizontal == "right" {
				offsetX = currentMaxWidth - getTextWidth(line, size, fontName, customWidths, unitsPerEm)
			}

			if style.Link != "" && sb == &ctx.currentSb {
				ctx.currentLinks = append(ctx.currentLinks, link{
					rect: []float64{x + currIndent + offsetX, y - size, x + currIndent + offsetX + getTextWidth(line, size, fontName, customWidths, unitsPerEm), y},
					url:  style.Link,
				})
			}

			if style.Background != "" {
				rb, gb, bb := hexToRGB(style.Background)
				width := getTextWidth(line, size, fontName, customWidths, unitsPerEm)
				// Draw highlight rectangle before text
				sb.WriteString("q ")
				if gs := p.getExtGState(ctx, style.Opacity); gs != "" {
					sb.WriteString(fmt.Sprintf("/%s gs ", gs))
				}
				sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f rg %.2f %.2f %.2f %.2f re f Q\n", rb, gb, bb, x+currIndent+offsetX, y-size, width, size))
			}

			// Use y-size as baseline
			font := p.getFont(style, ctx)
			sb.WriteString("BT ")
			if gs := p.getExtGState(ctx, style.Opacity); gs != "" {
				sb.WriteString(fmt.Sprintf("/%s gs ", gs))
			}
			sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f rg %s %.2f Tf ", r, g, b, font, size))
			if style.Horizontal == "justify" && i < len(lines)-1 {
				wordCount := strings.Count(line, " ")
				if wordCount > 0 {
					extraSpace := currentMaxWidth - getTextWidth(line, size, fontName, customWidths, unitsPerEm)
					wordSpacing := extraSpace / float64(wordCount)
					sb.WriteString(fmt.Sprintf("%.2f Tw ", wordSpacing))
				}
			}
			sb.WriteString(fmt.Sprintf("%.2f %.2f Td (%s) Tj 0 Tw ET\n", x+currIndent+offsetX, y-size, escapePDF(line)))
			y -= (size * lineSpacing)
			totalHeight += (size * lineSpacing)
		}
	}
	sb.WriteString("EMC\n")

	y -= style.SpacingAfter
	totalHeight += style.SpacingAfter
	return totalHeight + style.SpacingBefore
}

func (p *textRenderer) renderList(ctx *renderingContext, sb *strings.Builder, items []string, ordered bool, style document.CellStyle, x, y, maxWidth, fontSize float64) float64 {
	startY := y
	listStyle := style
	fontName := p.getFontName(style)
	customWidths, unitsPerEm := p.getCustomWidths(style, ctx)
	size := fontSize
	if style.Size > 0 {
		size = float64(style.Size)
	}

	for i, item := range items {
		prefix := "• "
		if ordered {
			prefix = fmt.Sprintf("%d. ", i+1)
		}
		prefixWidth := getTextWidth(prefix, size, fontName, customWidths, unitsPerEm)
		itemStyle := listStyle
		if itemStyle.Hanging == 0 {
			itemStyle.Hanging = prefixWidth + 2 // Add small gap
		}
		y -= p.renderParagraph(ctx, sb, prefix+item, itemStyle, x, y, maxWidth, fontSize)
	}
	return startY - y
}

func (p *textRenderer) renderRichParagraph(ctx *renderingContext, sb *strings.Builder, spans []document.TextSpan, style document.CellStyle, x, y, maxWidth, fontSize float64, dryRun bool) (float64, float64) {
	if len(spans) == 0 {
		return 0, 0
	}

	startY := y
	if !dryRun {
		y -= style.SpacingBefore
		// Accessibility: Tagged content
		mcid := ctx.mcidCounter
		ctx.mcidCounter++
		tag := "P"
		if strings.HasPrefix(style.Name, "Heading") {
			tag = "H" + strings.TrimPrefix(style.Name, "Heading")
		} else if style.Name == "Title" {
			tag = "H1"
		}
		sb.WriteString(fmt.Sprintf("/%s << /MCID %d >> BDC\n", tag, mcid))

		// Add to structure tree
		ctx.currentStructs = append(ctx.currentStructs, objects.Dictionary{
			"Type": objects.Name("StructElem"),
			"S":    objects.Name(tag),
			"P":    objects.Integer(0), // Parent (placeholder)
			"Pg":   objects.Integer(0), // Page (placeholder)
			"K":    objects.Integer(mcid),
		})
	}

	// Add bookmark for headings in rich text
	if style.Name != "" && (strings.HasPrefix(style.Name, "Heading") || style.Name == "Title") && !dryRun {
		if sb == &ctx.currentSb {
			lvl := 1
			if strings.HasPrefix(style.Name, "Heading") {
				if n := strings.TrimPrefix(style.Name, "Heading"); n != "" {
					if v, err := strconv.Atoi(n); err == nil {
						lvl = v
					}
				}
			}
			ctx.bookmarks = append(ctx.bookmarks, bookmark{title: concatSpanText(spans), level: lvl, page: len(ctx.pages), posY: y})
		}
	}

	// Simple word-by-word wrapping for rich text
	var words []wordInfo
	for i, span := range spans {
		parts := strings.Split(span.Text, " ")
		for j, part := range parts {
			t := part
			if j < len(parts)-1 {
				t += " "
			}
			sz := fontSize
			if span.Style.Size > 0 {
				sz = float64(span.Style.Size)
			}
			if span.Style.Superscript || span.Style.Subscript {
				sz *= 0.6
			}
			fName := p.getFontName(span.Style)
			cWidths, uPE := p.getCustomWidths(span.Style, ctx)
			w := getTextWidth(t, sz, fName, cWidths, uPE)
			words = append(words, wordInfo{text: t, width: w, spanIdx: i})
		}
	}

	lineIdx := 0
	firstLineMaxWidth := maxWidth - style.Indent
	otherLinesMaxWidth := maxWidth - style.Indent - style.Hanging
	if firstLineMaxWidth < 20 {
		firstLineMaxWidth = 20
	}
	if otherLinesMaxWidth < 20 {
		otherLinesMaxWidth = 20
	}

	currentLine := []wordInfo{}
	currentLineWidth := 0.0

	for _, w := range words {
		limit := firstLineMaxWidth
		if lineIdx > 0 {
			limit = otherLinesMaxWidth
		}
		if currentLineWidth+w.width <= limit || len(currentLine) == 0 {
			currentLine = append(currentLine, w)
			currentLineWidth += w.width
		} else {
			// Process full line
			lh, _ := p.renderRichLine(ctx, sb, currentLine, spans, style, x, y, limit, lineIdx, false, dryRun)
			y -= lh
			currentLine = []wordInfo{w}
			currentLineWidth = w.width
			lineIdx++

			// Page break check (only if not dryRun and in main content)
			m := p.getMargins()
			if !dryRun && !style.KeepTogether && sb == &ctx.currentSb && y < m.Bottom {
				p.finishPage(ctx)
				sb = &ctx.currentSb
				y = float64(ctx.h) - m.Top
				x = m.Left // Reset to left margin
			}
		}
	}
	// Process last line
	lh, _ := p.renderRichLine(ctx, sb, currentLine, spans, style, x, y, otherLinesMaxWidth, lineIdx, true, dryRun)
	y -= lh

	if !dryRun {
		sb.WriteString("EMC\n")
		y -= style.SpacingAfter
	}

	return startY - y, startY - y
}

func (p *textRenderer) renderRichLine(ctx *renderingContext, sb *strings.Builder, line []wordInfo, spans []document.TextSpan, style document.CellStyle, x, y, limit float64, lineIdx int, isLast, dryRun bool) (float64, float64) {
	maxLineH := 0.0
	lineWidth := 0.0
	for _, w := range line {
		sz := 12.0
		if spans[w.spanIdx].Style.Size > 0 {
			sz = float64(spans[w.spanIdx].Style.Size)
		}
		lineSpacing := 1.2
		if style.LineSpacing > 0 {
			lineSpacing = style.LineSpacing
		}
		if sz*lineSpacing > maxLineH {
			maxLineH = sz * lineSpacing
		}
		lineWidth += w.width
	}

	if dryRun {
		return maxLineH, lineWidth
	}

	currIndent := style.Indent
	if lineIdx > 0 {
		currIndent += style.Hanging
	}

	offsetX := 0.0
	if style.Horizontal == "center" {
		offsetX = (limit - lineWidth) / 2.0
	} else if style.Horizontal == "right" {
		offsetX = limit - lineWidth
	}

	wordSpacing := 0.0
	if style.Horizontal == "justify" && !isLast {
		spaces := 0
		for i := range len(line) - 1 {
			if strings.HasSuffix(line[i].text, " ") {
				spaces++
			}
		}
		if spaces > 0 {
			wordSpacing = (limit - lineWidth) / float64(spaces)
		}
	}

	currX := x + currIndent + offsetX
	for i, w := range line {
		span := spans[w.spanIdx]
		sz := 12.0
		if span.Style.Size > 0 {
			sz = float64(span.Style.Size)
		}
		if span.Style.Superscript || span.Style.Subscript {
			sz *= 0.6
		}

		font := p.getFont(span.Style, ctx)
		r, g, b := 0.0, 0.0, 0.0
		if span.Style.Color != "" {
			r, g, b = hexToRGB(span.Style.Color)
		}

		baselineY := y - sz
		if span.Style.Superscript {
			baselineY = y - (maxLineH * 0.4)
		} else if span.Style.Subscript {
			baselineY = y - sz - (maxLineH * 0.1)
		}

		if span.Style.Link != "" && sb == &ctx.currentSb {
			ctx.currentLinks = append(ctx.currentLinks, link{
				rect: []float64{currX, baselineY, currX + w.width, baselineY + sz},
				url:  span.Style.Link,
			})
		}

		if span.Style.Background != "" {
			rb, gb, bb := hexToRGB(span.Style.Background)
			sb.WriteString("q ")
			if gs := p.getExtGState(ctx, span.Style.Opacity); gs != "" {
				sb.WriteString(fmt.Sprintf("/%s gs ", gs))
			}
			sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f rg %.2f %.2f %.2f %.2f re f Q\n", rb, gb, bb, currX, baselineY, w.width, sz))
		}

		sb.WriteString("BT ")
		if gs := p.getExtGState(ctx, span.Style.Opacity); gs != "" {
			sb.WriteString(fmt.Sprintf("/%s gs ", gs))
		}
		sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f rg %s %.2f Tf ", r, g, b, font, sz))
		if strings.HasSuffix(w.text, " ") && i < len(line)-1 {
			sb.WriteString(fmt.Sprintf("%.2f Tw ", wordSpacing))
		}
		sb.WriteString(fmt.Sprintf("%.2f %.2f Td (%s) Tj 0 Tw ET\n", currX, baselineY, escapePDF(w.text)))
		currX += w.width
		if strings.HasSuffix(w.text, " ") && i < len(line)-1 {
			currX += wordSpacing
		}
	}

	return maxLineH, lineWidth
}

func concatSpanText(spans []document.TextSpan) string {
	b := getSB()
	defer putSB(b)
	for _, s := range spans {
		b.WriteString(s.Text)
	}
	return b.String()
}

func (p *textRenderer) renderHeader(ctx *renderingContext, sb *strings.Builder) {
	m := p.getMargins()
	currY := float64(ctx.h) - (m.Top / 2.0)
	for _, item := range p.header {
		if item.isParagraph {
			currY -= p.renderParagraph(ctx, sb, item.text, item.style, m.Left, currY, float64(ctx.w)-(m.Left+m.Right), 10)
		}
	}
}

func (p *textRenderer) renderFooter(ctx *renderingContext, sb *strings.Builder, pageNum, totalPages int) {
	m := p.getMargins()
	currY := m.Bottom / 2.0
	for _, item := range p.footer {
		if item.isParagraph {
			text := strings.ReplaceAll(item.text, "{n}", fmt.Sprintf("%d", pageNum))
			text = strings.ReplaceAll(text, "{nb}", fmt.Sprintf("%d", totalPages))
			p.renderParagraph(ctx, sb, text, item.style, m.Left, currY, float64(ctx.w)-(m.Left+m.Right), 10)
			currY -= 12
		}
	}
}

func (p *textRenderer) renderWatermark(ctx *renderingContext, sb *strings.Builder) {
	if p.watermark == nil {
		return
	}

	style := p.watermark.style
	if style.Size == 0 {
		style.Size = 60
	}
	if style.Color == "" {
		style.Color = "CCCCCC" // Light gray
	}

	// Center of page
	x := float64(ctx.w) / 2.0
	y := float64(ctx.h) / 2.0

	// Rotation is a bit complex in PDF content stream without a helper.
	// We'll use a transformation matrix for 45 degrees rotation.
	// cos(45) = sin(45) = 0.707
	sb.WriteString("q\n")
	sb.WriteString(fmt.Sprintf("0.707 0.707 -0.707 0.707 %.2f %.2f cm\n", x, y))
	// After rotation, (0,0) is at center of page.
	// We want to center the text at (0,0).
	fontName := p.getFontName(style)
	customWidths, unitsPerEm := p.getCustomWidths(style, ctx)
	textWidth := getTextWidth(p.watermark.text, float64(style.Size), fontName, customWidths, unitsPerEm)
	p.renderParagraph(ctx, sb, p.watermark.text, style, -textWidth/2.0, 0, textWidth+10, float64(style.Size))
	sb.WriteString("Q\n")
}

func (p *textRenderer) renderTOC(ctx *renderingContext, colWidth float64) {
	if len(ctx.bookmarks) == 0 {
		return
	}

	m := p.getMargins()
	titleStyle := document.CellStyle{Size: 16, Bold: true}
	ctx.posY -= p.renderParagraph(ctx, &ctx.currentSb, "Table of Contents", titleStyle, m.Left, ctx.posY, colWidth, 16)
	ctx.posY -= 20

	for _, bm := range ctx.bookmarks {
		if bm.level > 3 || bm.level == 0 {
			// Skip bookmarks that aren't headings or are too deep
			if bm.level == 0 && !strings.Contains(strings.ToLower(bm.title), "bookmark") {
				// unless it looks like a manual bookmark we want in TOC?
				// Usually only headings go to TOC.
				continue
			}
			if bm.level == 0 {
				continue
			}
		}

		indent := float64(bm.level-1) * 20.0
		entryStyle := document.CellStyle{Size: 10, Link: "#" + bm.title}

		// Entry title
		h := p.renderParagraph(ctx, &ctx.currentSb, bm.title, entryStyle, m.Left+indent, ctx.posY, colWidth-indent-40, 10)

		// Page number
		pageStr := fmt.Sprintf("%d", bm.page+1)
		p.renderParagraph(ctx, &ctx.currentSb, pageStr, document.CellStyle{Size: 10, Horizontal: "right"}, m.Left, ctx.posY, colWidth, 10)

		ctx.posY -= h
		ctx.posY -= 5

		if ctx.posY < m.Bottom+30 {
			p.finishPage(ctx)
		}
	}
	ctx.posY -= 20
}

// Delegation helpers

func (p *textRenderer) getFontName(style document.CellStyle) string {
	return (&renderer{p.state}).getFontName(style)
}

func (p *textRenderer) getFont(style document.CellStyle, ctx *renderingContext) string {
	return (&renderer{p.state}).getFont(style, ctx)
}

func (p *textRenderer) getExtGState(ctx *renderingContext, opacity float64) string {
	return (&renderer{p.state}).getExtGState(ctx, opacity)
}

func (p *textRenderer) getCustomWidths(style document.CellStyle, ctx *renderingContext) (map[rune]uint16, uint16) {
	return (&renderer{p.state}).getCustomWidths(style, ctx)
}

func (p *textRenderer) finishPage(ctx *renderingContext) {
	(&pageRenderer{p.state}).finishPage(ctx)
}
