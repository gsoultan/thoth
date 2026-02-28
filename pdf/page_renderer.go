package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

type pageRenderer struct {
	*state
}

func (p *pageRenderer) finishPage(ctx *renderingContext) {
	m := p.getMargins()
	if ctx.currentSb.Len() == 0 && len(ctx.pages) > 0 && len(ctx.currentFields) == 0 && len(ctx.currentFootnotes) == 0 {
		return
	}

	// Multi-column handling: check if we should just move to next column
	if p.pageSettings.Columns > 1 && ctx.currentColumn < p.pageSettings.Columns-1 {
		ctx.currentColumn++
		ctx.posY = float64(ctx.h) - m.Top
		return
	}

	ctx.pages = append(ctx.pages, pageInfo{
		sb:          ctx.currentSb.String(),
		links:       ctx.currentLinks,
		fields:      ctx.currentFields,
		footnotes:   ctx.currentFootnotes,
		structItems: ctx.currentStructs,
		w:           ctx.w,
		h:           ctx.h,
	})

	ctx.currentSb.Reset()
	ctx.currentLinks = nil
	ctx.currentFields = nil
	ctx.currentFootnotes = nil
	ctx.currentStructs = nil
	ctx.currentColumn = 0
	ctx.posY = float64(ctx.h) - m.Top
}

