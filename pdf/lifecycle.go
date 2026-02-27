package pdf

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
	"github.com/gsoultan/thoth/pdf/internal/parser"
)

// lifecycle handles document lifecycle operations.
type lifecycle struct{ *state }

type bookmark struct {
	title string
	page  int
	posY  float64
}

type link struct {
	rect []float64
	url  string
}

type pageInfo struct {
	contentItems []contentItem // Actually we already have contentItems in state
	// But we need to know which items went to which page?
	// No, the current approach of rendering into a strings.Builder is easier
	// to determine where things fit.
	sb    string
	links []link
	posY  float64 // Not really needed if we already rendered
}

type renderingContext struct {
	mgr          *objects.ObjectManager
	imageRefs    map[string]objects.Reference
	imageNames   map[string]string
	pageRefs     []objects.Reference
	pages        []pageInfo
	currentSb    strings.Builder
	currentLinks []link
	bookmarks    []bookmark
	posY         float64
	w, h         int
}

// Open loads a document from a reader.
func (p *lifecycle) Open(ctx context.Context, reader io.Reader) error {
	l := parser.NewLexer(reader)
	pr := parser.NewParser(l)

	for {
		obj, err := pr.ParseObject()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			// Skip single non-objects (e.g., trailer keyword, %%EOF)
			continue
		}
		if obj == nil {
			// Check if we hit EOF via nil obj
			break
		}
		p.objects = append(p.objects, obj)
	}

	return nil
}

// Save writes the document to a writer.
func (p *lifecycle) Save(ctx context.Context, writer io.Writer) error {
	if len(p.objects) > 0 {
		return p.saveModified(ctx, writer)
	}

	renderCtx := p.newRenderingContext()

	// Pre-process images
	p.collectImages(renderCtx, p.contentItems)
	p.collectImages(renderCtx, p.header)
	p.collectImages(renderCtx, p.footer)

	// Render content
	p.renderContent(renderCtx)

	p.finishPage(renderCtx)

	p.finalizePages(renderCtx)

	return p.writePDF(renderCtx, writer)
}

func (p *lifecycle) newRenderingContext() *renderingContext {
	w, h := 612, 792 // Letter default (points)
	switch p.pageSettings.PaperType {
	case document.PaperA4:
		w, h = 595, 842
	case document.PaperLetter:
		w, h = 612, 792
	}
	if p.pageSettings.Orientation == document.OrientationLandscape {
		w, h = h, w
	}

	return &renderingContext{
		mgr:        objects.NewObjectManager(),
		imageRefs:  make(map[string]objects.Reference),
		imageNames: make(map[string]string),
		posY:       float64(h - 80),
		w:          w,
		h:          h,
	}
}

func (p *lifecycle) collectImages(ctx *renderingContext, items []contentItem) {
	for _, item := range items {
		if item.isImage && item.path != "" {
			p.ensureImageInContext(ctx, item.path)
		}
		if item.isTable {
			for r := range item.rows {
				for c := range item.cols {
					for _, ci := range item.cells[r][c] {
						if ci.isImage && ci.path != "" {
							p.ensureImageInContext(ctx, ci.path)
						}
					}
				}
			}
		}
	}
}

func (p *lifecycle) ensureImageInContext(ctx *renderingContext, path string) {
	if _, ok := ctx.imageRefs[path]; ok {
		return
	}
	data, imgW, imgH, format, err := getImageData(path)
	if err != nil {
		fmt.Printf("Warning: failed to load image %s: %v\n", path, err)
		return
	}

	dict := objects.Dictionary{
		"Type":             objects.Name("XObject"),
		"Subtype":          objects.Name("Image"),
		"Width":            objects.Integer(imgW),
		"Height":           objects.Integer(imgH),
		"ColorSpace":       objects.Name("DeviceRGB"),
		"BitsPerComponent": objects.Integer(8),
		"Interpolate":      objects.Name("true"),
	}

	if format == "jpeg" {
		dict["Filter"] = objects.Name("DCTDecode")
	} else {
		// Use FlateDecode for other formats (which we converted to raw RGB)
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write(data)
		zw.Close()
		data = buf.Bytes()
		dict["Filter"] = objects.Name("FlateDecode")
	}

	imgObj := objects.Stream{Dict: dict, Data: data}
	ref := ctx.mgr.AddObject(imgObj)
	ctx.imageRefs[path] = ref
	ctx.imageNames[path] = fmt.Sprintf("Img%d", len(ctx.imageNames)+1)
}

