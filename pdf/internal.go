package pdf

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

var sbPool = sync.Pool{
	New: func() any {
		return new(strings.Builder)
	},
}

func getSB() *strings.Builder {
	sb := sbPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

func putSB(sb *strings.Builder) {
	sb.Reset()
	sbPool.Put(sb)
}

func formatPDFDate(t time.Time) string {
	return fmt.Sprintf("D:%s", t.Format("20060102150405-07'00'"))
}

func extractTextFromStream(data []byte) string {
	sb := getSB()
	defer putSB(sb)
	s := string(data)

	// Very naive contentItems stream parser
	for line := range strings.SplitSeq(s, "\n") {
		if strings.Contains(line, "Tj") {
			start := strings.Index(line, "(")
			end := strings.LastIndex(line, ")")
			if start != -1 && end != -1 && end > start {
				sb.WriteString(line[start+1 : end])
				sb.WriteString(" ")
			}
		} else if strings.Contains(line, "TJ") {
			start := strings.Index(line, "[")
			end := strings.LastIndex(line, "]")
			if start != -1 && end != -1 && end > start {
				content := line[start+1 : end]
				for {
					sIdx := strings.Index(content, "(")
					eIdx := strings.Index(content, ")")
					if sIdx == -1 || eIdx == -1 {
						break
					}
					sb.WriteString(content[sIdx+1 : eIdx])
					content = content[eIdx+1:]
				}
				sb.WriteString(" ")
			}
		}
	}
	return sb.String()
}

func (p *state) saveModified(ctx context.Context, writer io.Writer) error {
	fmt.Fprintf(writer, "%%PDF-2.0\n%%\xe2\xe3\xcf\xd3\n")

	type objOffset struct {
		offset int64
		gen    int
	}
	offsets := make(map[int]objOffset)
	currOffset := int64(15) // %PDF-1.4\n%âãÏÓ\n

	for _, obj := range p.objects {
		if ind, ok := obj.(*objects.IndirectObject); ok {
			offsets[ind.Number] = objOffset{offset: currOffset, gen: ind.Generation}
			n64, _ := ind.WriteTo(writer)
			currOffset += n64
			n, _ := fmt.Fprint(writer, "\n")
			currOffset += int64(n)
		}
	}

	xrefStart := currOffset
	maxID := 0
	if len(offsets) > 0 {
		maxID = slices.Max(slices.Collect(maps.Keys(offsets)))
	}

	fmt.Fprintf(writer, "xref\n0 %d\n0000000000 65535 f\r\n", maxID+1)
	for i := range maxID {
		objID := i + 1
		if oo, ok := offsets[objID]; ok {
			fmt.Fprintf(writer, "%010d %05d n\r\n", oo.offset, oo.gen)
		} else {
			fmt.Fprintf(writer, "0000000000 65535 f\r\n")
		}
	}

	trailer := objects.Dictionary{
		"Size": objects.Integer(maxID + 1),
	}
	// Find root and info
	for _, obj := range p.objects {
		if ind, ok := obj.(*objects.IndirectObject); ok {
			if dict, ok := ind.Data.(objects.Dictionary); ok {
				if dict["Type"] == objects.Name("Catalog") {
					trailer["Root"] = objects.Reference{Number: ind.Number, Generation: ind.Generation}
				}
				if dict["Type"] == objects.Name("Info") {
					trailer["Info"] = objects.Reference{Number: ind.Number, Generation: ind.Generation}
				}
			}
		}
	}

	fmt.Fprintf(writer, "trailer\n")
	_, _ = trailer.WriteTo(writer)
	fmt.Fprintf(writer, "\nstartxref\n%d\n%%%%EOF\n", xrefStart)
	return nil
}

func getImageData(path string) ([]byte, []byte, int, int, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, 0, 0, "", err
	}
	defer file.Close()

	cfg, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, nil, 0, 0, "", err
	}

	if format == "jpeg" {
		file.Close()
		data, err := os.ReadFile(path)
		return data, nil, cfg.Width, cfg.Height, "jpeg", err
	}

	file.Seek(0, 0)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, nil, 0, 0, "", err
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	rgb := make([]byte, w*h*3)
	alpha := make([]byte, w*h)
	idx := 0
	aIdx := 0
	hasAlpha := false
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgb[idx] = byte(r >> 8)
			rgb[idx+1] = byte(g >> 8)
			rgb[idx+2] = byte(b >> 8)
			idx += 3

			alphaVal := byte(a >> 8)
			alpha[aIdx] = alphaVal
			aIdx++
			if alphaVal < 255 {
				hasAlpha = true
			}
		}
	}
	if !hasAlpha {
		alpha = nil
	}
	return rgb, alpha, w, h, format, nil
}

