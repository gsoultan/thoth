package excel

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

// lifecycle handles document lifecycle operations.
type lifecycle struct{ *state }

// Open loads a document from a reader.
func (e *lifecycle) Open(ctx context.Context, reader io.Reader) error {
	// zip.NewReader requires ReaderAt, so we buffer to temp file
	tmp, err := os.CreateTemp("", "thoth-excel-*.xlsx")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	e.tempFile = tmp

	_, err = io.Copy(tmp, reader)
	if err != nil {
		return fmt.Errorf("buffer reader: %w", err)
	}

	zr, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return fmt.Errorf("open zip reader: %w", err)
	}
	e.reader = zr

	return e.loadCore(ctx)
}

// Save writes the document to a writer.
func (e *lifecycle) Save(ctx context.Context, writer io.Writer) error {
	e.prepareSheets()
	e.prepareContentTypes()
	zw := zip.NewWriter(writer)
	defer zw.Close()

	// Keep track of files we handle manually
	handled := make(map[string]bool)

	// Save main parts
	if err := e.saveCoreParts(zw, handled); err != nil {
		return err
	}

	// Save sheets and relationships
	if err := e.saveSheets(zw, handled); err != nil {
		return err
	}

	// Save media
	if err := e.saveMedia(zw, handled); err != nil {
		return err
	}

	// Save drawings
	if err := e.saveDrawings(zw, handled); err != nil {
		return err
	}

	// Save tables
	if err := e.saveTables(zw, handled); err != nil {
		return err
	}

	// Copy remaining files from original reader
	return e.copyRemainingFiles(zw, handled)
}

func (e *lifecycle) saveCoreParts(zw *zip.Writer, handled map[string]bool) error {
	if e.workbook != nil {
		if err := e.writeXML(zw, "xl/workbook.xml", e.workbook); err != nil {
			return err
		}
		handled["xl/workbook.xml"] = true
	}
	if e.sharedStrings != nil {
		if err := e.writeXML(zw, "xl/sharedStrings.xml", e.sharedStrings); err != nil {
			return err
		}
		handled["xl/sharedStrings.xml"] = true

		// Ensure sharedStrings relationship exists in workbookRels
		found := false
		for _, rel := range e.workbookRels.Rels {
			if rel.Type == "http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings" {
				found = true
				break
			}
		}
		if !found {
			e.workbookRels.AddRelationship("http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings", "sharedStrings.xml")
		}
	}
	if e.coreProperties != nil {
		if err := e.writeXML(zw, "docProps/core.xml", e.coreProperties); err != nil {
			return err
		}
		handled["docProps/core.xml"] = true

		// Ensure core properties relationship exists in rootRels
		found := false
		for _, rel := range e.rootRels.Rels {
			if rel.Type == "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" {
				found = true
				break
			}
		}
		if !found {
			e.rootRels.AddRelationship("http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties", "docProps/core.xml")
		}
	}
	if e.styles != nil {
		if err := e.writeXML(zw, "xl/styles.xml", e.styles); err != nil {
			return err
		}
		handled["xl/styles.xml"] = true

		// Ensure styles relationship exists
		found := false
		for _, rel := range e.workbookRels.Rels {
			if rel.Target == "styles.xml" || rel.Target == "/xl/styles.xml" {
				found = true
				break
			}
		}
		if !found {
			e.workbookRels.AddRelationship("http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles", "styles.xml")
		}
	}
	if e.workbookRels != nil && e.wbRelsPath != "" {
		if err := e.writeXML(zw, e.wbRelsPath, e.workbookRels); err != nil {
			return err
		}
		handled[e.wbRelsPath] = true
	}

	// Finalize and write root-level files that might have been updated during the process
	if e.rootRels != nil && e.rootRelsPath != "" {
		if err := e.writeXML(zw, e.rootRelsPath, e.rootRels); err != nil {
			return err
		}
		handled[e.rootRelsPath] = true
	}

	if e.contentTypes != nil {
		if err := e.writeXML(zw, "[Content_Types].xml", e.contentTypes); err != nil {
			return err
		}
		handled["[Content_Types].xml"] = true
	}

	return nil
}

func (e *lifecycle) saveSheets(zw *zip.Writer, handled map[string]bool) error {
	for name, ws := range e.sheets {
		if ws.XMLNS_R == "" {
			ws.XMLNS_R = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
		}
		for i, s := range e.workbook.Sheets {
			if s.Name == name {
				path := fmt.Sprintf("xl/worksheets/sheet%d.xml", i+1)
				if err := e.writeXML(zw, path, ws); err != nil {
					return err
				}
				handled[path] = true
				break
			}
		}
	}
	for key, rels := range e.sheetRels {
		if strings.HasSuffix(key, ".rels") {
			// This is already a path
			if err := e.writeXML(zw, key, rels); err != nil {
				return err
			}
			handled[key] = true
			continue
		}
		// This is a sheet name
		for i, s := range e.workbook.Sheets {
			if s.Name == key {
				path := fmt.Sprintf("xl/worksheets/_rels/sheet%d.xml.rels", i+1)
				if err := e.writeXML(zw, path, rels); err != nil {
					return err
				}
				handled[path] = true
				break
			}
		}
	}
	return nil
}