func (p *lifecycle) getFont(style document.CellStyle) string {
	if style.Bold && style.Italic {
		return "/F4"
	} else if style.Bold {
		return "/F2"
	} else if style.Italic {
		return "/F3"
	}
	return "/F1"
}

func (p *lifecycle) renderParagraph(ctx *renderingContext, sb *strings.Builder, text string, style document.CellStyle, x, y, maxWidth, fontSize float64) float64 {
	if text == "" {
		return 0
	}

	// Add bookmark for headings
	if style.Name != "" && (strings.HasPrefix(style.Name, "Heading") || style.Name == "Title") {
		// Only add bookmarks for contentItems, not headers/footers
		// Check if we are rendering into currentSb (main content)
		if sb == &ctx.currentSb {
			ctx.bookmarks = append(ctx.bookmarks, bookmark{
				title: text,
				page:  len(ctx.pageRefs),
				posY:  y,
			})
		}
	}

	font := p.getFont(style)
	size := fontSize
	if style.Size > 0 {
		size = float64(style.Size)
	}
	lines := wrapText(text, maxWidth, size)
	currY := y
	for _, line := range lines {
		r, g, b := 0.0, 0.0, 0.0
		if style.Color != "" {
			r, g, b = hexToRGB(style.Color)
		}
		offsetX := 0.0
		if style.Horizontal == "center" {
			offsetX = (maxWidth - getTextWidth(line, size)) / 2.0
		} else if style.Horizontal == "right" {
			offsetX = maxWidth - getTextWidth(line, size)
		}

		if style.Link != "" && sb == &ctx.currentSb {
			ctx.currentLinks = append(ctx.currentLinks, link{
				rect: []float64{x + offsetX, currY - size, x + offsetX + getTextWidth(line, size), currY},
				url:  style.Link,
			})
		}

		// Use y-size as baseline
		sb.WriteString(fmt.Sprintf("BT %.2f %.2f %.2f rg %s %.2f Tf %.2f %.2f Td (%s) Tj ET\n", r, g, b, font, size, x+offsetX, currY-size, escapePDF(line)))
		currY -= (size * 1.2)
	}
	return y - currY
}

func (p *lifecycle) renderHeader(ctx *renderingContext, sb *strings.Builder) {
	currY := float64(ctx.h - 40)
	for _, item := range p.header {
		if item.isParagraph {
			currY -= p.renderParagraph(ctx, sb, item.text, item.style, 50, currY, float64(ctx.w-100), 10)
		}
	}
}

func (p *lifecycle) renderFooter(ctx *renderingContext, sb *strings.Builder, pageNum int, totalPages int) {
	currY := 40.0
	for _, item := range p.footer {
		if item.isParagraph {
			text := strings.ReplaceAll(item.text, "{n}", fmt.Sprintf("%d", pageNum))
			text = strings.ReplaceAll(text, "{nb}", fmt.Sprintf("%d", totalPages))
			p.renderParagraph(ctx, sb, text, item.style, 50, currY, float64(ctx.w-100), 10)
			currY -= 12
		}
	}
}

func (p *lifecycle) finishPage(ctx *renderingContext) {
	if ctx.currentSb.Len() == 0 && len(ctx.pages) > 0 {
		return
	}

	ctx.pages = append(ctx.pages, pageInfo{
		sb:    ctx.currentSb.String(),
		links: ctx.currentLinks,
	})

	ctx.currentSb.Reset()
	ctx.currentLinks = nil
	ctx.posY = float64(ctx.h - 80)
}

