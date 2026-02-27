package excel

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
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

	// Copy remaining files from original reader
	return e.copyRemainingFiles(zw, handled)
}

func (e *lifecycle) saveCoreParts(zw *zip.Writer, handled map[string]bool) error {
	if e.contentTypes != nil {
		if err := e.writeXML(zw, "[Content_Types].xml", e.contentTypes); err != nil {
			return err
		}
		handled["[Content_Types].xml"] = true
	}
	if e.rootRels != nil && e.rootRelsPath != "" {
		if err := e.writeXML(zw, e.rootRelsPath, e.rootRels); err != nil {
			return err
		}
		handled[e.rootRelsPath] = true
	}
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
	}
	if e.coreProperties != nil {
		if err := e.writeXML(zw, "docProps/core.xml", e.coreProperties); err != nil {
			return err
		}
		handled["docProps/core.xml"] = true
	}
	if e.styles != nil {
		if err := e.writeXML(zw, "xl/styles.xml", e.styles); err != nil {
			return err
		}
		handled["xl/styles.xml"] = true
	}
	if e.workbookRels != nil && e.wbRelsPath != "" {
		if err := e.writeXML(zw, e.wbRelsPath, e.workbookRels); err != nil {
			return err
		}
		handled[e.wbRelsPath] = true
	}
	return nil
}

func (e *lifecycle) saveSheets(zw *zip.Writer, handled map[string]bool) error {
	for name, ws := range e.sheets {
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
	for sheet, rels := range e.sheetRels {
		for i, s := range e.workbook.Sheets {
			if s.Name == sheet {
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
