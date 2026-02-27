package core

// DocumentType represents the type of document.
type DocumentType int

const (
	// TypeUnknown represents an unknown document type.
	TypeUnknown DocumentType = iota
	// TypeExcel represents an Excel document type.
	TypeExcel
	// TypeWord represents a Word document type.
	TypeWord
	// TypePDF represents a PDF document type.
	TypePDF
)
