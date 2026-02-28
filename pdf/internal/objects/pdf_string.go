package objects

import (
	"fmt"
	"io"
	"strings"
)

// PDFString is a PDF literal string object.
// Renamed from String to avoid conflict with string type.
type PDFString string

func (s PDFString) String() string {
	str := string(s)
	if isNonASCII(str) {
		// Prepend UTF-8 BOM for PDF 2.0 / common reader support
		str = "\xEF\xBB\xBF" + str
	}
	return "(" + EscapeString(str) + ")"
}

func isNonASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}
	return false
}

func (s PDFString) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, s.String())
	return int64(n), err
}

func (s PDFString) WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error) {
	if ec == nil {
		return s.WriteTo(w)
	}
	return ec.WriteEncrypted(w, []byte(s), objNum, objGen)
}

func EscapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}
