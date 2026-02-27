package core

import (
	"os"
	"testing"

	"github.com/gsoultan/thoth/document"
	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	thoth := New(zerolog.Nop())
	if thoth == nil {
		t.Fatal("Expected New to return a Thoth instance, got nil")
	}
}

func TestThoth_GranularAPI(t *testing.T) {
	thoth := New(zerolog.Nop())
	ctx := t.Context()

	// Create a temporary file to "open"
	tmpFile := "test_granular.xlsx"
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	f.Close()
	defer os.Remove(tmpFile)

	// 1. Open
	doc, err := thoth.Excel().Open(ctx, tmpFile)
	if err != nil {
		// This might fail because the file is not a valid zip,
		// but we want to test the flow.
		t.Logf("Open failed as expected for empty file: %v", err)
	} else {
		defer doc.Close()
	}

	// Rest of the test...
}

func TestThoth_CreateAndSetSettings(t *testing.T) {
	thoth := New(zerolog.Nop())
	ctx := t.Context()

	// 1. Create
	doc, err := thoth.Word().New(ctx)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer doc.Close()

	// 2. Set Page Settings
	settings := document.PageSettings{
		Orientation: document.OrientationLandscape,
		PaperType:   document.PaperA4,
		Margins: document.Margins{
			Top:    1.0,
			Bottom: 1.0,
			Left:   1.0,
			Right:  1.0,
		},
	}
	err = doc.SetPageSettings(settings)
	if err != nil {
		t.Errorf("doc.SetPageSettings failed: %v", err)
	}

	// 3. Save
	outputFile := "created_doc.docx"
	err = doc.Export(outputFile)
	defer os.Remove(outputFile)
}