func (p *pageRenderer) finalizePages(ctx *renderingContext) {
	m := p.getMargins()
	total := len(ctx.pages)

	// First pass: create all page object placeholders to get their references for annotations/links
	for range total {
		ref := ctx.mgr.AddObject(objects.Dictionary{"Type": objects.Name("Page")})
		ctx.pageRefs = append(ctx.pageRefs, ref)
	}

	for i := range total {
		pageNum := i + 1
		sb := getSB()

		// Render header and footer into a new builder
		p.renderHeader(ctx, sb)
		p.renderFooter(ctx, sb, pageNum, total)
		p.renderWatermark(ctx, sb)

		// Append original page contentItems
		sb.WriteString(ctx.pages[i].sb)

		// Render footnotes
		if len(ctx.pages[i].footnotes) > 0 {
			fy := m.Bottom + 10
			sb.WriteString(fmt.Sprintf("q 0.5 w %.2f %.2f m %.2f %.2f l S Q\n", m.Left, m.Bottom+15, float64(ctx.w)-m.Right, m.Bottom+15)) // Line above footnotes
			for _, fn := range ctx.pages[i].footnotes {
				p.renderParagraph(ctx, sb, fn, document.CellStyle{Size: 8, Color: "666666"}, m.Left, fy, float64(ctx.w)-(m.Left+m.Right), 8)
				fy -= 10
			}
		}

		// Stream compression
		data := []byte(sb.String())
		putSB(sb)
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write(data)
		zw.Close()
		compressed := buf.Bytes()

		stream := objects.Stream{
			Dict: objects.Dictionary{"Filter": objects.Name("FlateDecode")},
			Data: compressed,
		}

		streamRef := ctx.mgr.AddObject(stream)
		xobjects := objects.Dictionary{}
		for path, name := range ctx.imageNames {
			xobjects[name] = ctx.imageRefs[path]
		}

		// Annotations (Links & Forms)
		annots := objects.Array{}
		for _, l := range ctx.pages[i].links {
			annot := objects.Dictionary{
				"Type":    objects.Name("Annot"),
				"Subtype": objects.Name("Link"),
				"Rect": objects.Array{
					objects.Integer(int(l.rect[0])),
					objects.Integer(int(l.rect[1])),
					objects.Integer(int(l.rect[2])),
					objects.Integer(int(l.rect[3])),
				},
				"Border": objects.Array{objects.Integer(0), objects.Integer(0), objects.Integer(0)},
			}

			if strings.HasPrefix(l.url, "#") {
				target := strings.TrimPrefix(l.url, "#")
				found := false
				for _, bm := range ctx.bookmarks {
					if bm.title == target {
						pageIdx := bm.page
						if pageIdx >= len(ctx.pageRefs) {
							pageIdx = len(ctx.pageRefs) - 1
						}
						if pageIdx < 0 {
							pageIdx = 0
						}
						annot["Dest"] = objects.Array{ctx.pageRefs[pageIdx], objects.Name("XYZ"), objects.Integer(0), objects.Integer(int(bm.posY)), objects.Integer(0)}
						found = true
						break
					}
				}
				if !found {
					// Fallback to URI
					annot["A"] = objects.Dictionary{
						"Type": objects.Name("Action"),
						"S":    objects.Name("URI"),
						"URI":  objects.PDFString(l.url),
					}
				}
			} else {
				annot["A"] = objects.Dictionary{
					"Type": objects.Name("Action"),
					"S":    objects.Name("URI"),
					"URI":  objects.PDFString(l.url),
				}
			}
			annots = append(annots, annot)
		}

		for _, f := range ctx.pages[i].fields {
			annot := objects.Dictionary{
				"Type":    objects.Name("Annot"),
				"Subtype": objects.Name("Widget"),
				"Rect": objects.Array{
					objects.Integer(int(f.x1)),
					objects.Integer(int(f.y1)),
					objects.Integer(int(f.x1 + f.width)),
					objects.Integer(int(f.y1 + f.height)),
				},
				"T": objects.PDFString(f.fieldName),
				"P": ctx.pageRefs[i], // Link to page
			}

			// Add basic appearance characteristics
			mk := objects.Dictionary{
				"BC": objects.Array{objects.Integer(0)}, // Black border
				"BG": objects.Array{objects.Integer(1)}, // White background
			}
			annot["MK"] = mk

			switch f.fieldType {
			case "text":
				annot["FT"] = objects.Name("Tx")
				annot["DA"] = objects.PDFString("/Helv 10 Tf 0 g")
			case "checkbox":
				annot["FT"] = objects.Name("Btn")
				annot["Ff"] = objects.Integer(0) // Checkbox
				annot["V"] = objects.Name("Off")
				annot["DV"] = objects.Name("Off")
			case "combobox":
				annot["FT"] = objects.Name("Ch")
				annot["Ff"] = objects.Integer(131072) // Combo
				opts := objects.Array{}
				for _, o := range f.options {
					opts = append(opts, objects.PDFString(o))
				}
				annot["Opt"] = opts
			case "radio":
				annot["FT"] = objects.Name("Btn")
				annot["Ff"] = objects.Integer(32768) // Radio
			}

			// Add to manager to get a reference
			fieldRef := ctx.mgr.AddObject(annot)
			ctx.allFields = append(ctx.allFields, fieldRef)
			annots = append(annots, fieldRef)
		}

		resDict := objects.Dictionary{
			"ProcSet": objects.Array{objects.Name("PDF"), objects.Name("Text"), objects.Name("ImageB"), objects.Name("ImageC"), objects.Name("ImageI")},
			"Font": objects.Dictionary{
				"F1": objects.Dictionary{"Type": objects.Name("Font"), "Subtype": objects.Name("Type1"), "BaseFont": objects.Name("Helvetica"), "Encoding": objects.Name("WinAnsiEncoding")},
				"F2": objects.Dictionary{"Type": objects.Name("Font"), "Subtype": objects.Name("Type1"), "BaseFont": objects.Name("Helvetica-Bold"), "Encoding": objects.Name("WinAnsiEncoding")},
				"F3": objects.Dictionary{"Type": objects.Name("Font"), "Subtype": objects.Name("Type1"), "BaseFont": objects.Name("Helvetica-Oblique"), "Encoding": objects.Name("WinAnsiEncoding")},
				"F4": objects.Dictionary{"Type": objects.Name("Font"), "Subtype": objects.Name("Type1"), "BaseFont": objects.Name("Helvetica-BoldOblique"), "Encoding": objects.Name("WinAnsiEncoding")},
			},
			"XObject": xobjects,
		}

		// Add custom fonts to resources
		for name, ref := range ctx.fontRefs {
			fontKey := ctx.fontNames[name]
			resDict["Font"].(objects.Dictionary)[fontKey] = ref
		}

		// Add ExtGStates to resources
		if len(ctx.extGStates) > 0 {
			egs := objects.Dictionary{}
			for op, ref := range ctx.extGStates {
				egs[fmt.Sprintf("GS%d", ref.Number)] = ref
				_ = op // use op if needed
			}
			resDict["ExtGState"] = egs
		}

		// Add imported pages to resources
		for _, ref := range ctx.importRefs {
			resDict["XObject"].(objects.Dictionary)[fmt.Sprintf("Imp%d", ref.Number)] = ref
		}

		page := objects.Dictionary{
			"Type":      objects.Name("Page"),
			"MediaBox":  objects.Array{objects.Integer(0), objects.Integer(0), objects.Integer(ctx.pages[i].w), objects.Integer(ctx.pages[i].h)},
			"Resources": resDict,
			"Contents":  streamRef,
		}

		if len(annots) > 0 {
			page["Annots"] = annots
		}

		// Accessibility: Structural order
		page["Tabs"] = objects.Name("S")

		pageRef := ctx.pageRefs[i]
		// Update the existing object in manager
		for j := range ctx.mgr.Objects {
			if ctx.mgr.Objects[j].Number == pageRef.Number {
				ctx.mgr.Objects[j].Data = page
			}
		}
	}
}

