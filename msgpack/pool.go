package msgpack

import (
	"bytes"
	"io"
	"net"

	"github.com/doubledutch/mux"
	"github.com/ugorji/go/codec"
)

// Pool provides for encoding and decoding
type Pool struct {
	mh *codec.MsgpackHandle
}

// Initialize ensures the pool is initialized
func (p *Pool) Initialize() {
	if p.mh == nil {
		p.mh = new(codec.MsgpackHandle)
	}
}

// NewBufferEncoder creates a new BufferEncoder
func (p *Pool) NewBufferEncoder() mux.BufferEncoder {
	p.Initialize()

	return NewBufferEncoder(new(bytes.Buffer), p.mh)
}

// NewBufferDecoder creates a new BufferDecoder
func (p *Pool) NewBufferDecoder() mux.BufferDecoder {
	p.Initialize()

	return NewBufferDecoder(new(bytes.Buffer), p.mh)
}

// NewEncoder creates a new Encoder
func (p *Pool) NewEncoder(w io.Writer) mux.Encoder {
	p.Initialize()

	return codec.NewEncoder(w, p.mh)
}

// NewDecoder creates a new Decoder
func (p *Pool) NewDecoder(r io.Reader) mux.Decoder {
	p.Initialize()

	return codec.NewDecoder(r, p.mh)
}

// NewReceiver creates a new receiver
func (p *Pool) NewReceiver(ch interface{}) mux.Receiver {
	return mux.NewReceiver(ch, p)
}

// NewServer creates a new Server
func (p *Pool) NewServer(conn net.Conn, config *mux.Config) (mux.Server, error) {
	return NewServer(conn, config)
}

// NewClient creates a new Client
func (p *Pool) NewClient(conn net.Conn, config *mux.Config) (mux.Client, error) {
	return NewClient(conn, config)
}
