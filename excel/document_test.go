package excel

import (
	"github.com/gsoultan/thoth/document"
	"testing"
)

func TestDocument_Search(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	keywords := []string{"test", "hello"}
	results, err := doc.Search(keywords)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Since it's a mock/skeleton, it currently returns nil, nil.
	// In a real implementation, we would check the results.
	if results != nil {
		t.Errorf("Expected nil results for mock, got %v", results)
	}
}

func TestDocument_Replace(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	replacements := map[string]string{"old": "new"}
	err := doc.Replace(replacements)
	if err != nil {
		t.Errorf("Replace failed: %v", err)
	}
}

func TestDocument_SetCellStyle(t *testing.T) {
	doc := NewDocument().(*Document)
	ctx := t.Context()
	doc.SetContext(ctx)

	doc.addSheet("Sheet1")

	style := document.CellStyle{
		Bold:  true,
		Size:  12,
		Color: "FF0000",
	}
	err := doc.setCellStyle("Sheet1", "A1", style)
	if err != nil {
		t.Fatalf("SetCellStyle failed: %v", err)
	}

	if doc.state.styles == nil {
		t.Fatal("Expected styles to be initialized")
	}

	if len(doc.state.styles.Fonts.Items) < 2 {
		t.Errorf("Expected at least 2 fonts (default + new), got %d", len(doc.state.styles.Fonts.Items))
	}
}
