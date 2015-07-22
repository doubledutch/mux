package msgpack

import (
	"bytes"

	"github.com/doubledutch/mux"
	"github.com/ugorji/go/codec"
)

// BufferEncoder is used to encode values to bytes
type BufferEncoder struct {
	*bytes.Buffer
	*codec.Encoder
}

// NewBufferEncoder creates a encoder
func NewBufferEncoder(buf *bytes.Buffer, mh *codec.MsgpackHandle) mux.BufferEncoder {
	enc := codec.NewEncoder(buf, mh)

	return &BufferEncoder{
		Buffer:  buf,
		Encoder: enc,
	}
}

// BufferDecoder is used to decode bytes to values
type BufferDecoder struct {
	*bytes.Buffer
	*codec.Decoder
}

// NewBufferDecoder creates a decoder
func NewBufferDecoder(buf *bytes.Buffer, mh *codec.MsgpackHandle) mux.BufferDecoder {
	dec := codec.NewDecoder(buf, mh)

	return &BufferDecoder{
		Buffer:  buf,
		Decoder: dec,
	}
}
