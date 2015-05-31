package gob

import (
	"bytes"
	"encoding/gob"
	"io"
	"net"

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

func (p *Pool) NewReceiver(ch interface{}) mux.Receiver {
	return mux.NewReceiver(ch, p)
}

func (p *Pool) NewServer(conn net.Conn, config *mux.Config) (mux.Server, error) {
	return NewServer(conn, config)
}

func (p *Pool) NewClient(conn net.Conn, config *mux.Config) (mux.Client, error) {
	return NewClient(conn, config)
}