func (p *lifecycle) finalizePages(ctx *renderingContext) {
	total := len(ctx.pages)
	for i := range total {
		pageNum := i + 1
		var sb strings.Builder

		// Render header and footer into a new builder
		p.renderHeader(ctx, &sb)
		p.renderFooter(ctx, &sb, pageNum, total)

		// Append original page contentItems
		sb.WriteString(ctx.pages[i].sb)

		// Stream compression
		data := []byte(sb.String())
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

		// Annotations (Links)
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
				"A": objects.Dictionary{
					"Type": objects.Name("Action"),
					"S":    objects.Name("URI"),
					"URI":  objects.PDFString(l.url),
				},
			}
			annots = append(annots, annot)
		}

		page := objects.Dictionary{
			"Type":     objects.Name("Page"),
			"MediaBox": objects.Array{objects.Integer(0), objects.Integer(0), objects.Integer(ctx.w), objects.Integer(ctx.h)},
			"Resources": objects.Dictionary{
				"ProcSet": objects.Array{objects.Name("PDF"), objects.Name("Text"), objects.Name("ImageB"), objects.Name("ImageC"), objects.Name("ImageI")},
				"Font": objects.Dictionary{
					"F1": objects.Dictionary{"Type": objects.Name("Font"), "Subtype": objects.Name("Type1"), "BaseFont": objects.Name("Helvetica")},
					"F2": objects.Dictionary{"Type": objects.Name("Font"), "Subtype": objects.Name("Type1"), "BaseFont": objects.Name("Helvetica-Bold")},
					"F3": objects.Dictionary{"Type": objects.Name("Font"), "Subtype": objects.Name("Type1"), "BaseFont": objects.Name("Helvetica-Oblique")},
					"F4": objects.Dictionary{"Type": objects.Name("Font"), "Subtype": objects.Name("Type1"), "BaseFont": objects.Name("Helvetica-BoldOblique")},
				},
				"XObject": xobjects,
			},
			"Contents": streamRef,
		}

		if len(annots) > 0 {
			page["Annots"] = annots
		}

		pageRef := ctx.mgr.AddObject(page)
		ctx.pageRefs = append(ctx.pageRefs, pageRef)
	}
}

func (p *lifecycle) renderContent(ctx *renderingContext) {
	for _, item := range p.contentItems {
		if item.isPageBreak {
			p.finishPage(ctx)
			continue
		}
		if item.isShape {
			p.renderShape(ctx, item)
			continue
		}

		itemHeight := p.calculateItemHeight(ctx, item)
		if ctx.posY-itemHeight < 60 {
			p.finishPage(ctx)
		}

		if item.isParagraph {
			ctx.posY -= p.renderParagraph(ctx, &ctx.currentSb, item.text, item.style, 50, ctx.posY, float64(ctx.w-100), 12)
			ctx.posY -= 10
		} else if item.isImage {
			imgName := ctx.imageNames[item.path]
			if imgName != "" {
				x := 50.0
				if item.style.Horizontal == "center" {
					x = (float64(ctx.w) - item.width) / 2.0
				} else if item.style.Horizontal == "right" {
					x = float64(ctx.w) - item.width - 50.0
				}
				ctx.currentSb.WriteString(fmt.Sprintf("q %.2f 0 0 %.2f %.2f %.2f cm /%s Do Q\n", item.width, item.height, x, ctx.posY-item.height, imgName))
			}
			ctx.posY -= (item.height + 20)
		} else if item.isTable {
			p.renderTable(ctx, item)
		}
	}
}

func (p *lifecycle) renderShape(ctx *renderingContext, item contentItem) {
	r, g, b := 0.0, 0.0, 0.0
	if item.style.Color != "" {
		r, g, b = hexToRGB(item.style.Color)
	}
	if item.shapeType == "line" {
		ctx.currentSb.WriteString(fmt.Sprintf("q %.2f %.2f %.2f RG %.2f %.2f m %.2f %.2f l S Q\n", r, g, b, item.x1, item.y1, item.x2, item.y2))
	} else if item.shapeType == "rect" {
		if item.style.Background != "" {
			rb, gb, bb := hexToRGB(item.style.Background)
			ctx.currentSb.WriteString(fmt.Sprintf("q %.2f %.2f %.2f rg %.2f %.2f %.2f %.2f re f Q\n", rb, gb, bb, item.x1, item.y1, item.width, item.height))
		}
		if item.style.Border {
			ctx.currentSb.WriteString(fmt.Sprintf("q %.2f %.2f %.2f RG %.2f %.2f %.2f %.2f re S Q\n", r, g, b, item.x1, item.y1, item.width, item.height))
		}
	}
}

func (p *lifecycle) calculateItemHeight(ctx *renderingContext, item contentItem) float64 {
	if item.isImage {
		return item.height + 10
	} else if item.isTable {
		// Calculate total table height
		colWidths := make([]float64, item.cols)
		for c := range item.cols {
			colWidths[c] = (float64(ctx.w) - 100.0) / float64(item.cols)
		}
		totalH := 0.0
		for r := range item.rows {
			totalH += p.calculateRowHeight(ctx, item, r, colWidths)
		}
		return totalH + 10
	} else if item.isParagraph {
		size := 12.0
		if item.style.Size > 0 {
			size = float64(item.style.Size)
		}
		lines := wrapText(item.text, float64(ctx.w-100), size)
		return float64(len(lines)) * size * 1.2
	}
	return 20.0
}

