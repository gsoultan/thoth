package word

import (
	"os"
	"strings"
	"testing"

	"github.com/gsoultan/thoth/word/internal/xmlstructs"
)

func TestDocument_AddParagraph(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	text := "Hello, World!"
	err := doc.AddParagraph(text)
	if err != nil {
		t.Fatalf("AddParagraph failed: %v", err)
	}

	content, err := doc.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent failed: %v", err)
	}

	if !strings.Contains(content, text) {
		t.Errorf("Expected content to contain %q, got %q", text, content)
	}
}

func TestDocument_InsertTable(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	handle, err := doc.AddTable(3, 2)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	if handle == nil {
		t.Fatal("Expected table handle, got nil")
	}

	var tableCount int
	for _, c := range doc.doc.Body.Content {
		if _, ok := c.(xmlstructs.Table); ok {
			tableCount++
		}
	}

	if tableCount != 1 {
		t.Errorf("Expected 1 table, got %d", tableCount)
	}

	tbl, err := doc.getTable(0)
	if err != nil {
		t.Fatalf("getTable failed: %v", err)
	}
	if len(tbl.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(tbl.Rows))
	}
}

func TestDocument_InsertImage(t *testing.T) {
	// create dummy image file
	tmpImg := "test_image.png"
	os.WriteFile(tmpImg, []byte("dummy image data"), 0644)
	defer os.Remove(tmpImg)

	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	err := doc.InsertImage(tmpImg, 100, 100)
	if err != nil {
		t.Fatalf("InsertImage failed: %v", err)
	}

	if len(doc.media) != 1 {
		t.Errorf("Expected 1 media item, got %d", len(doc.media))
	}
}
