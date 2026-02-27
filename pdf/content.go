package pdf

import (
	"bytes"
	"compress/zlib"
	"io"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/pdf/internal/objects"
)

// content handles reading and searching document content.
type content struct{ *state }

// ReadContent returns the text contentItems of the document.
func (p *content) ReadContent() (string, error) {
	var sb strings.Builder
	for _, obj := range p.objects {
		var data []byte
		var dict objects.Dictionary

		if stream, ok := obj.(objects.Stream); ok {
			data = stream.Data
			dict = stream.Dict
		} else if indObj, ok := obj.(*objects.IndirectObject); ok {
			if stream, ok := indObj.Data.(objects.Stream); ok {
				data = stream.Data
				dict = stream.Dict
			}
		}

		if data != nil {
			if filter, ok := dict["Filter"]; ok {
				if name, ok := filter.(objects.Name); ok && name == "FlateDecode" {
					zr, err := zlib.NewReader(bytes.NewReader(data))
					if err == nil {
						decoded, _ := io.ReadAll(zr)
						data = decoded
						zr.Close()
					}
				}
			}
			sb.WriteString(extractTextFromStream(data))
		}
	}

	if sb.Len() == 0 {
		// Fallback to simple string extraction
		for _, obj := range p.objects {
			if s, ok := obj.(objects.PDFString); ok {
				sb.WriteString(string(s))
				sb.WriteString(" ")
			} else if io, ok := obj.(*objects.IndirectObject); ok {
				if s, ok := io.Data.(objects.PDFString); ok {
					sb.WriteString(string(s))
					sb.WriteString(" ")
				}
			}
		}
	}

	return strings.TrimSpace(sb.String()), nil
}

// Search finds keywords in the document.
func (p *content) Search(keywords []string) ([]document.SearchResult, error) {
	content, err := p.ReadContent()
	if err != nil {
		return nil, err
	}
	strategy := document.NewSimpleSearchStrategy()
	return strategy.Execute(p.ctx, content, keywords)
}

// Replace replaces keywords with new values.
func (p *content) Replace(replacements map[string]string) error {
	for i := range p.objects {
		obj := p.objects[i]
		var data []byte
		var dict objects.Dictionary
		var isIndirect bool
		var indObj *objects.IndirectObject

		if stream, ok := obj.(objects.Stream); ok {
			data = stream.Data
			dict = stream.Dict
		} else if ind, ok := obj.(*objects.IndirectObject); ok {
			indObj = ind
			isIndirect = true
			if stream, ok := ind.Data.(objects.Stream); ok {
				data = stream.Data
				dict = stream.Dict
			}
		}

		if data != nil {
			// If it's FlateDecoded, we need to decompress, replace, then re-compress
			isCompressed := false
			if filter, ok := dict["Filter"]; ok {
				if name, ok := filter.(objects.Name); ok && name == "FlateDecode" {
					isCompressed = true
					zr, err := zlib.NewReader(bytes.NewReader(data))
					if err == nil {
						decoded, _ := io.ReadAll(zr)
						data = decoded
						zr.Close()
					}
				}
			}

			content := string(data)
			modified := false
			for old, new := range replacements {
				if strings.Contains(content, old) {
					content = strings.ReplaceAll(content, old, new)
					modified = true
				}
			}

			if modified {
				newData := []byte(content)
				if isCompressed {
					var buf bytes.Buffer
					zw := zlib.NewWriter(&buf)
					zw.Write(newData)
					zw.Close()
					newData = buf.Bytes()
				}

				if isIndirect {
					if stream, ok := indObj.Data.(objects.Stream); ok {
						stream.Data = newData
						stream.Dict["Length"] = objects.Integer(len(newData))
						indObj.Data = stream
					}
				} else {
					if stream, ok := obj.(objects.Stream); ok {
						stream.Data = newData
						stream.Dict["Length"] = objects.Integer(len(newData))
						p.objects[i] = stream
					}
				}
			}
		}
	}
	return nil
}