func (p *lifecycle) calculateRowHeight(ctx *renderingContext, item contentItem, row int, colWidths []float64) float64 {
	maxH := 20.0
	for c := range item.cols {
		if len(item.cells[row][c]) > 0 && item.cells[row][c][0].hidden {
			continue
		}
		if len(item.cells[row][c]) > 0 && item.cells[row][c][0].rowSpan > 1 {
			continue
		}
		cW := colWidths[c]
		if len(item.cells[row][c]) > 0 && item.cells[row][c][0].colSpan > 1 {
			for i := 1; i < item.cells[row][c][0].colSpan && c+i < item.cols; i++ {
				cW += colWidths[c+i]
			}
		}
		cellH := p.calculateCellHeight(item.cells[row][c], cW)
		if cellH > maxH {
			maxH = cellH
		}
	}
	return maxH
}

func (p *lifecycle) calculateCellHeight(cells []cellItem, width float64) float64 {
	cellH := 0.0
	for _, ci := range cells {
		if ci.isImage {
			cellH += ci.height
		} else if ci.text != "" {
			size := 10.0
			if ci.style.Size > 0 {
				size = float64(ci.style.Size)
			}
			lines := wrapText(ci.text, width-4.0, size)
			cellH += float64(len(lines)) * size * 1.2
		}
		cellH += 2.0
	}
	return cellH + 4.0 // padding
}

func (p *lifecycle) renderTable(ctx *renderingContext, item contentItem) {
	colWidths := make([]float64, item.cols)
	for c := range item.cols {
		colWidths[c] = (float64(ctx.w) - 100.0) / float64(item.cols)
	}

	rowHeights := make([]float64, item.rows)
	for r := range item.rows {
		rowHeights[r] = p.calculateRowHeight(ctx, item, r, colWidths)
	}

	// Adjust row heights for rowSpan
	for r := range item.rows {
		for c := range item.cols {
			if len(item.cells[r][c]) > 0 && item.cells[r][c][0].rowSpan > 1 {
				rs := item.cells[r][c][0].rowSpan
				cW := colWidths[c]
				if item.cells[r][c][0].colSpan > 1 {
					for i := 1; i < item.cells[r][c][0].colSpan && c+i < item.cols; i++ {
						cW += colWidths[c+i]
					}
				}
				cellH := p.calculateCellHeight(item.cells[r][c], cW)

				currentTotalH := 0.0
				for i := 0; i < rs && r+i < item.rows; i++ {
					currentTotalH += rowHeights[r+i]
				}

				if cellH > currentTotalH {
					extra := cellH - currentTotalH
					lastRow := r + rs - 1
					if lastRow >= item.rows {
						lastRow = item.rows - 1
					}
					rowHeights[lastRow] += extra
				}
			}
		}
	}

	for r := range item.rows {
		rh := rowHeights[r]
		if ctx.posY-rh < 60 {
			p.finishPage(ctx)
			if r > 0 { // Repeat header
				p.renderTableRow(ctx, item, 0, ctx.posY, colWidths, rowHeights)
				ctx.posY -= rowHeights[0]
			}
		}
		p.renderTableRow(ctx, item, r, ctx.posY, colWidths, rowHeights)
		ctx.posY -= rh
	}
	ctx.posY -= 10
}

func (p *lifecycle) renderTableRow(ctx *renderingContext, item contentItem, r int, y float64, colWidths []float64, rowHeights []float64) {
	rowHeight := rowHeights[r]
	for c := range item.cols {
		if len(item.cells[r][c]) > 0 && item.cells[r][c][0].hidden {
			continue
		}
		cellX := 50.0
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
				for i := 0; i < item.cells[r][c][0].rowSpan && r+i < item.rows; i++ {
					cH += rowHeights[r+i]
				}
			}
		}

		// Background
		if len(item.cells[r][c]) > 0 && item.cells[r][c][0].style.Background != "" {
			rb, gb, bb := hexToRGB(item.cells[r][c][0].style.Background)
			ctx.currentSb.WriteString(fmt.Sprintf("q %.2f %.2f %.2f rg %.2f %.2f %.2f %.2f re f Q\n", rb, gb, bb, cellX, y-cH, cW, cH))
		}
		// Border
		if len(item.cells[r][c]) > 0 && item.cells[r][c][0].style.Border {
			ctx.currentSb.WriteString(fmt.Sprintf("q 0.5 w 0 G %.2f %.2f %.2f %.2f re S Q\n", cellX, y-cH, cW, cH))
		}

		currY := y - 2.0
		for _, ci := range item.cells[r][c] {
			if ci.isImage {
				ix := cellX + 2.0
				if ci.style.Horizontal == "center" {
					ix = cellX + (cW-ci.width)/2.0
				} else if ci.style.Horizontal == "right" {
					ix = cellX + cW - ci.width - 2.0
				}
				imgName := ctx.imageNames[ci.path]
				if imgName != "" {
					ctx.currentSb.WriteString(fmt.Sprintf("q %.2f 0 0 %.2f %.2f %.2f cm /%s Do Q\n", ci.width, ci.height, ix, currY-ci.height, imgName))
				}
				currY -= ci.height + 2.0
			} else if ci.text != "" {
				size := 10.0
				if ci.style.Size > 0 {
					size = float64(ci.style.Size)
				}
				currY -= p.renderParagraph(ctx, &ctx.currentSb, ci.text, ci.style, cellX+2.0, currY, cW-4.0, size)
				currY -= 2.0
			}
		}
	}
}

