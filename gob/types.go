package gob

import (
	"net"

	"github.com/doubledutch/mux"
)

// NewDefaultServer creates a *mux.Client using gob encoding
func NewDefaultServer(conn net.Conn) (mux.Server, error) {
	gc, err := NewDefaultConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewServer(gc)
}

// NewDefaultClient creates a *mux.Client using gob encoding
func NewDefaultClient(conn net.Conn) (mux.Client, error) {
	gc, err := NewDefaultConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewClient(gc)
}
