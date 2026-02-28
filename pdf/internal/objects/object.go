package objects

import "io"

// Object is the base interface for all PDF objects.
// Primitives include Name, Integer, String, Array, and Dictionary.
type Object interface {
	// String returns the PDF representation of the object.
	String() string
	// WriteTo writes the PDF representation to the writer.
	WriteTo(w io.Writer) (int64, error)
	// WriteEncrypted writes the PDF representation to the writer, encrypted if needed.
	WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error)
}