var standardFontWidths = map[string]map[rune]int{
	"Helvetica": {
		' ': 278, '!': 278, '"': 355, '#': 556, '$': 556, '%': 889, '&': 667, '\'': 191,
		'(': 333, ')': 333, '*': 389, '+': 584, ',': 278, '-': 333, '.': 278, '/': 278,
		'0': 556, '1': 556, '2': 556, '3': 556, '4': 556, '5': 556, '6': 556, '7': 556,
		'8': 556, '9': 556, ':': 278, ';': 278, '<': 584, '=': 584, '>': 584, '?': 556,
		'@': 1015, 'A': 667, 'B': 667, 'C': 722, 'D': 722, 'E': 667, 'F': 611, 'G': 778,
		'H': 722, 'I': 278, 'J': 500, 'K': 667, 'L': 556, 'M': 833, 'N': 722, 'O': 778,
		'P': 667, 'Q': 778, 'R': 722, 'S': 667, 'T': 611, 'U': 722, 'V': 667, 'W': 944,
		'X': 667, 'Y': 667, 'Z': 611, '[': 278, '\\': 278, ']': 278, '^': 469, '_': 556,
		'`': 333, 'a': 556, 'b': 556, 'c': 500, 'd': 556, 'e': 556, 'f': 278, 'g': 556,
		'h': 556, 'i': 222, 'j': 222, 'k': 500, 'l': 222, 'm': 833, 'n': 556, 'o': 556,
		'p': 556, 'q': 556, 'r': 333, 's': 500, 't': 278, 'u': 556, 'v': 500, 'w': 722,
		'x': 500, 'y': 500, 'z': 500, '{': 334, '|': 260, '}': 334, '~': 584,
		0x95: 350,
	},
	"Helvetica-Bold": {
		' ': 278, '!': 333, '"': 474, '#': 556, '$': 556, '%': 889, '&': 722, '\'': 238,
		'(': 333, ')': 333, '*': 389, '+': 584, ',': 333, '-': 333, '.': 333, '/': 278,
		'0': 556, '1': 556, '2': 556, '3': 556, '4': 556, '5': 556, '6': 556, '7': 556,
		'8': 556, '9': 556, ':': 333, ';': 333, '<': 584, '=': 584, '>': 584, '?': 611,
		'@': 975, 'A': 722, 'B': 722, 'C': 722, 'D': 722, 'E': 667, 'F': 611, 'G': 778,
		'H': 722, 'I': 278, 'J': 556, 'K': 722, 'L': 611, 'M': 833, 'N': 722, 'O': 778,
		'P': 667, 'Q': 778, 'R': 722, 'S': 667, 'T': 611, 'U': 722, 'V': 722, 'W': 944,
		'X': 722, 'Y': 722, 'Z': 611, '[': 333, '\\': 278, ']': 333, '^': 584, '_': 556,
		'`': 333, 'a': 556, 'b': 611, 'c': 556, 'd': 611, 'e': 556, 'f': 333, 'g': 611,
		'h': 611, 'i': 278, 'j': 278, 'k': 556, 'l': 278, 'm': 889, 'n': 611, 'o': 611,
		'p': 611, 'q': 611, 'r': 389, 's': 556, 't': 333, 'u': 611, 'v': 556, 'w': 778,
		'x': 556, 'y': 556, 'z': 500, '{': 389, '|': 280, '}': 389, '~': 584,
		0x95: 350,
	},
	"Times-Roman": {
		' ': 250, '!': 333, '"': 408, '#': 500, '$': 500, '%': 833, '&': 778, '\'': 180,
		'(': 333, ')': 333, '*': 500, '+': 564, ',': 250, '-': 333, '.': 250, '/': 278,
		'0': 500, '1': 500, '2': 500, '3': 500, '4': 500, '5': 500, '6': 500, '7': 500,
		'8': 500, '9': 500, ':': 278, ';': 278, '<': 564, '=': 564, '>': 564, '?': 444,
		'@': 921, 'A': 722, 'B': 667, 'C': 667, 'D': 722, 'E': 611, 'F': 556, 'G': 722,
		'H': 722, 'I': 333, 'J': 389, 'K': 722, 'L': 611, 'M': 889, 'N': 722, 'O': 722,
		'P': 556, 'Q': 722, 'R': 667, 'S': 556, 'T': 611, 'U': 722, 'V': 722, 'W': 944,
		'X': 722, 'Y': 722, 'Z': 611, '[': 333, '\\': 278, ']': 333, '^': 469, '_': 500,
		'`': 333, 'a': 444, 'b': 500, 'c': 444, 'd': 500, 'e': 444, 'f': 333, 'g': 500,
		'h': 500, 'i': 278, 'j': 278, 'k': 500, 'l': 278, 'm': 778, 'n': 500, 'o': 500,
		'p': 500, 'q': 500, 'r': 333, 's': 389, 't': 278, 'u': 500, 'v': 500, 'w': 722,
		'x': 500, 'y': 500, 'z': 444, '{': 480, '|': 200, '}': 480, '~': 541,
		0x95: 350,
	},
}

