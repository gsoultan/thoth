package document

// Metadatable defines an interface for documents that support metadata operations.
type Metadatable interface {
	GetMetadata() (Metadata, error)
	SetMetadata(metadata Metadata) error
}
