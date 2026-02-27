package objects

import (
	"fmt"
	"io"
	"slices"
)

// ObjectManager handles PDF objects and their generation.
type ObjectManager struct {
	Objects []IndirectObject
	nextID  int
}

// NewObjectManager creates a new ObjectManager.
func NewObjectManager() *ObjectManager {
	return &ObjectManager{
		Objects: make([]IndirectObject, 0),
		nextID:  1,
	}
}

// AddObject adds an object and returns its reference.
func (m *ObjectManager) AddObject(obj Object) Reference {
	id := m.nextID
	m.nextID++
	iobj := IndirectObject{
		Number:     id,
		Generation: 0,
		Data:       obj,
	}
	m.Objects = append(m.Objects, iobj)
	return Reference{Number: id, Generation: 0}
}

// Write writes the PDF to the writer.
func (m *ObjectManager) Write(w io.Writer, catalogRef Reference, infoRef *Reference) error {
	offsets := make(map[int]int64)
	var currentOffset int64

	// Write header
	n, _ := fmt.Fprintf(w, "%%PDF-1.7\n")
	currentOffset += int64(n)

	// Sort objects by ID for consistency
	slices.SortFunc(m.Objects, func(a, b IndirectObject) int {
		return a.Number - b.Number
	})

	// Write objects
	for _, obj := range m.Objects {
		offsets[obj.Number] = currentOffset
		n64, err := obj.WriteTo(w)
		if err != nil {
			return err
		}
		currentOffset += n64
		n, err := fmt.Fprint(w, "\n")
		if err != nil {
			return err
		}
		currentOffset += int64(n)
	}

	// Write XRef table
	startXRef := currentOffset
	maxID := 0
	for id := range offsets {
		if id > maxID {
			maxID = id
		}
	}

	fmt.Fprintf(w, "xref\n0 %d\n0000000000 65535 f\r\n", maxID+1)
	for i := 1; i <= maxID; i++ {
		if offset, ok := offsets[i]; ok {
			fmt.Fprintf(w, "%010d 00000 n\r\n", offset)
		} else {
			fmt.Fprintf(w, "0000000000 65535 f\r\n")
		}
	}

	// Write Trailer
	trailer := Dictionary{
		"Size": Integer(maxID + 1),
		"Root": catalogRef,
	}
	if infoRef != nil {
		trailer["Info"] = *infoRef
	}

	fmt.Fprintf(w, "trailer\n")
	_, err := trailer.WriteTo(w)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\nstartxref\n%d\n%%%%EOF\n", startXRef)
	return nil
}
