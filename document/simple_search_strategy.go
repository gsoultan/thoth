package document

import (
	"context"
	"strings"
)

// SimpleSearchStrategy implements a basic text search.
type SimpleSearchStrategy struct{}

// NewSimpleSearchStrategy creates a new SimpleSearchStrategy.
func NewSimpleSearchStrategy() SearchStrategy {
	return &SimpleSearchStrategy{}
}

// Execute performs a simple substring search.
func (s *SimpleSearchStrategy) Execute(ctx context.Context, content string, patterns []string) ([]SearchResult, error) {
	var results []SearchResult
	for _, p := range patterns {
		if strings.Contains(content, p) {
			results = append(results, SearchResult{
				Keyword: p,
			})
		}
	}
	return results, nil
}
