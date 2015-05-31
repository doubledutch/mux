package gob

import (
	"bytes"
	"encoding/gob"

	"github.com/doubledutch/mux"
)

// BufferEncoder is used to encode values to bytes
type BufferEncoder struct {
	*bytes.Buffer
	*gob.Encoder
}

// NewBufferEncoder creates a encoder
func NewBufferEncoder(buf *bytes.Buffer) mux.BufferEncoder {
	enc := gob.NewEncoder(buf)

	return &BufferEncoder{
		Buffer:  buf,
		Encoder: enc,
	}
}

// BufferDecoder is used to decode bytes to values
type BufferDecoder struct {
	*bytes.Buffer
	*gob.Decoder
}

// NewBufferDecoder creates a decoder
func NewBufferDecoder(buf *bytes.Buffer) mux.BufferDecoder {
	dec := gob.NewDecoder(buf)

	return &BufferDecoder{
		Buffer:  buf,
		Decoder: dec,
	}
}