func (p *lifecycle) writePDF(ctx *renderingContext, writer io.Writer) error {
	// Pages
	pages := objects.Dictionary{"Type": objects.Name("Pages"), "Count": objects.Integer(len(ctx.pageRefs)), "Kids": objects.Array{}}
	pagesRef := ctx.mgr.AddObject(pages)
	for _, ref := range ctx.pageRefs {
		pages["Kids"] = append(pages["Kids"].(objects.Array), ref)
		for i := range ctx.mgr.Objects {
			if ctx.mgr.Objects[i].Number == ref.Number {
				if d, ok := ctx.mgr.Objects[i].Data.(objects.Dictionary); ok {
					d["Parent"] = pagesRef
				}
			}
		}
	}

	// Outlines (Bookmarks)
	var outlinesRef *objects.Reference
	if len(ctx.bookmarks) > 0 {
		outlinesDict := objects.Dictionary{"Type": objects.Name("Outlines"), "Count": objects.Integer(len(ctx.bookmarks))}
		outRef := ctx.mgr.AddObject(outlinesDict)
		outlinesRef = &outRef

		entryRefs := make([]objects.Reference, len(ctx.bookmarks))
		for i := range ctx.bookmarks {
			// Ensure page index is valid
			pageIdx := ctx.bookmarks[i].page
			if pageIdx >= len(ctx.pageRefs) {
				pageIdx = len(ctx.pageRefs) - 1
			}
			if pageIdx < 0 {
				pageIdx = 0
			}

			entryRefs[i] = ctx.mgr.AddObject(objects.Dictionary{
				"Title":  objects.PDFString(ctx.bookmarks[i].title),
				"Parent": *outlinesRef,
				"Dest":   objects.Array{ctx.pageRefs[pageIdx], objects.Name("XYZ"), objects.Integer(0), objects.Integer(int(ctx.bookmarks[i].posY)), objects.Integer(0)},
			})
		}

		for i := range entryRefs {
			// Find the object in manager
			for j := range ctx.mgr.Objects {
				if ctx.mgr.Objects[j].Number == entryRefs[i].Number {
					dict := ctx.mgr.Objects[j].Data.(objects.Dictionary)
					if i > 0 {
						dict["Prev"] = entryRefs[i-1]
					}
					if i < len(entryRefs)-1 {
						dict["Next"] = entryRefs[i+1]
					}
					ctx.mgr.Objects[j].Data = dict
					break
				}
			}
		}

		outlinesDict["First"] = entryRefs[0]
		outlinesDict["Last"] = entryRefs[len(entryRefs)-1]
		// Update outlinesDict in manager
		for j := range ctx.mgr.Objects {
			if ctx.mgr.Objects[j].Number == outRef.Number {
				ctx.mgr.Objects[j].Data = outlinesDict
				break
			}
		}
	}

	catalog := objects.Dictionary{"Type": objects.Name("Catalog"), "Pages": pagesRef}
	if outlinesRef != nil {
		catalog["Outlines"] = *outlinesRef
	}
	catRef := ctx.mgr.AddObject(catalog)
	var infoRef *objects.Reference
	if p.info != nil {
		ref := ctx.mgr.AddObject(p.info)
		infoRef = &ref
	}
	return ctx.mgr.Write(writer, catRef, infoRef)
}

// Close releases any resources used by the document.
func (p *lifecycle) Close() error {
	return nil
}
