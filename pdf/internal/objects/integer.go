package objects

import (
	"fmt"
	"io"
)

// Integer is a PDF integer object.
type Integer int

func (i Integer) String() string {
	return fmt.Sprintf("%d", i)
}

func (i Integer) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, i.String())
	return int64(n), err
}

func (i Integer) WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error) {
	return i.WriteTo(w)
}
