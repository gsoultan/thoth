package document

import (
	"context"
	"io"
)

// Document is the common interface for processing Excel, Word, and PDF files.
// It provides methods for opening, reading, searching, modifying, and saving documents.
type Document interface {
	Metadatable

	// Open loads a document from a reader.
	Open(ctx context.Context, reader io.Reader) error

	// ReadContent returns the text content of the document.
	ReadContent() (string, error)

	// Search finds keywords in the document.
	Search(keywords []string) ([]SearchResult, error)

	// Replace replaces keywords with new values throughout the document.
	Replace(replacements map[string]string) error

	// SetPageSettings configures the document's layout.
	SetPageSettings(settings PageSettings) error

	// Save writes the modified document to the provided writer.
	Save(ctx context.Context, writer io.Writer) error

	// Export saves the document using the configured storage system.
	Export(uri string) error

	// SetPassword sets the document password for encryption.
	SetPassword(password string) error

	// Close releases any resources used by the document.
	Close() error
}
