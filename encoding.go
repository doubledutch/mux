package mux

import "io"

// Encoder encodes
type Encoder interface {
	Encode(v interface{}) error
}

// Decoder decodes
type Decoder interface {
	Decode(v interface{}) error
}

// Buffer is an interface for buffering
type Buffer interface {
	Write(b []byte) (int, error)
	Bytes() []byte
	Len() int
	Reset()
}

// BufferEncoder is a Buffer and Encoder
type BufferEncoder interface {
	Buffer
	Encoder
}

// BufferDecoder is a Buffer and Decoder
type BufferDecoder interface {
	Buffer
	Decoder
}

// Pool is an interface for interacting with encoding implementations
type Pool interface {
	NewBufferEncoder() BufferEncoder
	NewBufferDecoder() BufferDecoder
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}
