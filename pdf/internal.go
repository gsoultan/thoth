package pdf

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"

	"github.com/gsoultan/thoth/pdf/internal/objects"
)

func extractTextFromStream(data []byte) string {
	var sb strings.Builder
	s := string(data)

	// Very naive contentItems stream parser
	lines := strings.Split(s, "\n")
	for _, line := range lines {
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
	fmt.Fprintf(writer, "%%PDF-1.7\n")

	offsets := make(map[string]int64)
	currOffset := int64(9) // %PDF-2.0\n

	for _, obj := range p.objects {
		if ind, ok := obj.(*objects.IndirectObject); ok {
			key := fmt.Sprintf("%d %d", ind.Number, ind.Generation)
			offsets[key] = currOffset
			s := ind.String() + "\n"
			n, _ := fmt.Fprint(writer, s)
			currOffset += int64(n)
		}
	}

	xrefStart := currOffset
	fmt.Fprintf(writer, "xref\n0 %d\n0000000000 65535 f\n", len(offsets)+1)
	for i := 1; i <= len(offsets); i++ {
		// This is simplified, assumes contiguous numbers starting from 1
		for key, offset := range offsets {
			var n, g int
			fmt.Sscanf(key, "%d %d", &n, &g)
			if n == i {
				fmt.Fprintf(writer, "%010d %05d n\r\n", offset, g)
				break
			}
		}
	}

	trailer := objects.Dictionary{
		"Size": objects.Integer(len(offsets) + 1),
	}
	// Find root and info
	for _, obj := range p.objects {
		if ind, ok := obj.(*objects.IndirectObject); ok {
			if dict, ok := ind.Data.(objects.Dictionary); ok {
				if dict["Type"] == objects.Name("Catalog") {
					trailer["Root"] = objects.Reference{Number: ind.Number, Generation: ind.Generation}
				}
			}
		}
	}

	fmt.Fprintf(writer, "trailer\n%s\nstartxref\n%d\n%%%%EOF\n", trailer.String(), xrefStart)
	return nil
}

func getImageData(path string) ([]byte, int, int, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, "", err
	}
	defer file.Close()

	cfg, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, 0, 0, "", err
	}

	if format == "jpeg" {
		file.Close()
		data, err := os.ReadFile(path)
		return data, cfg.Width, cfg.Height, "jpeg", err
	}

	file.Seek(0, 0)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, 0, "", err
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	rgb := make([]byte, w*h*3)
	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgb[idx] = byte(r >> 8)
			rgb[idx+1] = byte(g >> 8)
			rgb[idx+2] = byte(b >> 8)
			idx += 3
		}
	}
	return rgb, w, h, format, nil
}

func wrapText(text string, maxWidth float64, fontSize float64) []string {
	var lines []string
	paragraphs := strings.Split(text, "\n")
	for _, p := range paragraphs {
		words := strings.Fields(p)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		currentLine := words[0]
		for _, word := range words[1:] {
			if getTextWidth(currentLine+" "+word, fontSize) <= maxWidth {
				currentLine += " " + word
			} else {
				lines = append(lines, currentLine)
				currentLine = word
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
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return 0, 0, 0
	}
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0
}

func getTextWidth(text string, fontSize float64) float64 {
	// Crude estimation for Helvetica
	width := 0.0
	for _, r := range text {
		if r >= 'A' && r <= 'Z' {
			width += 0.7
		} else if r >= 'a' && r <= 'z' {
			width += 0.5
		} else if r >= '0' && r <= '9' {
			width += 0.55
		} else if r == ' ' {
			width += 0.3
		} else {
			width += 0.6
		}
	}
	return width * fontSize
}

func escapePDF(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}
