package document

import "errors"

// Common errors for document processing.
var (
	ErrUnsupportedFormat = errors.New("unsupported document format")
	ErrInvalidFormat     = errors.New("invalid document format")
	ErrEncryptedDocument = errors.New("encrypted document not supported")
	ErrSheetNotFound     = errors.New("sheet not found")
	ErrDocumentNotLoaded = errors.New("document not loaded")
	ErrSaveFailed        = errors.New("save failed")
	ErrOpenFailed        = errors.New("open failed")
)
