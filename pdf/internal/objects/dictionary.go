package objects

import (
	"fmt"
	"io"
	"strings"
)

// Dictionary is a PDF dictionary object.
type Dictionary map[string]Object

func (d Dictionary) String() string {
	var sb strings.Builder
	sb.WriteString("<<")
	for k, v := range d {
		sb.WriteString(fmt.Sprintf("/%s %s", k, v.String()))
	}
	sb.WriteString(">>")
	return sb.String()
}

func (d Dictionary) WriteTo(w io.Writer) (int64, error) {
	var total int64
	n, err := w.Write([]byte("<<"))
	total += int64(n)
	if err != nil {
		return total, err
	}
	for k, v := range d {
		n, err = fmt.Fprintf(w, "/%s ", k)
		total += int64(n)
		if err != nil {
			return total, err
		}
		n64, err := v.WriteTo(w)
		total += n64
		if err != nil {
			return total, err
		}
	}
	n, err = w.Write([]byte(">>"))
	total += int64(n)
	return total, err
}
