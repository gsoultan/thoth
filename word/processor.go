package word

import (
	"github.com/gsoultan/thoth/document"
	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// processor handles document manipulation.
type processor struct{ *state }

func (p *processor) SetPageSettings(settings document.PageSettings) error {
	if p.xmlDoc == nil || p.xmlDoc.Body.SectPr == nil {
		return nil
	}
	sect := p.xmlDoc.Body.SectPr
	if sect.PgSz == nil {
		sect.PgSz = &xmlstructs.PgSz{}
	}
	if sect.PgMar == nil {
		sect.PgMar = &xmlstructs.PgMar{}
	}

	// Conversion: 1 inch = 1440 twips
	// 1 mm = 1/25.4 inch = 1440/25.4 = 56.69 twips
	// We assume points (1/72 inch) -> 1 point = 20 twips
	sect.PgMar.Top = int(settings.Margins.Top * 20)
	sect.PgMar.Bottom = int(settings.Margins.Bottom * 20)
	sect.PgMar.Left = int(settings.Margins.Left * 20)
	sect.PgMar.Right = int(settings.Margins.Right * 20)

	if settings.Orientation == document.OrientationLandscape {
		sect.PgSz.Orient = "landscape"
		// Swap width and height if needed or use standard values
	} else {
		sect.PgSz.Orient = "portrait"
	}

	switch settings.PaperType {
	case document.PaperA4:
		sect.PgSz.W = 11906
		sect.PgSz.H = 16838
	case document.PaperLetter:
		sect.PgSz.W = 12240
		sect.PgSz.H = 15840
	}

	return nil
}

func (p *processor) RegisterFont(name, path string) error {
	// Word typically uses fonts available on the system.
	// We can add font information to styles if needed.
	return nil
}

func (p *processor) AttachFile(path, name, description string) error {
	// Not implemented for Word yet
	return nil
}
