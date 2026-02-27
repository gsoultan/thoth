package document

// SearchResult contains the information about a found keyword.
type SearchResult struct {
	Keyword  string
	Location string // e.g., "Page 1", "Sheet1!A1", etc.
	Index    int    // Character index or sequence number
}
