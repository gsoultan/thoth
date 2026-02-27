package document

import (
	"time"
)

// Metadata represents document-level information.
type Metadata struct {
	Title       string
	Author      string
	Subject     string
	Description string
	Keywords    []string
	Created     time.Time
	Modified    time.Time
	Generator   string
}