func (p *pageRenderer) renderContent(ctx *renderingContext) {
	m := p.getMargins()
	for _, item := range p.contentItems {
		marginX := m.Left + m.Right
		colWidth := float64(ctx.w) - marginX
		if p.pageSettings.Columns > 1 {
			totalGap := float64(p.pageSettings.Columns-1) * p.pageSettings.ColumnGap
			colWidth = (float64(ctx.w) - marginX - totalGap) / float64(p.pageSettings.Columns)
		}

		x := m.Left + (float64(ctx.currentColumn) * (colWidth + p.pageSettings.ColumnGap))
		y := ctx.posY

		if item.style.Absolute {
			x = item.style.X
			y = item.style.Y
		}

		if item.isParagraph || item.isHeading {
			h := p.renderParagraph(ctx, &ctx.currentSb, item.text, item.style, x, y, colWidth, 12)
			if !item.style.Absolute {
				ctx.posY -= h
			}
		} else if item.isRich {
			h, _ := p.renderRichParagraph(ctx, &ctx.currentSb, item.spans, item.style, x, y, colWidth, 12, false)
			if !item.style.Absolute {
				ctx.posY -= h
			}
		} else if item.isList {
			h := p.renderList(ctx, &ctx.currentSb, item.listItems, item.ordered, item.style, x, y, colWidth, 12)
			if !item.style.Absolute {
				ctx.posY -= h
			}
		} else if item.isImage {
			if !item.style.Absolute && ctx.posY-item.height < m.Bottom {
				p.finishPage(ctx)
				x = m.Left + (float64(ctx.currentColumn) * (colWidth + p.pageSettings.ColumnGap))
				y = ctx.posY
			}
			if item.style.Absolute {
				x = item.style.X
				y = item.style.Y
			}

			imgName := ctx.imageNames[item.path]
			if imgName != "" {
				// Accessibility: Tagged content for Image
				mcid := ctx.mcidCounter
				ctx.mcidCounter++
				ctx.currentSb.WriteString(fmt.Sprintf("/Figure << /MCID %d >> BDC\n", mcid))

				dict := objects.Dictionary{
					"Type": objects.Name("StructElem"),
					"S":    objects.Name("Figure"),
					"P":    objects.Integer(0), // Parent (placeholder)
					"Pg":   objects.Integer(0), // Page (placeholder)
					"K":    objects.Integer(mcid),
				}
				if item.style.Alt != "" {
					dict["Alt"] = objects.PDFString(item.style.Alt)
				}
				ctx.currentStructs = append(ctx.currentStructs, dict)

				ctx.currentSb.WriteString(fmt.Sprintf("q %.2f 0 0 %.2f %.2f %.2f cm /%s Do Q\n", item.width, item.height, x, y-item.height, imgName))
				ctx.currentSb.WriteString("EMC\n")
			}
			if !item.style.Absolute {
				ctx.posY -= item.height + 10
			}
		} else if item.isTable {
			p.renderTable(ctx, item, x, colWidth)
		} else if item.isPageBreak {
			p.finishPage(ctx)
			ctx.currentColumn = 0 // Reset column on hard page break
			ctx.posY = float64(ctx.h) - m.Top
		} else if item.isShape {
			p.renderShape(ctx, item)
		} else if item.isTOC {
			p.renderTOC(ctx, colWidth)
		} else if item.isBookmark {
			ctx.bookmarks = append(ctx.bookmarks, bookmark{
				title: item.bookmarkName,
				page:  len(ctx.pages),
				posY:  ctx.posY,
			})
		} else if item.isImported {
			// Placeholder for imported page rendering
			imgKey := fmt.Sprintf("Imp%d", ctx.importRefs[fmt.Sprintf("%s:%d", item.importPath, item.importPage)].Number)
			ctx.currentSb.WriteString(fmt.Sprintf("q 1 0 0 1 0 0 cm /%s Do Q\n", imgKey))
		} else if item.isFormField {
			// Form fields positioning
			f := item
			f.x1 = x + f.x1 // Offset relative to column
			ctx.currentFields = append(ctx.currentFields, f)
		} else if item.isFootnote {
			ctx.currentFootnotes = append(ctx.currentFootnotes, item.text)
		}
	}
}

