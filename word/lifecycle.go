package word

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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
func (w *lifecycle) Save(ctx context.Context, writer io.Writer) (retErr error) {
	zw := zip.NewWriter(writer)
	defer func() { retErr = errors.Join(retErr, zw.Close()) }()

	handled := make(map[string]bool)

	// Save media
	if err := w.saveMedia(zw, handled); err != nil {
		return err
	}

	// Save specialized parts
	if err := w.saveParts(zw, handled); err != nil {
		return err
	}

	// Copy remaining files from original reader
	return w.copyRemainingFiles(zw, handled)
}

func (w *lifecycle) saveParts(zw *zip.Writer, handled map[string]bool) error {
	// 1) Write static parts first (they don't depend on late mutations)
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
	if w.styles != nil {
		if err := w.writeXML(zw, "word/styles.xml", w.styles); err != nil {
			return err
		}
		handled["word/styles.xml"] = true
	}
	if w.settings != nil {
		if err := w.writeXML(zw, "word/settings.xml", w.settings); err != nil {
			return err
		}
		handled["word/settings.xml"] = true
		if w.contentTypes != nil {
			w.contentTypes.AddOverride("/word/settings.xml", "application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml")
		}
		if w.docRels != nil {
			w.docRels.AddRelationship("http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings", "settings.xml")
		}
	}

	// 2) Dynamic parts that contribute overrides/relationships
	for id, header := range w.headers {
		path := "word/" + id
		if err := w.writeXML(zw, path, header); err != nil {
			return err
		}
		handled[path] = true
		if w.contentTypes != nil {
			w.contentTypes.AddOverride("/"+path, "application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml")
		}
		// Write header relationships if they exist
		if rels, ok := w.headerRels[id]; ok && rels != nil && len(rels.Rels) > 0 {
			relPath := "word/_rels/" + id + ".rels"
			if err := w.writeXML(zw, relPath, rels); err != nil {
				return err
			}
			handled[relPath] = true
		}
	}

	for id, footer := range w.footers {
		path := "word/" + id
		if err := w.writeXML(zw, path, footer); err != nil {
			return err
		}
		handled[path] = true
		if w.contentTypes != nil {
			w.contentTypes.AddOverride("/"+path, "application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml")
		}
		// Write footer relationships if they exist
		if rels, ok := w.footerRels[id]; ok && rels != nil && len(rels.Rels) > 0 {
			relPath := "word/_rels/" + id + ".rels"
			if err := w.writeXML(zw, relPath, rels); err != nil {
				return err
			}
			handled[relPath] = true
		}
	}

	if w.numbering != nil {
		if err := w.writeXML(zw, "word/numbering.xml", w.numbering); err != nil {
			return err
		}
		handled["word/numbering.xml"] = true
		if w.contentTypes != nil {
			w.contentTypes.AddOverride("/word/numbering.xml", "application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml")
		}
		if w.docRels != nil {
			w.docRels.AddRelationship("http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering", "numbering.xml")
		}
	}

	if w.footnotes != nil {
		if err := w.writeXML(zw, "word/footnotes.xml", w.footnotes); err != nil {
			return err
		}
		handled["word/footnotes.xml"] = true
		if w.contentTypes != nil {
			w.contentTypes.AddOverride("/word/footnotes.xml", "application/vnd.openxmlformats-officedocument.wordprocessingml.footnotes+xml")
		}
		if w.docRels != nil {
			w.docRels.AddRelationship("http://schemas.openxmlformats.org/officeDocument/2006/relationships/footnotes", "footnotes.xml")
		}
	}

	// 3) Write relationships and content types LAST so they include all mutations above
	if w.docRels != nil {
		if err := w.writeXML(zw, "word/_rels/document.xml.rels", w.docRels); err != nil {
			return err
		}
		handled["word/_rels/document.xml.rels"] = true
	}
	if w.contentTypes != nil {
		if err := w.writeXML(zw, "[Content_Types].xml", w.contentTypes); err != nil {
			return err
		}
		handled["[Content_Types].xml"] = true
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

		// Ensure extension is in Content Types
		if w.contentTypes != nil {
			if idx := strings.LastIndex(name, "."); idx != -1 {
				ext := name[idx+1:]
				switch ext {
				case "png":
					w.contentTypes.AddDefault(ext, "image/png")
				case "jpeg", "jpg":
					w.contentTypes.AddDefault(ext, "image/jpeg")
				case "gif":
					w.contentTypes.AddDefault(ext, "image/gif")
				}
			}
		}
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
