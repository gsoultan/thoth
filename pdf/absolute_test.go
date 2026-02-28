package pdf

import (
	"bytes"
	"compress/zlib"
	"io"
	"strings"
	"testing"

	"github.com/gsoultan/thoth/document"
)

func TestDocument_AbsolutePositioning(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()

	// Add absolute paragraph
	_ = doc.AddParagraph("Absolute Text", document.NewCellStyleBuilder().Pos(100, 500).Build())

	var buf bytes.Buffer
	err := doc.Save(ctx, &buf)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data := buf.Bytes()

	// Find the stream content
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
		t.Fatalf("Failed to decompress stream: %v", err)
	}

	content := decompressed.String()
	// Check for the absolute position command in the PDF (Td)
	// For Pos(100, 500) and default font size 12, the baseline is at 488
	if !strings.Contains(content, "100.00 488.00 Td") && !strings.Contains(content, "100.00 488.00 Tm") {
		t.Errorf("Expected absolute position (100.00, 488.00), not found in decompressed output: %s", content)
	}
}