var fontAliases = map[string]string{
	"Helvetica-Oblique":     "Helvetica",
	"Helvetica-BoldOblique": "Helvetica-Bold",
}

func getTextWidth(text string, fontSize float64, fontName string, customWidths map[rune]uint16, unitsPerEm uint16) float64 {
	if customWidths != nil {
		totalWidth := 0
		for _, r := range text {
			if w, ok := customWidths[r]; ok {
				totalWidth += int(w)
			} else {
				totalWidth += int(unitsPerEm) / 2 // fallback to 0.5 em
			}
		}
		if unitsPerEm == 0 {
			unitsPerEm = 1000
		}
		return (float64(totalWidth) / float64(unitsPerEm)) * fontSize
	}

	if alias, ok := fontAliases[fontName]; ok {
		fontName = alias
	}
	widths, ok := standardFontWidths[fontName]
	if !ok {
		// Fallback to Helvetica
		widths = standardFontWidths["Helvetica"]
	}

	totalWidth := 0
	for _, r := range text {
		if w, ok := widths[r]; ok {
			totalWidth += w
		} else {
			totalWidth += 500 // fallback
		}
	}
	return (float64(totalWidth) / 1000.0) * fontSize
}

func wrapText(text string, maxWidths []float64, fontSize float64, fontName string, customWidths map[rune]uint16, unitsPerEm uint16) []string {
	if len(maxWidths) == 0 {
		return []string{text}
	}

	var lines []string
	paragraphs := strings.Split(text, "\n")
	for _, p := range paragraphs {
		words := strings.Fields(p)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		lineIdx := 0
		currentMaxWidth := maxWidths[0]
		currentLine := words[0]

		for i := 1; i < len(words); i++ {
			word := words[i]
			if getTextWidth(currentLine+" "+word, fontSize, fontName, customWidths, unitsPerEm) <= currentMaxWidth {
				currentLine += " " + word
			} else {
				lines = append(lines, currentLine)
				currentLine = word
				lineIdx++
				if lineIdx < len(maxWidths) {
					currentMaxWidth = maxWidths[lineIdx]
				} else {
					currentMaxWidth = maxWidths[len(maxWidths)-1]
				}
			}
		}
		lines = append(lines, currentLine)
	}
	return lines
}

