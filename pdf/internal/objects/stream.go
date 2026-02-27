package objects

import (
	"fmt"
	"io"
)

// Stream represents a PDF stream object.
type Stream struct {
	Dict Dictionary
	Data []byte
}

func (s Stream) String() string {
	if s.Dict == nil {
		s.Dict = make(Dictionary)
	}
	s.Dict["Length"] = Integer(len(s.Data))
	return fmt.Sprintf("%s\nstream\n%s\nendstream", s.Dict.String(), string(s.Data))
}

func (s Stream) WriteTo(w io.Writer) (int64, error) {
	if s.Dict == nil {
		s.Dict = make(Dictionary)
	}
	s.Dict["Length"] = Integer(len(s.Data))

	var total int64
	n64, err := s.Dict.WriteTo(w)
	total += n64
	if err != nil {
		return total, err
	}

	n, err := fmt.Fprint(w, "\nstream\n")
	total += int64(n)
	if err != nil {
		return total, err
	}

	n, err = w.Write(s.Data)
	total += int64(n)
	if err != nil {
		return total, err
	}

	n, err = fmt.Fprint(w, "\nendstream")
	total += int64(n)
	return total, err
}
