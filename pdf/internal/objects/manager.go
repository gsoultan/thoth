package objects

import (
	"cmp"
	"fmt"
	"io"
	"maps"
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
func (m *ObjectManager) Write(w io.Writer, catalogRef Reference, infoRef *Reference, ec *EncryptionContext, encryptRef *Reference, fileID []byte) error {
	type offsetGen struct {
		offset int64
		gen    int
	}
	offsets := make(map[int]offsetGen)
	var currentOffset int64

	// Write header
	n, _ := fmt.Fprintf(w, "%%PDF-2.0\n")
	currentOffset += int64(n)
	n, _ = fmt.Fprintf(w, "%%\xe2\xe3\xcf\xd3\n")
	currentOffset += int64(n)

	// Sort objects by ID for consistency
	slices.SortFunc(m.Objects, func(a, b IndirectObject) int {
		return cmp.Compare(a.Number, b.Number)
	})

	// Write objects
	for _, obj := range m.Objects {
		offsets[obj.Number] = offsetGen{offset: currentOffset, gen: obj.Generation}
		var n64 int64
		var err error
		if encryptRef != nil && obj.Number == encryptRef.Number {
			// Encrypt dictionary itself is NOT encrypted
			n64, err = obj.WriteEncrypted(w, nil, obj.Number, obj.Generation)
		} else {
			n64, err = obj.WriteEncrypted(w, ec, obj.Number, obj.Generation)
		}
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
	if len(offsets) > 0 {
		maxID = slices.Max(slices.Collect(maps.Keys(offsets)))
	}

	fmt.Fprintf(w, "xref\n0 %d\n0000000000 65535 f\r\n", maxID+1)
	for i := range maxID {
		objID := i + 1
		if og, ok := offsets[objID]; ok {
			fmt.Fprintf(w, "%010d %05d n\r\n", og.offset, og.gen)
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
	if encryptRef != nil {
		trailer["Encrypt"] = *encryptRef
	}
	if fileID != nil {
		id := PDFString(fileID)
		trailer["ID"] = Array{id, id}
	}

	fmt.Fprintf(w, "trailer\n")
	_, err := trailer.WriteTo(w)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\nstartxref\n%d\n%%%%EOF\n", startXRef)
	return nil
}
