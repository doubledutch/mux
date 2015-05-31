package json

import (
	"bytes"
	"encoding/json"

	"github.com/doubledutch/mux"
)

// BufferEncoder is used to encode values to bytes
type BufferEncoder struct {
	*bytes.Buffer
	mux.Encoder
}

// NewBufferEncoder creates a encoder
func NewBufferEncoder(buf *bytes.Buffer) mux.BufferEncoder {
	enc := json.NewEncoder(buf)

	return &BufferEncoder{
		Buffer:  buf,
		Encoder: enc,
	}
}

// BufferDecoder is used to decode bytes to values
type BufferDecoder struct {
	*bytes.Buffer
	mux.Decoder
}

// NewBufferDecoder creates a decoder
func NewBufferDecoder(buf *bytes.Buffer) mux.BufferDecoder {
	dec := json.NewDecoder(buf)

	return &BufferDecoder{
		Buffer:  buf,
		Decoder: dec,
	}
}
