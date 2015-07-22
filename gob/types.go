package gob

import (
	"net"

	"github.com/doubledutch/mux"
)

// NewDefaultServer creates a default mux.Server using gob encoding
func NewDefaultServer(conn net.Conn) (mux.Server, error) {
	gc, err := NewDefaultConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewServer(gc)
}

// NewDefaultClient creates a default mux.Client using gob encoding
func NewDefaultClient(conn net.Conn) (mux.Client, error) {
	gc, err := NewDefaultConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewClient(gc)
}

// NewClient creates a new mux.Client using gob encoding
func NewClient(conn net.Conn, config *mux.Config) (mux.Client, error) {
	gc, err := NewConn(conn, config)
	if err != nil {
		return nil, err
	}

	return mux.NewClient(gc)
}

// NewServer creates a new mux.Server using gob encoding
func NewServer(conn net.Conn, config *mux.Config) (mux.Server, error) {
	gc, err := NewConn(conn, config)
	if err != nil {
		return nil, err
	}

	return mux.NewServer(gc)
}
