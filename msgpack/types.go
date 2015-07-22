package msgpack

import (
	"net"

	"github.com/doubledutch/mux"
)

// NewDefaultServer creates a default mux.Server using msgpack encoding
func NewDefaultServer(conn net.Conn) (mux.Server, error) {
	gc, err := NewDefaultConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewServer(gc)
}

// NewDefaultClient creates a default mux.Client using msgpack encoding
func NewDefaultClient(conn net.Conn) (mux.Client, error) {
	gc, err := NewDefaultConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewClient(gc)
}

// NewClient creates a mux.Client using msgpack encoding
func NewClient(conn net.Conn, config *mux.Config) (mux.Client, error) {
	gc, err := NewConn(conn, config)
	if err != nil {
		return nil, err
	}

	return mux.NewClient(gc)
}

// NewServer creates a mux.Server using msgpack encoding
func NewServer(conn net.Conn, config *mux.Config) (mux.Server, error) {
	gc, err := NewConn(conn, config)
	if err != nil {
		return nil, err
	}

	return mux.NewServer(gc)
}
