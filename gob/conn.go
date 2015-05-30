package gob

import (
	"net"

	"github.com/doubledutch/mux"
)

func NewDefaultNetConn(conn net.Conn) (mux.Conn, error) {
	return mux.NewNetConn(conn, new(Pool), mux.DefaultConfig())
}

func NewNetConn(conn net.Conn, config *mux.Config) (mux.Conn, error) {
	return mux.NewNetConn(conn, new(Pool), config)
}
