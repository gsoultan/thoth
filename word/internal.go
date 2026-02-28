package word

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func (p *processor) mapTableCellProperties(s document.CellStyle) *xmlstructs.TableCellProperties {
	tcPr := &xmlstructs.TableCellProperties{}
	if s.Background != "" {
		tcPr.Shd = &xmlstructs.TableCellShading{
			Val:  "clear",
			Fill: s.Background,
		}
	}
	if s.Vertical != "" {
		val := s.Vertical
		if val == "middle" {
			val = "center"
		}
		tcPr.VAlign = &xmlstructs.VAlign{Val: val}
	}
	if s.Border || s.BorderTop || s.BorderBottom || s.BorderLeft || s.BorderRight {
		tcPr.TcBorders = &xmlstructs.TableCellBorders{}
		borderVal := "single"
		borderSz := 4
		if s.BorderWidth > 0 {
			borderSz = int(s.BorderWidth * 8)
		}
		borderColor := "auto"
		if s.BorderColor != "" {
			borderColor = s.BorderColor
		}

		line := &xmlstructs.BorderLine{Val: borderVal, Sz: borderSz, Color: borderColor}

		if s.Border || s.BorderTop {
			tcPr.TcBorders.Top = line
		}
		if s.Border || s.BorderLeft {
			tcPr.TcBorders.Left = line
		}
		if s.Border || s.BorderBottom {
			tcPr.TcBorders.Bottom = line
		}
		if s.Border || s.BorderRight {
			tcPr.TcBorders.Right = line
		}
	}
	if s.Padding > 0 {
		twips := int(s.Padding * 20)
		tcPr.TcMar = &xmlstructs.TableCellMargins{
			Top:    &xmlstructs.TableCellW{W: twips, Type: "dxa"},
			Left:   &xmlstructs.TableCellW{W: twips, Type: "dxa"},
			Bottom: &xmlstructs.TableCellW{W: twips, Type: "dxa"},
			Right:  &xmlstructs.TableCellW{W: twips, Type: "dxa"},
		}
	}
	return tcPr
}

func (p *processor) mapParagraphProperties(s document.CellStyle) *xmlstructs.ParagraphProperties {
	pPr := &xmlstructs.ParagraphProperties{}
	if s.Name != "" {
		pPr.PStyle = &xmlstructs.ParagraphStyle{Val: s.Name}
	}
	if s.Horizontal != "" {
		val := s.Horizontal
		if val == "justify" {
			val = "both"
		}
		pPr.Jc = &xmlstructs.Justification{Val: val}
	}
	if s.KeepWithNext {
		pPr.KeepNext = &struct{}{}
	}
	if s.KeepTogether {
		pPr.KeepLines = &struct{}{}
	}
	if s.Indent != 0 {
		pPr.Ind = &xmlstructs.Ind{Left: int(s.Indent * 20)}
	}
	if s.Hanging != 0 {
		if pPr.Ind == nil {
			pPr.Ind = &xmlstructs.Ind{}
		}
		pPr.Ind.Hanging = int(s.Hanging * 20)
	}

	if s.SpacingBefore != 0 || s.SpacingAfter != 0 || s.LineSpacing != 0 {
		pPr.Spacing = &xmlstructs.Spacing{
			Before: int(s.SpacingBefore * 20),
			After:  int(s.SpacingAfter * 20),
		}
		if s.LineSpacing != 0 {
			pPr.Spacing.Line = int(s.LineSpacing * 240)
			pPr.Spacing.LineRule = "auto"
		}
	}

	return pPr
}

func (p *processor) mapRunProperties(s document.CellStyle) *xmlstructs.RunProperties {
	rPr := &xmlstructs.RunProperties{}
	if s.Bold {
		rPr.Bold = &struct{}{}
	}
	if s.Italic {
		rPr.Italic = &struct{}{}
	}
	if s.Size > 0 {
		rPr.Sz = &xmlstructs.ValInt{Val: s.Size * 2}
		rPr.SzCs = &xmlstructs.ValInt{Val: s.Size * 2}
	}
	if s.Color != "" {
		rPr.Color = &xmlstructs.Color{Val: s.Color}
	}
	if s.Superscript {
		rPr.VertAlign = &xmlstructs.ValStr{Val: "superscript"}
	}
	if s.Subscript {
		rPr.VertAlign = &xmlstructs.ValStr{Val: "subscript"}
	}
	if s.Background != "" {
		rPr.Shd = &xmlstructs.TableCellShading{
			Val:  "clear",
			Fill: s.Background,
		}
	}
	return rPr
}