func (p *pageRenderer) renderShape(ctx *renderingContext, item *contentItem) {
	r, g, b := 0.0, 0.0, 0.0
	if item.style.Color != "" {
		r, g, b = hexToRGB(item.style.Color)
	}
	rb, gb, bb := 0.0, 0.0, 0.0
	hasBg := false
	if item.style.Background != "" {
		rb, gb, bb = hexToRGB(item.style.Background)
		hasBg = true
	}

	// Accessibility: Tagged content for Shape
	mcid := ctx.mcidCounter
	ctx.mcidCounter++
	ctx.currentSb.WriteString(fmt.Sprintf("/Figure << /MCID %d >> BDC\n", mcid))

	dict := objects.Dictionary{
		"Type": objects.Name("StructElem"),
		"S":    objects.Name("Figure"),
		"P":    objects.Integer(0), // Parent (placeholder)
		"Pg":   objects.Integer(0), // Page (placeholder)
		"K":    objects.Integer(mcid),
	}
	if item.style.Alt != "" {
		dict["Alt"] = objects.PDFString(item.style.Alt)
	}
	ctx.currentStructs = append(ctx.currentStructs, dict)

	ctx.currentSb.WriteString("q ")
	if gs := p.getExtGState(ctx, item.style.Opacity); gs != "" {
		ctx.currentSb.WriteString(fmt.Sprintf("/%s gs ", gs))
	}
	if len(item.style.DashPattern) > 0 {
		ctx.currentSb.WriteString("[")
		for i, d := range item.style.DashPattern {
			if i > 0 {
				ctx.currentSb.WriteString(" ")
			}
			ctx.currentSb.WriteString(fmt.Sprintf("%.2f", d))
		}
		ctx.currentSb.WriteString("] 0 d ")
	}

	if item.style.BorderWidth > 0 {
		ctx.currentSb.WriteString(fmt.Sprintf("%.2f w ", item.style.BorderWidth))
	}

	if item.shapeType == "line" {
		ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f RG %.2f %.2f m %.2f %.2f l S", r, g, b, item.x1, item.y1, item.x2, item.y2))
	} else if item.shapeType == "rect" {
		if hasBg {
			ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f rg %.2f %.2f %.2f %.2f re f ", rb, gb, bb, item.x1, item.y1, item.width, item.height))
		}
		if item.style.Border {
			ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f RG %.2f %.2f %.2f %.2f re S ", r, g, b, item.x1, item.y1, item.width, item.height))
		}
	} else if item.shapeType == "ellipse" {
		// Approximation of ellipse with 4 bezier curves
		kappa := 0.552284749831
		cx, cy := item.x1+item.width/2.0, item.y1+item.height/2.0
		rx, ry := item.width/2.0, item.height/2.0
		ox, oy := rx*kappa, ry*kappa

		op := "S"
		if hasBg && item.style.Border {
			ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f rg %.2f %.2f %.2f RG ", rb, gb, bb, r, g, b))
			op = "B"
		} else if hasBg {
			ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f rg ", rb, gb, bb))
			op = "f"
		} else {
			ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f RG ", r, g, b))
		}

		ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f m ", cx-rx, cy))
		ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f c ", cx-rx, cy+oy, cx-ox, cy+ry, cx, cy+ry))
		ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f c ", cx+ox, cy+ry, cx+rx, cy+oy, cx+rx, cy))
		ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f c ", cx+rx, cy-oy, cx+ox, cy-ry, cx, cy-ry))
		ctx.currentSb.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f c %s", cx-ox, cy-ry, cx-rx, cy-oy, cx-rx, cy, op))
	}
	ctx.currentSb.WriteString(" Q\n")
	ctx.currentSb.WriteString("EMC\n")
}

