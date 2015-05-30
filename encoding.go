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

type Buffer interface {
	Write(b []byte) (int, error)
	Bytes() []byte
	Len() int
	Reset()
}

type BufferEncoder interface {
	Buffer
	Encoder
}

type BufferDecoder interface {
	Buffer
	Decoder
}

type Pool interface {
	NewBufferEncoder() BufferEncoder
	NewBufferDecoder() BufferDecoder
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}
