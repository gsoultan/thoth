package objects

import (
	"bytes"
	"io"
)

// Stream represents a PDF stream object.
type Stream struct {
	Dict Dictionary
	Data []byte
}

func (s Stream) String() string {
	// Fallback to WriteTo if possible, but String() should ideally not be used for binary data.
	var buf bytes.Buffer
	_, _ = s.WriteTo(&buf)
	return buf.String()
}

func (s Stream) WriteTo(w io.Writer) (int64, error) {
	return s.WriteEncrypted(w, nil, 0, 0)
}

func (s Stream) WriteEncrypted(w io.Writer, ec *EncryptionContext, objNum, objGen int) (int64, error) {
	if s.Dict == nil {
		s.Dict = make(Dictionary)
	}

	data := s.Data
	if ec != nil {
		data = ec.Encrypt(s.Data, objNum, objGen)
	}

	s.Dict["Length"] = Integer(len(data))

	var total int64
	n64, err := s.Dict.WriteEncrypted(w, ec, objNum, objGen)
	total += n64
	if err != nil {
		return total, err
	}

	n, err := w.Write([]byte("\nstream\n"))
	total += int64(n)
	if err != nil {
		return total, err
	}

	n, err = w.Write(data)
	total += int64(n)
	if err != nil {
		return total, err
	}

	n, err = w.Write([]byte("\nendstream"))
	total += int64(n)
	return total, err
}
