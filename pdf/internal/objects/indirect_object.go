package objects

import (
	"fmt"
	"io"
)

// IndirectObject is a PDF indirect object.
type IndirectObject struct {
	Number     int
	Generation int
	Data       Object
}

func (io *IndirectObject) String() string {
	return fmt.Sprintf("%d %d obj\n%s\nendobj", io.Number, io.Generation, io.Data.String())
}

func (io *IndirectObject) WriteTo(w io.Writer) (int64, error) {
	return io.WriteEncrypted(w, nil, 0, 0)
}

func (io *IndirectObject) WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error) {
	var total int64
	n, err := fmt.Fprintf(w, "%d %d obj\n", io.Number, io.Generation)
	total += int64(n)
	if err != nil {
		return total, err
	}
	n64, err := io.Data.WriteEncrypted(w, ec, io.Number, io.Generation)
	total += n64
	if err != nil {
		return total, err
	}
	n, err = fmt.Fprint(w, "\nendobj")
	total += int64(n)
	return total, err
}
