package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"strings"
	"testing"
)

func TestDocument_Save(t *testing.T) {
	doc := NewDocument()
	ctx := t.Context()

	var buf bytes.Buffer
	err := doc.Save(ctx, &buf)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	content := buf.String()
	if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF-1.7")) {
		t.Errorf("Expected %%PDF-1.7 header, got %s", content)
	}
}

func TestDocument_OpenAndRead(t *testing.T) {
	// Create a minimal PDF in memory
	pdfContent := "%PDF-2.0\n" +
		"1 0 obj\n(Hello World)\nendobj\n" +
		"trailer\n<< /Root 1 0 R >>\n" +
		"%%EOF"

	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	err := doc.Open(ctx, bytes.NewReader([]byte(pdfContent)))
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	content, err := doc.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent failed: %v", err)
	}

	if content != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", content)
	}
}

func TestDocument_Search(t *testing.T) {
	// PDF with text in a stream
	pdfContent := "%PDF-2.0\n" +
		"1 0 obj\n<< /Length 25 >>\nstream\nBT (Hello Search) Tj ET\nendstream\nendobj\n" +
		"trailer\n<< /Root 1 0 R >>\n" +
		"%%EOF"

	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	err := doc.Open(ctx, bytes.NewReader([]byte(pdfContent)))
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	results, err := doc.Search([]string{"Search"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) == 0 {
		t.Errorf("Expected search result for 'Search', got none")
	}
}

func TestDocument_FlateDecode(t *testing.T) {
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	zw.Write([]byte("BT (Hello Flate) Tj ET"))
	zw.Close()
	data := buf.Bytes()

	pdfContent := fmt.Sprintf("%%PDF-2.0\n"+
		"1 0 obj\n<< /Length %d /Filter /FlateDecode >>\nstream\n%s\nendstream\nendobj\n"+
		"trailer\n<< /Root 1 0 R >>\n"+
		"%%%%EOF", len(data), string(data))

	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	err := doc.Open(ctx, bytes.NewReader([]byte(pdfContent)))
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	content, err := doc.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent failed: %v", err)
	}

	if !strings.Contains(content, "Hello Flate") {
		t.Errorf("Expected contentItems to contain 'Hello Flate', got '%s'", content)
	}
}
