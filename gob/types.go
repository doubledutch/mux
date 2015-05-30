package gob

import (
	"net"

	"github.com/doubledutch/mux"
)

func NewDefaultServer(conn net.Conn) (*mux.Server, error) {
	gc, err := NewDefaultNetConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewServer(gc)
}

func NewDefaultClient(conn net.Conn) (*mux.Client, error) {
	gc, err := NewDefaultNetConn(conn)
	if err != nil {
		return nil, err
	}

	return mux.NewClient(gc)
}
