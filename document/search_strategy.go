package document

import (
	"context"
)

// SearchStrategy defines the interface for different search algorithms.
type SearchStrategy interface {
	Execute(ctx context.Context, content string, patterns []string) ([]SearchResult, error)
}
