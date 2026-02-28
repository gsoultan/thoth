package objects

import (
	"io"
	"strconv"
)

// Boolean represents a PDF boolean object (true or false).
type Boolean bool

func (b Boolean) String() string {
	return strconv.FormatBool(bool(b))
}

func (b Boolean) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(b.String()))
	return int64(n), err
}

func (b Boolean) WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error) {
	return b.WriteTo(w)
}
