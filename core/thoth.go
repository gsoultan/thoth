package core

import (
	"context"
	"fmt"
	"io"

	"github.com/gsoultan/thoth/document"
	"github.com/rs/zerolog"
)

// Thoth provides high-level document processing operations.
// It acts as a central hub for document creation, storage, and processing.
type Thoth struct {
	factory *DocumentFactory
	storage StorageProvider
	logger  zerolog.Logger
}

// New creates a new Thoth instance with a default logger and factory.
// It uses a default ThothStorage which handles local files and basic URLs.
func New(logger zerolog.Logger) *Thoth {
	return &Thoth{
		factory: NewDocumentFactory(),
		storage: NewThothStorage(nil),
		logger:  logger,
	}
}

// WithS3Config sets the S3 configuration for Thoth by updating its storage provider.
// This allows Thoth to open documents from s3:// URIs.
func (t *Thoth) WithS3Config(config S3Config) *Thoth {
	t.storage = NewThothStorage(&config)
	return t
}

// WithStorage sets a custom storage provider for Thoth.
// This is useful for implementing custom storage backends (e.g., GCS, Azure Blob Storage).
func (t *Thoth) WithStorage(storage StorageProvider) *Thoth {
	t.storage = storage
	return t
}

// openDocument opens an existing document from a local path, HTTP/HTTPS URL, or S3 URI.
func (t *Thoth) openDocument(ctx context.Context, uri string) (document.Document, error) {
	t.logger.Info().Str("uri", uri).Msg("Opening document")

	doc, err := t.factory.Create(uri)
	if err != nil {
		return nil, fmt.Errorf("factory create: %w", err)
	}

	reader, err := t.storage.Open(ctx, uri)
	if err != nil {
		doc.Close()
		return nil, fmt.Errorf("storage open: %w", err)
	}
	defer reader.Close()

	if err := doc.Open(ctx, reader); err != nil {
		doc.Close()
		return nil, fmt.Errorf("document open: %w", err)
	}

	t.bind(ctx, doc)
	return doc, nil
}

// saveDocument saves a document to a local path or S3 URI using the configured storage provider.
func (t *Thoth) saveDocument(ctx context.Context, doc document.Document, uri string) error {
	t.logger.Info().Str("uri", uri).Msg("Saving document")

	pr, pw := io.Pipe()
	errChan := make(chan error, 1)

	go func() {
		defer pw.Close()
		if err := doc.Save(ctx, pw); err != nil {
			errChan <- fmt.Errorf("document save: %w", err)
			return
		}
		errChan <- nil
	}()

	if err := t.storage.Save(ctx, uri, pr); err != nil {
		pr.Close() // Ensure the pipe is closed if storage save fails
		return fmt.Errorf("storage save: %w", err)
	}

	return <-errChan
}

// Excel returns fluent entry points for Excel.
func (t *Thoth) Excel() ExcelFluent {
	return ExcelFluent{thoth: t}
}

// Word returns fluent entry points for Word.
func (t *Thoth) Word() WordFluent {
	return WordFluent{thoth: t}
}

// PDF returns fluent entry points for PDF.
func (t *Thoth) PDF() PDFFluent {
	return PDFFluent{thoth: t}
}

// createNewDocument creates a new, empty document instance in-memory based on the file extension.
func (t *Thoth) createNewDocument(ctx context.Context, uri string) (document.Document, error) {
	t.logger.Debug().Str("uri", uri).Msg("Creating new document")
	doc, err := t.factory.Create(uri)
	if err != nil {
		return nil, fmt.Errorf("factory create: %w", err)
	}
	t.bind(ctx, doc)
	return doc, nil
}

type contextBinder interface {
	SetContext(ctx context.Context)
}

type exportBinder interface {
	SetExportFunc(fn func(doc document.Document, uri string) error)
}

func (t *Thoth) bind(ctx context.Context, doc document.Document) {
	if b, ok := doc.(contextBinder); ok {
		b.SetContext(ctx)
	}
	if b, ok := doc.(exportBinder); ok {
		b.SetExportFunc(func(d document.Document, uri string) error {
			return t.saveDocument(ctx, d, uri)
		})
	}
}
