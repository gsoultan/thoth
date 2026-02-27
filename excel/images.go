package excel

import (
	"fmt"
	"os"

	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// imageProcessor handles image insertion operations.
type imageProcessor struct{ *state }

func (e *imageProcessor) insertImage(sheet string, path string, x, y float64) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read image file: %w", err)
	}

	// 1. Add to media
	imgName := fmt.Sprintf("image%d.png", len(e.media)+1)
	mediaPath := "xl/media/" + imgName
	e.media[mediaPath] = data

	// 2. Add relationship to sheet
	ws, ok := e.sheets[sheet]
	if !ok {
		return fmt.Errorf("sheet %s not found", sheet)
	}

	if e.sheetRels[sheet] == nil {
		e.sheetRels[sheet] = &xmlstructs.Relationships{}
	}

	rID := e.sheetRels[sheet].AddRelationship(
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/image",
		"../media/"+imgName,
	)

	// 3. Mark worksheet with a drawing placeholder
	if ws.Drawing == nil {
		ws.Drawing = &xmlstructs.WsDrawing{RID: rID}
	}

	return nil
}