func (p *pageRenderer) calculateItemHeight(ctx *renderingContext, item *contentItem) float64 {
	if item.style.Absolute {
		return 0
	}
	m := p.getMargins()
	marginX := m.Left + m.Right

	if item.isImage {
		return item.height + 10
	} else if item.isTable {
		// Calculate total table height
		colWidths := make([]float64, item.cols)
		if len(item.colWidths) == item.cols {
			copy(colWidths, item.colWidths)
		} else {
			for c := range item.cols {
				colWidths[c] = (float64(ctx.w) - marginX) / float64(item.cols)
			}
		}
		totalH := 0.0
		for r := range item.rows {
			totalH += p.calculateRowHeight(ctx, item, r, colWidths)
		}
		return totalH + 10
	} else if item.isParagraph || item.isHeading {
		size := 12.0
		if item.style.Size > 0 {
			size = float64(item.style.Size)
		}
		fontName := p.getFontName(item.style)
		customWidths, unitsPerEm := p.getCustomWidths(item.style, ctx)
		maxWidth := float64(ctx.w) - marginX
		if p.pageSettings.Columns > 1 {
			totalGap := float64(p.pageSettings.Columns-1) * p.pageSettings.ColumnGap
			maxWidth = (float64(ctx.w) - marginX - totalGap) / float64(p.pageSettings.Columns)
		}

		firstLineMaxWidth := maxWidth - item.style.Indent
		otherLinesMaxWidth := maxWidth - item.style.Indent - item.style.Hanging
		if firstLineMaxWidth < 20 {
			firstLineMaxWidth = 20
		}
		if otherLinesMaxWidth < 20 {
			otherLinesMaxWidth = 20
		}

		lines := wrapText(item.text, []float64{firstLineMaxWidth, otherLinesMaxWidth}, size, fontName, customWidths, unitsPerEm)
		lineSpacing := 1.2
		if item.style.LineSpacing > 0 {
			lineSpacing = item.style.LineSpacing
		}
		h := float64(len(lines)) * size * lineSpacing
		return h + item.style.SpacingBefore + item.style.SpacingAfter
	} else if item.isRich {
		m := p.getMargins()
		marginX := m.Left + m.Right
		maxWidth := float64(ctx.w) - marginX
		if p.pageSettings.Columns > 1 {
			totalGap := float64(p.pageSettings.Columns-1) * p.pageSettings.ColumnGap
			maxWidth = (float64(ctx.w) - marginX - totalGap) / float64(p.pageSettings.Columns)
		}
		h, _ := p.renderRichParagraph(ctx, nil, item.spans, item.style, 0, 0, maxWidth, 12, true)
		return h
	} else if item.isList {
		m := p.getMargins()
		marginX := m.Left + m.Right
		size := 12.0
		if item.style.Size > 0 {
			size = float64(item.style.Size)
		}
		fontName := p.getFontName(item.style)
		customWidths, unitsPerEm := p.getCustomWidths(item.style, ctx)
		maxWidth := float64(ctx.w) - marginX
		if p.pageSettings.Columns > 1 {
			totalGap := float64(p.pageSettings.Columns-1) * p.pageSettings.ColumnGap
			maxWidth = (float64(ctx.w) - marginX - totalGap) / float64(p.pageSettings.Columns)
		}
		totalH := 0.0
		for _, listItem := range item.listItems {
			prefix := "• "
			if item.ordered {
				prefix = "1. "
			}
			prefixWidth := getTextWidth(prefix, size, fontName, customWidths, unitsPerEm)
			firstLineMaxWidth := maxWidth - item.style.Indent
			otherLinesMaxWidth := maxWidth - item.style.Indent - item.style.Hanging - prefixWidth
			if firstLineMaxWidth < 20 {
				firstLineMaxWidth = 20
			}
			if otherLinesMaxWidth < 20 {
				otherLinesMaxWidth = 20
			}
			lines := wrapText(prefix+listItem, []float64{firstLineMaxWidth, otherLinesMaxWidth}, size, fontName, customWidths, unitsPerEm)
			lineSpacing := 1.2
			if item.style.LineSpacing > 0 {
				lineSpacing = item.style.LineSpacing
			}
			totalH += float64(len(lines)) * size * lineSpacing
			totalH += item.style.SpacingBefore + item.style.SpacingAfter
		}
		return totalH
	}
	return 20.0
}

