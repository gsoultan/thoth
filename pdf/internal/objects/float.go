package objects

import (
	"fmt"
	"io"
)

// Float represents a PDF real number object.
type Float float64

func (f Float) String() string {
	return fmt.Sprintf("%.2f", f)
}

func (f Float) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, f.String())
	return int64(n), err
}

func (f Float) WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error) {
	return f.WriteTo(w)
}
