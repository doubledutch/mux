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
	"io"
	"net"
	"os"
	"time"
)

const (
	// ErrType for errors
	ErrType uint8 = iota
	// LogType for logs
	LogType
	// SignalType for signals
	SignalType
)

var (
	// ErrInvalidTimeout defines an error for an invalid Config Timeout value
	ErrInvalidTimeout = errors.New("Invalid Config Timeout")
	// ErrInvalidLogOutput defines an error for an invalid Config LogOutput value
	ErrInvalidLogOutput = errors.New("Invalid LogOutput")
)

// Config configures a Server or Client
type Config struct {
	// Timeout for receiving frames
	Timeout time.Duration

	// LogOutput is used to control the log destination
	LogOutput io.Writer
}

// Verify validates the config
func (c *Config) Verify() error {
	if c.Timeout == 0 {
		return ErrInvalidTimeout
	}

	if c.LogOutput == nil {
		return ErrInvalidLogOutput
	}

	return nil
}

// DefaultConfig creates config with default settings
func DefaultConfig() *Config {
	return &Config{
		Timeout:   100 * time.Millisecond,
		LogOutput: os.Stderr,
	}
}

// Server wraps Conn, adding Done
type Server interface {
	Conn
	Done(err error)
}

// GobServer is a server that uses GobConn
type GobServer struct {
	Conn
}

// NewDefaultServer creates a new Server with default configuration
func NewDefaultServer(conn net.Conn) (Server, error) {
	return NewGobServer(conn, DefaultConfig())
}

// NewGobServer creates a new GobServer
func NewGobServer(conn net.Conn, config *Config) (Server, error) {
	gc, err := NewGobConn(conn, config)
	if err != nil {
		return nil, err
	}
	return &GobServer{
		Conn: gc,
	}, nil
}

// Done sends err to client. This marks the end of the server's work
// The server should not send further, the client may not receive it.
func (s *GobServer) Done(err error) {
	var errStr string
	if err == nil {
		errStr = ""
	} else {
		errStr = err.Error()
	}
	s.Send(ErrType, errStr)
}

// Client wraps Conn, adding Wait
type Client interface {
	Conn
	Wait() error
}

// GobClient is a client that uses GobConn
type GobClient struct {
	Conn
	errCh chan string
}

// NewDefaultClient creates a new client with default configuration
func NewDefaultClient(conn net.Conn) (Client, error) {
	return NewGobClient(conn, DefaultConfig())
}

// NewGobClient returns a new GobClient
func NewGobClient(conn net.Conn, config *Config) (Client, error) {
	gc, err := NewGobConn(conn, config)
	if err != nil {
		return nil, err
	}

	errCh := make(chan string, 1)
	errR := StringReceiver{
		dec: NewDecoder(),
		ch:  errCh,
	}

	gc.Receive(ErrType, errR)

	return &GobClient{
		Conn:  gc,
		errCh: errCh,
	}, nil
}

// Wait waits for an error from Server then closes the connection.
// When this returns, the server is done sending.
func (c *GobClient) Wait() error {
	errStr := <-c.errCh

	c.Shutdown()

	if errStr == "" {
		return nil
	}
	return errors.New(errStr)
}
