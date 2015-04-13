/*
Copyright 2015 Doubledutch

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mux

import (
	"errors"
	"net"
)

const (
	// ErrType for errors
	ErrType uint8 = iota
	// LogType for logs
	LogType
	// SignalType for signals
	SignalType
)

// Server wraps GobConn, adding Done
type Server struct {
	*GobConn
}

// NewServer returns a new server
func NewServer(conn net.Conn) *Server {
	gc := NewGobConn(conn)

	return &Server{
		GobConn: gc,
	}
}

// Done sends err to client. This marks the end of the server's work
// The server should not send further, the client may not receive it.
func (s *Server) Done(err error) {
	var errStr string
	if err == nil {
		errStr = ""
	} else {
		errStr = err.Error()
	}
	s.Send(ErrType, errStr)
}

// Client wraps GobConn, adding Wait
type Client struct {
	*GobConn
	errCh chan string
}

// NewClient returns a new client
func NewClient(conn net.Conn) *Client {
	gc := NewGobConn(conn)
	errCh := make(chan string, 1)

	errR := StringReceiver{
		dec: NewDecoder(),
		ch:  errCh,
	}

	gc.Receive(ErrType, errR)

	return &Client{
		GobConn: gc,
		errCh:   errCh,
	}
}

// Wait waits for an error from Server then closes the connection.
// When this returns, the server is done sending.
func (c *Client) Wait() error {
	errStr := <-c.errCh

	c.Shutdown()

	if errStr == "" {
		return nil
	}
	return errors.New(errStr)
}