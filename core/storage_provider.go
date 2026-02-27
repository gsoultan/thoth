package core

import (
	"context"
	"io"
)

// StorageProvider defines the interface for opening and saving documents from different storage types.
type StorageProvider interface {
	Open(ctx context.Context, uri string) (io.ReadCloser, error)
	Save(ctx context.Context, uri string, reader io.Reader) error
}
