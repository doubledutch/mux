package gob

import (
	"bytes"
	"encoding/gob"
	"io"
	"net"

	"github.com/doubledutch/mux"
)

// Pool provides objects for encoding and decoding with gob
type Pool struct{}

// NewBufferEncoder creates a buffer encoder using gob
func (p *Pool) NewBufferEncoder() mux.BufferEncoder {
	return NewBufferEncoder(new(bytes.Buffer))
}

// NewBufferDecoder creates a buffer decoder using gob
func (p *Pool) NewBufferDecoder() mux.BufferDecoder {
	return NewBufferDecoder(new(bytes.Buffer))
}

// NewEncoder creates a encoder using gob
func (p *Pool) NewEncoder(w io.Writer) mux.Encoder {
	return gob.NewEncoder(w)
}

// NewDecoder creates a decoder using gob
func (p *Pool) NewDecoder(r io.Reader) mux.Decoder {
	return gob.NewDecoder(r)
}

// NewReceiver creates a new Receiver using gob
func (p *Pool) NewReceiver(ch interface{}) mux.Receiver {
	return mux.NewReceiver(ch, p)
}

// NewServer creates a new Server using gob
func (p *Pool) NewServer(conn net.Conn, config *mux.Config) (mux.Server, error) {
	return NewServer(conn, config)
}

// NewClient creates a new Client using gob
func (p *Pool) NewClient(conn net.Conn, config *mux.Config) (mux.Client, error) {
	return NewClient(conn, config)
}