func (p *processor) ensureNumbering() {
	if p.numbering != nil {
		return
	}

	bulletLevels := make([]xmlstructs.Level, 9)
	bulletTexts := []string{"•", "○", "■", "•", "○", "■", "•", "○", "■"}
	for i := 0; i < 9; i++ {
		bulletLevels[i] = xmlstructs.Level{
			ILvl:    i,
			Start:   &xmlstructs.ValInt{Val: 1},
			NumFmt:  &xmlstructs.ValStr{Val: "bullet"},
			LvlText: &xmlstructs.ValStr{Val: bulletTexts[i]},
			LvlJc:   &xmlstructs.ValStr{Val: "left"},
			PPr: &xmlstructs.ParagraphProperties{
				Ind: &xmlstructs.Ind{Left: (i + 1) * 720, Hanging: 360},
			},
		}
	}

	orderedLevels := make([]xmlstructs.Level, 9)
	fmts := []string{"decimal", "lowerLetter", "lowerRoman", "decimal", "lowerLetter", "lowerRoman", "decimal", "lowerLetter", "lowerRoman"}
	for i := 0; i < 9; i++ {
		orderedLevels[i] = xmlstructs.Level{
			ILvl:    i,
			Start:   &xmlstructs.ValInt{Val: 1},
			NumFmt:  &xmlstructs.ValStr{Val: fmts[i]},
			LvlText: &xmlstructs.ValStr{Val: fmt.Sprintf("%%%d.", i+1)},
			LvlJc:   &xmlstructs.ValStr{Val: "left"},
			PPr: &xmlstructs.ParagraphProperties{
				Ind: &xmlstructs.Ind{Left: (i + 1) * 720, Hanging: 360},
			},
		}
	}

	p.numbering = &xmlstructs.Numbering{
		W: "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
		AbstractNums: []xmlstructs.AbstractNum{
			{
				AbstractNumID: 1,
				Levels:        bulletLevels,
			},
			{
				AbstractNumID: 2,
				Levels:        orderedLevels,
			},
		},
		Nums: []xmlstructs.Num{
			{
				NumID:         1,
				AbstractNumID: &xmlstructs.ValInt{Val: 1},
			},
			{
				NumID:         2,
				AbstractNumID: &xmlstructs.ValInt{Val: 2},
			},
		},
	}
}

func (w *state) loadCore(ctx context.Context) error {
	// 0. Load Content Types
	var ct xmlstructs.ContentTypes
	if err := w.loadXML("[Content_Types].xml", &ct); err == nil {
		w.contentTypes = &ct
	} else {
		w.contentTypes = xmlstructs.NewContentTypes()
	}

	// 1. Load root relationships to find document.xml
	var rootRels xmlstructs.Relationships
	if err := w.loadXML("_rels/.rels", &rootRels); err != nil {
		// Fallback to hardcoded path if .rels is missing (not standard but for robustness)
		var doc xmlstructs.Document
		if err := w.loadXML("word/document.xml", &doc); err != nil {
			return fmt.Errorf("load document.xml fallback: %w", err)
		}
		w.doc = &doc
		w.xmlDoc = w.doc
		w.rootRels = &xmlstructs.Relationships{
			Rels: []xmlstructs.Relationship{
				{
					ID:     "rId1",
					Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument",
					Target: "word/document.xml",
				},
			},
		}
	} else {
		w.rootRels = &rootRels
		docPath := rootRels.TargetByType("http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument")
		if docPath == "" {
			docPath = "word/document.xml"
		} else if strings.HasPrefix(docPath, "/") {
			docPath = docPath[1:]
		}

		var doc xmlstructs.Document
		if err := w.loadXML(docPath, &doc); err != nil {
			return fmt.Errorf("load document.xml: %w", err)
		}
		w.doc = &doc
		w.xmlDoc = w.doc

		// Document Relationships
		var drPath string
		if idx := strings.LastIndex(docPath, "/"); idx != -1 {
			drPath = docPath[:idx] + "/_rels/" + docPath[idx+1:] + ".rels"
		} else {
			drPath = "_rels/" + docPath + ".rels"
		}
		var dr xmlstructs.Relationships
		if err := w.loadXML(drPath, &dr); err == nil {
			w.docRels = &dr
		}

		// Core Properties
		cpPath := rootRels.TargetByType("http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties")
		if cpPath != "" {
			if strings.HasPrefix(cpPath, "/") {
				cpPath = cpPath[1:]
			}
			var cp xmlstructs.CoreProperties
			if err := w.loadXML(cpPath, &cp); err == nil {
				w.coreProperties = &cp
			}
		}
	}

	return nil
}

func (w *state) loadXML(name string, target any) error {
	f, err := w.reader.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return xml.NewDecoder(f).Decode(target)
}

func (w *state) writeXML(zw *zip.Writer, name string, data any) error {
	wtr, err := zw.Create(name)
	if err != nil {
		return err
	}
	fmt.Fprint(wtr, xml.Header)
	return xml.NewEncoder(wtr).Encode(data)
}

func (w *state) copyFile(f *zip.File, zw *zip.Writer) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	wtr, err := zw.Create(f.Name)
	if err != nil {
		return err
	}
	_, err = io.Copy(wtr, rc)
	return err
}
