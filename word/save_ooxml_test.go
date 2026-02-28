package word

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

// Test_Save_AddsOverridesAndRels ensures dynamic parts are reflected in
// [Content_Types].xml and word/_rels/document.xml.rels
func Test_Save_AddsOverridesAndRels(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	// Create dynamic parts
	if err := doc.SetHeader("Header"); err != nil {
		t.Fatalf("SetHeader: %v", err)
	}
	if err := doc.SetFooter("Footer"); err != nil {
		t.Fatalf("SetFooter: %v", err)
	}
	if err := doc.AddList([]string{"A", "B"}, true); err != nil { // numbering.xml
		t.Fatalf("AddList: %v", err)
	}
	if err := doc.AddFootnote("Note"); err != nil { // footnotes.xml
		t.Fatalf("AddFootnote: %v", err)
	}

	// Save into memory
	var buf bytes.Buffer
	if err := doc.Save(ctx, &buf); err != nil {
		t.Fatalf("Save: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader: %v", err)
	}

	// Helper to open and decode XML
	openXML := func(name string, v any) error {
		for _, f := range zr.File {
			if f.Name == name {
				r, err := f.Open()
				if err != nil {
					return err
				}
				defer r.Close()
				return xml.NewDecoder(r).Decode(v)
			}
		}
		return nil
	}

	// Check [Content_Types].xml overrides
	var ct xmlstructs.ContentTypes
	if err := openXML("[Content_Types].xml", &ct); err != nil {
		t.Fatalf("read [Content_Types].xml: %v", err)
	}

	reqParts := map[string]bool{
		"/word/header1.xml":   false,
		"/word/footer1.xml":   false,
		"/word/numbering.xml": false,
		"/word/footnotes.xml": false,
	}
	for _, ov := range ct.Override {
		if _, ok := reqParts[ov.PartName]; ok {
			reqParts[ov.PartName] = true
		}
	}
	for part, seen := range reqParts {
		if !seen {
			t.Fatalf("missing Override for %s", part)
		}
	}

	// Check document relationships
	var rels xmlstructs.Relationships
	if err := openXML("word/_rels/document.xml.rels", &rels); err != nil {
		t.Fatalf("read document.xml.rels: %v", err)
	}

	reqRels := map[string]bool{
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering:numbering.xml": false,
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/footnotes:footnotes.xml": false,
	}
	for _, r := range rels.Rels {
		key := r.Type + ":" + r.Target
		if _, ok := reqRels[key]; ok {
			reqRels[key] = true
		}
	}
	for key, seen := range reqRels {
		if !seen {
			t.Fatalf("missing relationship %s", key)
		}
	}
}
