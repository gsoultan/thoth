package objects

import "io"

// Name represents a PDF name object.
type Name string

// String returns the name in PDF format (starting with /).
func (n Name) String() string {
	return "/" + string(n)
}

func (n Name) WriteTo(w io.Writer) (int64, error) {
	m, err := w.Write([]byte(n.String()))
	return int64(m), err
}

func (n Name) WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error) {
	return n.WriteTo(w)
}
