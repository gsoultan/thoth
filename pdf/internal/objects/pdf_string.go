package objects

import (
	"fmt"
	"io"
)

// PDFString is a PDF literal string object.
// Renamed from String to avoid conflict with string type.
type PDFString string

func (s PDFString) String() string {
	return fmt.Sprintf("(%s)", string(s))
}

func (s PDFString) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, s.String())
	return int64(n), err
}
