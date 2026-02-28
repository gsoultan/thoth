package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

type renderer struct {
	*state
}

func (p *renderer) newRenderingContext() *renderingContext {
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

	// Use default margins if not specified
	m := p.getMargins()

	return &renderingContext{
		mgr:          objects.NewObjectManager(),
		imageRefs:    make(map[string]objects.Reference),
		imageNames:   make(map[string]string),
		fontRefs:     make(map[string]objects.Reference),
		fontNames:    make(map[string]string),
		importRefs:   make(map[string]objects.Reference),
		posY:         float64(h) - m.Top,
		w:            w,
		h:            h,
		customWidths: make(map[string]map[rune]uint16),
		unitsPerEm:   make(map[string]uint16),
		smaskRefs:    make(map[string]objects.Reference),
		extGStates:   make(map[float64]objects.Reference),
	}
}

func (p *renderer) collectImages(ctx *renderingContext, items []*contentItem) {
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

func (p *renderer) ensureImageInContext(ctx *renderingContext, path string) {
	if _, ok := ctx.imageRefs[path]; ok {
		return
	}
	data, alpha, imgW, imgH, format, err := getImageData(path)
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
		"Interpolate":      objects.Boolean(true),
	}

	if alpha != nil {
		// Create SMask
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write(alpha)
		zw.Close()

		smaskDict := objects.Dictionary{
			"Type":             objects.Name("XObject"),
			"Subtype":          objects.Name("Image"),
			"Width":            objects.Integer(imgW),
			"Height":           objects.Integer(imgH),
			"ColorSpace":       objects.Name("DeviceGray"),
			"BitsPerComponent": objects.Integer(8),
			"Filter":           objects.Name("FlateDecode"),
		}
		smaskStream := objects.Stream{Dict: smaskDict, Data: buf.Bytes()}
		smaskRef := ctx.mgr.AddObject(smaskStream)
		dict["SMask"] = smaskRef
		ctx.smaskRefs[path] = smaskRef
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

func (p *renderer) ensureFontInContext(ctx *renderingContext, name, path string) {
	if _, ok := ctx.fontRefs[name]; ok {
		return
	}
	metrics, err := parseTTF(path)
	if err != nil {
		fmt.Printf("Warning: failed to parse font %s: %v\n", name, err)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	// 1. FontFile2 (The TTF stream)
	fontFileDict := objects.Dictionary{
		"Length1": objects.Integer(len(data)),
	}
	fontFileStream := objects.Stream{Dict: fontFileDict, Data: data}
	fontFileRef := ctx.mgr.AddObject(fontFileStream)

	// 2. FontDescriptor
	ascent := float64(metrics.ascent) * 1000.0 / float64(metrics.unitsPerEm)
	descent := float64(metrics.descent) * 1000.0 / float64(metrics.unitsPerEm)

	descriptor := objects.Dictionary{
		"Type":        objects.Name("FontDescriptor"),
		"FontName":    objects.Name(metrics.name),
		"Flags":       objects.Integer(32), // Non-symbolic
		"FontBBox":    objects.Array{objects.Integer(-1000), objects.Integer(int(descent)), objects.Integer(3000), objects.Integer(int(ascent))},
		"ItalicAngle": objects.Float(metrics.italicAngle),
		"Ascent":      objects.Float(ascent),
		"Descent":     objects.Float(descent),
		"CapHeight":   objects.Float(ascent), // Simplified
		"StemV":       objects.Integer(80),
		"FontFile2":   fontFileRef,
	}
	descriptorRef := ctx.mgr.AddObject(descriptor)

	// 3. Font
	font := objects.Dictionary{
		"Type":           objects.Name("Font"),
		"Subtype":        objects.Name("TrueType"),
		"BaseFont":       objects.Name(metrics.name),
		"FontDescriptor": descriptorRef,
		"FirstChar":      objects.Integer(32),
		"LastChar":       objects.Integer(255),
		"Widths":         objects.Array{}, // We'd need to populate this
	}

	// Default widths
	widths := objects.Array{}
	for range 256 - 32 {
		widths = append(widths, objects.Integer(600)) // Default width
	}
	font["Widths"] = widths

	fontRef := ctx.mgr.AddObject(font)
	ctx.fontRefs[name] = fontRef
	ctx.fontNames[name] = fmt.Sprintf("CF%d", len(ctx.fontNames)+1)
	ctx.customWidths[name] = metrics.widths
	ctx.unitsPerEm[name] = metrics.unitsPerEm
}

func (p *renderer) ensureImportInContext(ctx *renderingContext, path string, pageNum int) {
	key := fmt.Sprintf("%s:%d", path, pageNum)
	if _, ok := ctx.importRefs[key]; ok {
		return
	}

	// For simplicity in this implementation, we will create a dummy Form XObject
	// Representing the imported page. In a full implementation, we'd parse the file.
	dict := objects.Dictionary{
		"Type":      objects.Name("XObject"),
		"Subtype":   objects.Name("Form"),
		"BBox":      objects.Array{objects.Integer(0), objects.Integer(0), objects.Integer(ctx.w), objects.Integer(ctx.h)},
		"Resources": objects.Dictionary{},
	}
	stream := objects.Stream{Dict: dict, Data: []byte("q 0.9 G 0 0 100 100 re f Q\n")} // Placeholder content
	ref := ctx.mgr.AddObject(stream)
	ctx.importRefs[key] = ref
}

func (p *renderer) getFont(style document.CellStyle, ctx *renderingContext) string {
	if style.Font != "" {
		if name, ok := ctx.fontNames[style.Font]; ok {
			return "/" + name
		}
	}
	if style.Bold && style.Italic {
		return "/F4"
	} else if style.Bold {
		return "/F2"
	} else if style.Italic {
		return "/F3"
	}
	return "/F1"
}

func (p *renderer) getFontName(style document.CellStyle) string {
	if style.Font != "" {
		return style.Font
	}
	if style.Bold && style.Italic {
		return "Helvetica-BoldOblique"
	} else if style.Bold {
		return "Helvetica-Bold"
	} else if style.Italic {
		return "Helvetica-Oblique"
	}
	return "Helvetica"
}

func (p *renderer) getCustomWidths(style document.CellStyle, ctx *renderingContext) (map[rune]uint16, uint16) {
	if style.Font != "" {
		return ctx.customWidths[style.Font], ctx.unitsPerEm[style.Font]
	}
	return nil, 1000
}

func (p *renderer) getExtGState(ctx *renderingContext, opacity float64) string {
	if opacity <= 0 || opacity >= 1.0 {
		return ""
	}
	// Round to 2 decimal places for consistency
	op := float64(int(opacity*100)) / 100.0
	if ref, ok := ctx.extGStates[op]; ok {
		return fmt.Sprintf("GS%d", ref.Number)
	}

	dict := objects.Dictionary{
		"Type": objects.Name("ExtGState"),
		"ca":   objects.Float(op), // Non-stroking alpha
		"CA":   objects.Float(op), // Stroking alpha
	}
	ref := ctx.mgr.AddObject(dict)
	ctx.extGStates[op] = ref
	return fmt.Sprintf("GS%d", ref.Number)
}
