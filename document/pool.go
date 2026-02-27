package document

import (
	"bytes"
	"encoding/xml"
	"io"
	"sync"
)

// BufferPool is a pool of bytes.Buffer objects to reduce allocations.
var BufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// GetBuffer returns a bytes.Buffer from the pool.
func GetBuffer() *bytes.Buffer {
	return BufferPool.Get().(*bytes.Buffer)
}

// PutBuffer returns a bytes.Buffer to the pool, resetting it first.
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	BufferPool.Put(buf)
}

// XMLEncoderPool is a pool of xml.Encoder objects.
var XMLEncoderPool = sync.Pool{
	New: func() any {
		return xml.NewEncoder(io.Discard)
	},
}

// XMLDecoderPool is a pool of xml.Decoder objects.
var XMLDecoderPool = sync.Pool{
	New: func() any {
		return xml.NewDecoder(bytes.NewReader(nil))
	},
}