func hexToRGB(hex string) (float64, float64, float64) {
	if hex == "" {
		return 0, 0, 0
	}
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}
	var r, g, b int
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return 0, 0, 0
	}
	return float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0
}

func escapePDF(s string) string {
	return objects.EscapeString(s)
}

func generateXMP(meta document.Metadata) string {
	sb := getSB()
	defer putSB(sb)

	now := time.Now().Format("2006-01-02T15:04:05-07:00")

	sb.WriteString(`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?>`)
	sb.WriteString(`<x:xmpmeta xmlns:x="adobe:ns:meta/">`)
	sb.WriteString(`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">`)

	// Dublin Core
	sb.WriteString(`<rdf:Description rdf:about="" xmlns:dc="http://purl.org/dc/elements/1.1/">`)
	sb.WriteString(`<dc:format>application/pdf</dc:format>`)
	if meta.Title != "" {
		sb.WriteString(fmt.Sprintf(`<dc:title><rdf:Alt><rdf:li xml:lang="x-default">%s</rdf:li></rdf:Alt></dc:title>`, meta.Title))
	}
	if meta.Author != "" {
		sb.WriteString(fmt.Sprintf(`<dc:creator><rdf:Seq><rdf:li>%s</rdf:li></rdf:Seq></dc:creator>`, meta.Author))
	}
	if meta.Description != "" {
		sb.WriteString(fmt.Sprintf(`<dc:description><rdf:Alt><rdf:li xml:lang="x-default">%s</rdf:li></rdf:Alt></dc:description>`, meta.Description))
	}
	if len(meta.Keywords) > 0 {
		sb.WriteString(`<dc:subject><rdf:Bag>`)
		for _, k := range meta.Keywords {
			sb.WriteString(fmt.Sprintf(`<rdf:li>%s</rdf:li>`, k))
		}
		sb.WriteString(`</rdf:Bag></dc:subject>`)
	}
	sb.WriteString(`</rdf:Description>`)

	// XMP Basic
	sb.WriteString(`<rdf:Description rdf:about="" xmlns:xmp="http://ns.adobe.com/xap/1.0/">`)
	sb.WriteString(fmt.Sprintf(`<xmp:CreateDate>%s</xmp:CreateDate>`, now))
	sb.WriteString(fmt.Sprintf(`<xmp:ModifyDate>%s</xmp:ModifyDate>`, now))
	sb.WriteString(fmt.Sprintf(`<xmp:MetadataDate>%s</xmp:MetadataDate>`, now))
	sb.WriteString(`<xmp:CreatorTool>Thoth PDF Engine</xmp:CreatorTool>`)
	sb.WriteString(`</rdf:Description>`)

	// PDF Schema
	sb.WriteString(`<rdf:Description rdf:about="" xmlns:pdf="http://ns.adobe.com/pdf/1.3/">`)
	sb.WriteString(`<pdf:Producer>Thoth PDF Engine</pdf:Producer>`)
	if meta.Keywords != nil {
		sb.WriteString(fmt.Sprintf(`<pdf:Keywords>%s</pdf:Keywords>`, strings.Join(meta.Keywords, ",")))
	}
	sb.WriteString(`</rdf:Description>`)

	sb.WriteString(`</rdf:RDF>`)
	sb.WriteString(`</x:xmpmeta>`)
	sb.WriteString(`<?xpacket end="w"?>`)

	return sb.String()
}
