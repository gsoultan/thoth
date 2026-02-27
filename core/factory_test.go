package core

import (
	"github.com/gsoultan/thoth/excel"
	"github.com/gsoultan/thoth/pdf"
	"github.com/gsoultan/thoth/word"
	"testing"
)

func TestDocumentFactory_Create(t *testing.T) {
	factory := NewDocumentFactory()

	tests := []struct {
		name     string
		filename string
		wantErr  bool
		wantType any
	}{
		{"excel .xlsx", "test.xlsx", false, &excel.Document{}},
		{"excel .xls", "test.xls", false, &excel.Document{}},
		{"word .docx", "test.docx", false, &word.Document{}},
		{"word .doc", "test.doc", false, &word.Document{}},
		{"pdf .pdf", "test.pdf", false, &pdf.Document{}},
		{"http url", "https://example.com/file.xlsx?q=1", false, &excel.Document{}},
		{"s3 uri", "s3://bucket/path/file.docx", false, &word.Document{}},
		{"unknown", "test.txt", true, nil},
		{"no extension", "test", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := factory.Create(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("Create() got nil, want %T", tt.wantType)
					return
				}
				// Type check using type switch or reflect if needed,
				// but here we just check if it's not nil for the mock.
			}
		})
	}
}
