package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

type mediaProcessor struct{ *state }

func (e *mediaProcessor) insertImage(sheet string, path string, x, y float64) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read image file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		ext = ".png"
	}
	imgName := fmt.Sprintf("image%d%s", len(e.media)+1, ext)
	mediaPath := "xl/media/" + imgName
	e.media[mediaPath] = data

	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	// 1. Get or create drawing for this sheet
	drawingPath, err := e.getOrCreateDrawing(sheet)
	if err != nil {
		return err
	}
	dr := e.drawings[drawingPath]

	// 2. Get or create drawing relationships
	drRelsPath := "xl/drawings/_rels/" + filepath.Base(drawingPath) + ".rels"
	if e.sheetRels[drRelsPath] == nil {
		e.sheetRels[drRelsPath] = &xmlstructs.Relationships{}
	}
	drRels := e.sheetRels[drRelsPath]

	// 3. Add image to drawing relationships
	rID := drRels.AddRelationship(
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/image",
		"../media/"+imgName,
	)

	// 4. Add anchor to drawing
	// For now, simple OneCellAnchor at (x, y) cells
	// In Excel, x and y here represent Col and Row
	anchor := xmlstructs.Anchor{
		OneCellAnchor: &xmlstructs.OneCellAnchor{
			From: xmlstructs.Marker{
				Col: int(x),
				Row: int(y),
			},
			Ext: xmlstructs.Extent{
				Cx: 600000, // Default size
				Cy: 600000,
			},
			Pic: &xmlstructs.Pic{
				NvPicPr: xmlstructs.NvPicPr{
					CNvPr: xmlstructs.CNvPr{
						ID:   len(dr.Anchors) + 1,
						Name: imgName,
					},
				},
				BlipFill: xmlstructs.BlipFill{
					Blip: xmlstructs.Blip{
						Embed: rID,
					},
				},
				SpPr: xmlstructs.SpPr{
					PrstGeom: xmlstructs.PrstGeom{
						Prst: "rect",
					},
				},
			},
			ClientData: &xmlstructs.Any{},
		},
	}
	dr.Anchors = append(dr.Anchors, anchor)

	// Ensure worksheet has drawing reference
	if ws.Drawing == nil {
		// Need to add relationship from sheet to drawing
		if e.sheetRels[sheet] == nil {
			e.sheetRels[sheet] = &xmlstructs.Relationships{}
		}
		sRels := e.sheetRels[sheet]

		// Find sheet index
		sheetIdx := 1
		for i, s := range e.workbook.Sheets {
			if s.Name == sheet {
				sheetIdx = i + 1
				break
			}
		}

		drRID := sRels.AddRelationship(
			"http://schemas.openxmlformats.org/officeDocument/2006/relationships/drawing",
			fmt.Sprintf("../drawings/drawing%d.xml", sheetIdx),
		)
		ws.Drawing = &xmlstructs.WsDrawing{RID: drRID}
	}

	return nil
}

func (e *mediaProcessor) getOrCreateDrawing(sheet string) (string, error) {
	sheetIdx := 1
	for i, s := range e.workbook.Sheets {
		if s.Name == sheet {
			sheetIdx = i + 1
			break
		}
	}
	path := fmt.Sprintf("xl/drawings/drawing%d.xml", sheetIdx)
	if _, ok := e.drawings[path]; !ok {
		e.drawings[path] = &xmlstructs.WsDr{}
	}
	return path, nil
}
