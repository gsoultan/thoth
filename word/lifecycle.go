package word

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
func (w *lifecycle) Open(ctx context.Context, reader io.Reader) error {
	tmp, err := os.CreateTemp("", "thoth-word-*.docx")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	w.tempFile = tmp

	_, err = io.Copy(tmp, reader)
	if err != nil {
		return fmt.Errorf("buffer reader: %w", err)
	}

	zr, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return fmt.Errorf("open zip reader: %w", err)
	}
	w.reader = zr

	return w.loadCore(ctx)
}

// Save writes the document to a writer.
func (w *lifecycle) Save(ctx context.Context, writer io.Writer) error {
	zw := zip.NewWriter(writer)
	defer zw.Close()

	handled := make(map[string]bool)

	// Save specialized parts
	if err := w.saveParts(zw, handled); err != nil {
		return err
	}

	// Save media
	if err := w.saveMedia(zw, handled); err != nil {
		return err
	}

	// Copy remaining files from original reader
	return w.copyRemainingFiles(zw, handled)
}

func (w *lifecycle) saveParts(zw *zip.Writer, handled map[string]bool) error {
	if w.contentTypes != nil {
		if err := w.writeXML(zw, "[Content_Types].xml", w.contentTypes); err != nil {
			return err
		}
		handled["[Content_Types].xml"] = true
	}
	if w.rootRels != nil {
		if err := w.writeXML(zw, "_rels/.rels", w.rootRels); err != nil {
			return err
		}
		handled["_rels/.rels"] = true
	}
	if w.doc != nil {
		if err := w.writeXML(zw, "word/document.xml", w.doc); err != nil {
			return err
		}
		handled["word/document.xml"] = true
	}
	if w.coreProperties != nil {
		if err := w.writeXML(zw, "docProps/core.xml", w.coreProperties); err != nil {
			return err
		}
		handled["docProps/core.xml"] = true
	}
	if w.appProperties != nil {
		if err := w.writeXML(zw, "docProps/app.xml", w.appProperties); err != nil {
			return err
		}
		handled["docProps/app.xml"] = true
	}
	if w.docRels != nil {
		if err := w.writeXML(zw, "word/_rels/document.xml.rels", w.docRels); err != nil {
			return err
		}
		handled["word/_rels/document.xml.rels"] = true
	}
	return nil
}

func (w *lifecycle) saveMedia(zw *zip.Writer, handled map[string]bool) error {
	for name, data := range w.media {
		f, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		if err != nil {
			return err
		}
		handled[name] = true
	}
	return nil
}

func (w *lifecycle) copyRemainingFiles(zw *zip.Writer, handled map[string]bool) error {
	if w.reader == nil {
		return nil
	}
	for _, f := range w.reader.File {
		if handled[f.Name] {
			continue
		}
		if err := w.copyFile(f, zw); err != nil {
			return err
		}
	}
	return nil
}

// Close releases any resources used by the document, such as temporary files.
func (w *lifecycle) Close() error {
	if w.reader != nil {
		w.reader.Close()
	}
	if w.tempFile != nil {
		name := w.tempFile.Name()
		w.tempFile.Close()
		os.Remove(name)
	}
	return nil
}
