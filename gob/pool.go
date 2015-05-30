package gob

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/doubledutch/mux"
)

type Pool struct {
}

func (p *Pool) NewBufferEncoder() mux.BufferEncoder {
	return NewBufferEncoder(new(bytes.Buffer))
}

func (p *Pool) NewBufferDecoder() mux.BufferDecoder {
	return NewBufferDecoder(new(bytes.Buffer))
}

func (p *Pool) NewEncoder(w io.Writer) mux.Encoder {
	return gob.NewEncoder(w)
}

func (p *Pool) NewDecoder(r io.Reader) mux.Decoder {
	return gob.NewDecoder(r)
}
