package word

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

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
