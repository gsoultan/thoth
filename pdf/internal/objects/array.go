package objects

import (
	"io"
	"strings"
)

// Array is a PDF array object.
type Array []Object

func (a Array) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range a {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}

func (a Array) WriteTo(w io.Writer) (int64, error) {
	var total int64
	n, err := w.Write([]byte("["))
	total += int64(n)
	if err != nil {
		return total, err
	}
	for i, v := range a {
		if i > 0 {
			n, err = w.Write([]byte(" "))
			total += int64(n)
			if err != nil {
				return total, err
			}
		}
		if v == nil {
			n, err = w.Write([]byte("null"))
			total += int64(n)
			if err != nil {
				return total, err
			}
		} else {
			n64, err := v.WriteTo(w)
			total += n64
			if err != nil {
				return total, err
			}
		}
	}
	n, err = w.Write([]byte("]"))
	total += int64(n)
	return total, err
}

func (a Array) WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error) {
	var total int64
	n, err := w.Write([]byte("["))
	total += int64(n)
	if err != nil {
		return total, err
	}
	for i, v := range a {
		if i > 0 {
			n, err = w.Write([]byte(" "))
			total += int64(n)
			if err != nil {
				return total, err
			}
		}
		if v == nil {
			n, err = w.Write([]byte("null"))
			total += int64(n)
			if err != nil {
				return total, err
			}
		} else {
			n64, err := v.WriteEncrypted(w, ec, objNum, objGen)
			total += n64
			if err != nil {
				return total, err
			}
		}
	}
	n, err = w.Write([]byte("]"))
	total += int64(n)
	return total, err
}
