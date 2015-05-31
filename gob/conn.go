package gob

import (
	"net"

	"github.com/doubledutch/mux"
)

// NewDefaultConn creates a new mux.Conn using gob encoding with default configuration
func NewDefaultConn(conn net.Conn) (mux.Conn, error) {
	return mux.NewConn(conn, new(Pool), mux.DefaultConfig())
}

// NewConn creates a new mux.Conn using gob encoding
func NewConn(conn net.Conn, config *mux.Config) (mux.Conn, error) {
	return mux.NewConn(conn, new(Pool), config)
}
