package json

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/doubledutch/mux"
)

type Pool struct {
	BasePool
}

func (p *Pool) NewBufferEncoder() mux.BufferEncoder {
	return NewBufferEncoder(new(bytes.Buffer))
}

func (p *Pool) NewBufferDecoder() mux.BufferDecoder {
	return NewBufferDecoder(new(bytes.Buffer))
}

func (p *Pool) NewEncoder(w io.Writer) mux.Encoder {
	return json.NewEncoder(w)
}

func (p *Pool) NewDecoder(r io.Reader) mux.Decoder {
	return json.NewDecoder(r)
}