func (e *lifecycle) saveMedia(zw *zip.Writer, handled map[string]bool) error {
	for name, data := range e.media {
		f, err := zw.Create(name)
		if err != nil {
			return fmt.Errorf("create media %s: %w", name, err)
		}
		_, err = f.Write(data)
		if err != nil {
			return fmt.Errorf("write media %s: %w", name, err)
		}
		handled[name] = true
	}
	return nil
}

func (e *lifecycle) saveDrawings(zw *zip.Writer, handled map[string]bool) error {
	for path, dr := range e.drawings {
		if dr.XMLNS_R == "" {
			dr.XMLNS_R = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
		}
		if err := e.writeXML(zw, path, dr); err != nil {
			return err
		}
		handled[path] = true
	}
	return nil
}

func (e *lifecycle) saveTables(zw *zip.Writer, handled map[string]bool) error {
	for path, table := range e.tables {
		if err := e.writeXML(zw, path, table); err != nil {
			return err
		}
		handled[path] = true
	}
	return nil
}

func (e *lifecycle) copyRemainingFiles(zw *zip.Writer, handled map[string]bool) error {
	if e.reader == nil {
		return nil
	}
	for _, f := range e.reader.File {
		if handled[f.Name] {
			continue
		}
		if err := e.copyFile(f, zw); err != nil {
			return fmt.Errorf("copy file %s: %w", f.Name, err)
		}
	}
	return nil
}

// Close releases any resources used by the document, such as temporary files.
func (e *lifecycle) prepareContentTypes() {
	if e.contentTypes == nil {
		e.contentTypes = xmlstructs.NewContentTypes()
	}

	if e.workbook != nil {
		e.contentTypes.AddOverride("/xl/workbook.xml", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml")
	}

	if e.styles != nil {
		e.contentTypes.AddOverride("/xl/styles.xml", "application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml")
	}

	if e.sharedStrings != nil {
		e.contentTypes.AddOverride("/xl/sharedStrings.xml", "application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml")
	}

	if e.coreProperties != nil {
		e.contentTypes.AddOverride("/docProps/core.xml", "application/vnd.openxmlformats-package.core-properties+xml")
	}

	for name := range e.drawings {
		e.contentTypes.AddOverride("/"+name, "application/vnd.openxmlformats-officedocument.drawing+xml")
	}

	for path := range e.tables {
		e.contentTypes.AddOverride("/"+path, "application/vnd.openxmlformats-officedocument.spreadsheetml.table+xml")
	}

	for name := range e.sheets {
		for i, s := range e.workbook.Sheets {
			if s.Name == name {
				path := fmt.Sprintf("/xl/worksheets/sheet%d.xml", i+1)
				e.contentTypes.AddOverride(path, "application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml")
				break
			}
		}
	}

	// Add defaults for media types
	for name := range e.media {
		ext := strings.ToLower(filepath.Ext(name))
		if strings.HasPrefix(ext, ".") {
			ext = ext[1:]
		}
		switch ext {
		case "png":
			e.contentTypes.AddDefault(ext, "image/png")
		case "jpg", "jpeg":
			e.contentTypes.AddDefault(ext, "image/jpeg")
		case "gif":
			e.contentTypes.AddDefault(ext, "image/gif")
		case "bmp":
			e.contentTypes.AddDefault(ext, "image/bmp")
		case "tif", "tiff":
			e.contentTypes.AddDefault(ext, "image/tiff")
		}
	}
}

func (e *lifecycle) prepareSheets() {
	for _, ws := range e.sheets {
		if ws.Dimension == nil {
			ws.Dimension = &xmlstructs.Dimension{Ref: e.calculateDimension(ws)}
		}
		if ws.SheetFormatPr == nil {
			ws.SheetFormatPr = &xmlstructs.SheetFormatPr{
				DefaultRowHeight: 15.0,
			}
		}
		if ws.PageMargins == nil {
			ws.PageMargins = &xmlstructs.PageMargins{
				Left:   0.7,
				Right:  0.7,
				Top:    0.75,
				Bottom: 0.75,
				Header: 0.3,
				Footer: 0.3,
			}
		}
	}
}

func (e *lifecycle) calculateDimension(ws *xmlstructs.Worksheet) string {
	if len(ws.SheetData.Rows) == 0 {
		return "A1"
	}

	minCol, minRow := "ZZZ", 999999
	maxCol, maxRow := "A", 0

	for _, row := range ws.SheetData.Rows {
		if row.R < minRow {
			minRow = row.R
		}
		if row.R > maxRow {
			maxRow = row.R
		}
		for _, cell := range row.Cells {
			col := getColumnFromAxis(cell.R)
			if compareColumns(col, minCol) < 0 {
				minCol = col
			}
			if compareColumns(col, maxCol) > 0 {
				maxCol = col
			}
		}
	}

	if maxRow == 0 {
		return "A1"
	}

	if minCol == maxCol && minRow == maxRow {
		return fmt.Sprintf("%s%d", minCol, minRow)
	}

	return fmt.Sprintf("%s%d:%s%d", minCol, minRow, maxCol, maxRow)
}

func (e *lifecycle) Close() error {
	if e.reader != nil {
		e.reader.Close()
	}
	if e.tempFile != nil {
		name := e.tempFile.Name()
		e.tempFile.Close()
		os.Remove(name)
	}
	return nil
}