// Delegation helpers

func (p *pageRenderer) getExtGState(ctx *renderingContext, opacity float64) string {
	return (&renderer{p.state}).getExtGState(ctx, opacity)
}

func (p *pageRenderer) getFontName(style document.CellStyle) string {
	return (&renderer{p.state}).getFontName(style)
}

func (p *pageRenderer) getCustomWidths(style document.CellStyle, ctx *renderingContext) (map[rune]uint16, uint16) {
	return (&renderer{p.state}).getCustomWidths(style, ctx)
}

func (p *pageRenderer) renderParagraph(ctx *renderingContext, sb *strings.Builder, text string, style document.CellStyle, x, y, maxWidth, fontSize float64) float64 {
	return (&textRenderer{p.state}).renderParagraph(ctx, sb, text, style, x, y, maxWidth, fontSize)
}

func (p *pageRenderer) renderRichParagraph(ctx *renderingContext, sb *strings.Builder, spans []document.TextSpan, style document.CellStyle, x, y, maxWidth, fontSize float64, dryRun bool) (float64, float64) {
	return (&textRenderer{p.state}).renderRichParagraph(ctx, sb, spans, style, x, y, maxWidth, fontSize, dryRun)
}

func (p *pageRenderer) renderList(ctx *renderingContext, sb *strings.Builder, items []string, ordered bool, style document.CellStyle, x, y, maxWidth, fontSize float64) float64 {
	return (&textRenderer{p.state}).renderList(ctx, sb, items, ordered, style, x, y, maxWidth, fontSize)
}

func (p *pageRenderer) renderHeader(ctx *renderingContext, sb *strings.Builder) {
	(&textRenderer{p.state}).renderHeader(ctx, sb)
}

func (p *pageRenderer) renderFooter(ctx *renderingContext, sb *strings.Builder, pageNum, totalPages int) {
	(&textRenderer{p.state}).renderFooter(ctx, sb, pageNum, totalPages)
}

func (p *pageRenderer) renderWatermark(ctx *renderingContext, sb *strings.Builder) {
	(&textRenderer{p.state}).renderWatermark(ctx, sb)
}

func (p *pageRenderer) renderTOC(ctx *renderingContext, colWidth float64) {
	(&textRenderer{p.state}).renderTOC(ctx, colWidth)
}

func (p *pageRenderer) renderTable(ctx *renderingContext, item *contentItem, x, width float64) {
	(&tableRenderer{p.state}).renderTable(ctx, item, x, width)
}

func (p *pageRenderer) renderTableRow(ctx *renderingContext, item *contentItem, r int, y, startX float64, colWidths []float64, rowHeights []float64) {
	(&tableRenderer{p.state}).renderTableRow(ctx, item, r, y, startX, colWidths, rowHeights)
}

func (p *pageRenderer) calculateRowHeight(ctx *renderingContext, item *contentItem, row int, colWidths []float64) float64 {
	return (&tableRenderer{p.state}).calculateRowHeight(ctx, item, row, colWidths)
}
