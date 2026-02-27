package objects

import (
	"fmt"
	"io"
)

// Reference is a PDF reference object.
type Reference struct {
	Number     int
	Generation int
}

func (r Reference) String() string {
	return fmt.Sprintf("%d %d R", r.Number, r.Generation)
}

func (r Reference) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, r.String())
	return int64(n), err
}
