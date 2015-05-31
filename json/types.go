package json

import (
	"net"

	"github.com/doubledutch/mux"
)

// NewDefaultServer creates a mux.Client using json encoding
func NewDefaultServer(conn net.Conn) (mux.Server, error) {
	gc, err := NewDefaultConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewServer(gc)
}

// NewDefaultClient creates a mux.Client using json encoding
func NewDefaultClient(conn net.Conn) (mux.Client, error) {
	gc, err := NewDefaultConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewClient(gc)
}

func NewClient(conn net.Conn, config *mux.Config) (mux.Client, error) {
	gc, err := NewConn(conn, config)
	if err != nil {
		return nil, err
	}

	return mux.NewClient(gc)
}

func NewServer(conn net.Conn, config *mux.Config) (mux.Server, error) {
	gc, err := NewConn(conn, config)
	if err != nil {
		return nil, err
	}

	return mux.NewServer(gc)
}
