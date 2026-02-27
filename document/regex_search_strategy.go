package document

import (
	"context"
	"regexp"
)

// RegexSearchStrategy implements a regex-based search.
type RegexSearchStrategy struct{}

// NewRegexSearchStrategy creates a new RegexSearchStrategy.
func NewRegexSearchStrategy() SearchStrategy {
	return &RegexSearchStrategy{}
}

// Execute performs a regex-based search.
func (s *RegexSearchStrategy) Execute(ctx context.Context, content string, patterns []string) ([]SearchResult, error) {
	var results []SearchResult
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		if re.MatchString(content) {
			results = append(results, SearchResult{
				Keyword: p,
			})
		}
	}
	return results, nil
}
