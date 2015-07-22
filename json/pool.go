package json

import (
	"bytes"
	"encoding/json"
	"io"
	"net"

	"github.com/doubledutch/mux"
)

// Pool provides objects for encoding and decoding with json
type Pool struct {
}

// NewBufferEncoder creates a buffer encoder using json
func (p *Pool) NewBufferEncoder() mux.BufferEncoder {
	return NewBufferEncoder(new(bytes.Buffer))
}

// NewBufferDecoder creates a buffer decoder using json
func (p *Pool) NewBufferDecoder() mux.BufferDecoder {
	return NewBufferDecoder(new(bytes.Buffer))
}

// NewEncoder creates a encoder using json
func (p *Pool) NewEncoder(w io.Writer) mux.Encoder {
	return json.NewEncoder(w)
}

// NewDecoder creates a decoder using json
func (p *Pool) NewDecoder(r io.Reader) mux.Decoder {
	return json.NewDecoder(r)
}

// NewReceiver creates a new Receiver using json
func (p *Pool) NewReceiver(ch interface{}) mux.Receiver {
	return mux.NewReceiver(ch, p)
}

// NewServer creates a new Server using json
func (p *Pool) NewServer(conn net.Conn, config *mux.Config) (mux.Server, error) {
	return NewServer(conn, config)
}

// NewClient creates a new Client using json
func (p *Pool) NewClient(conn net.Conn, config *mux.Config) (mux.Client, error) {
	return NewClient(conn, config)
}
