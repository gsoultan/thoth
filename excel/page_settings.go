package excel

import (
	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// pageSettings handles page layout settings.
type pageSettings struct{ *state }

// SetPageSettings configures the document's layout.
func (e *pageSettings) SetPageSettings(settings document.PageSettings) error {
	if e.workbook == nil {
		e.workbook = &xmlstructs.Workbook{
			Sheets: make([]xmlstructs.Sheet, 0),
		}
	}

	for _, ws := range e.sheets {
		if ws.PageSetup == nil {
			ws.PageSetup = &xmlstructs.PageSetup{}
		}

		if settings.Orientation == document.OrientationLandscape {
			ws.PageSetup.Orientation = "landscape"
		} else {
			ws.PageSetup.Orientation = "portrait"
		}

		switch settings.PaperType {
		case document.PaperA4:
			ws.PageSetup.PaperSize = 9
		case document.PaperLetter:
			ws.PageSetup.PaperSize = 1
		}

		ws.PageMargins = &xmlstructs.PageMargins{
			Top:    settings.Margins.Top,
			Bottom: settings.Margins.Bottom,
			Left:   settings.Margins.Left,
			Right:  settings.Margins.Right,
		}
	}
	return nil
}
