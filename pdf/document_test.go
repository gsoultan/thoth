package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/gsoultan/thoth/document"
	"io"
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
	if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF-2.0")) {
		t.Errorf("Expected %%PDF-2.0 header, got %s", content)
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

func TestDocument_Highlight(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()

	// Add highlighted paragraph
	_ = doc.AddParagraph("Highlighted Text", document.CellStyle{Background: "FFFF00"})

	var buf bytes.Buffer
	err := doc.Save(ctx, &buf)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data := buf.Bytes()

	// Find the stream content. Content stream is usually the first long stream.
	streamMarker := []byte("stream\n")
	startIdx := bytes.Index(data, streamMarker)
	if startIdx == -1 {
		t.Fatalf("No stream found in PDF")
	}
	startIdx += len(streamMarker)

	endIdx := bytes.Index(data[startIdx:], []byte("\nendstream"))
	if endIdx == -1 {
		t.Fatalf("No endstream found in PDF")
	}
	endIdx += startIdx

	compressedData := data[startIdx:endIdx]

	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		t.Fatalf("Failed to create zlib reader: %v", err)
	}
	defer r.Close()

	var decompressed bytes.Buffer
	_, err = io.Copy(&decompressed, r)
	if err != nil {
		// If decompression fails, it might not be the content stream or it might not be compressed.
		// For this test, it should be compressed.
		t.Fatalf("Failed to decompress stream: %v", err)
	}

	content := decompressed.String()
	// Check for the highlight rectangle command in the PDF
	// It should look something like: 1.00 1.00 0.00 rg ... re f
	if !strings.Contains(content, "1.00 1.00 0.00 rg") {
		t.Errorf("Expected highlight color (yellow: 1.0, 1.0, 0.0), not found in decompressed output: %s", content)
	}
	if !strings.Contains(content, "re f") {
		t.Errorf("Expected rectangle fill command (re f), not found in decompressed output: %s", content)
	}
}
