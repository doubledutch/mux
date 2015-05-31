package mux

import (
	"errors"
	"time"

	"github.com/doubledutch/lager"
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
	// ErrInvalidLager defines an error for an invalid Config Lager value
	ErrInvalidLager = errors.New("Invalid Lager")
)

// Config configures a Server or Client
type Config struct {
	// Timeout for receiving frames
	Timeout time.Duration

	// Lager is used to control the log destination
	Lager lager.Lager
}

// Verify validates the config
func (c *Config) Verify() error {
	if c.Timeout == 0 {
		return ErrInvalidTimeout
	}

	if c.Lager == nil {
		return ErrInvalidLager
	}

	return nil
}

// DefaultConfig creates config with default settings
func DefaultConfig() *Config {
	return &Config{
		Timeout: 100 * time.Millisecond,
		Lager:   lager.NewLogLager(nil),
	}
}

// Server is a server that uses GobConn
type Server struct {
	Conn
}

// NewServer creates a new GobServer
func NewServer(conn Conn) (*Server, error) {
	return &Server{
		Conn: conn,
	}, nil
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

// Client is a client that uses GobConn
type Client struct {
	Conn
	errCh chan string
}

// NewClient returns a new GobClient
func NewClient(conn Conn) (*Client, error) {
	errCh := make(chan string, 1)
	errR := StringReceiver{
		dec: conn.Pool().NewBufferDecoder(),
		ch:  errCh,
	}

	conn.Receive(ErrType, errR)

	return &Client{
		Conn:  conn,
		errCh: errCh,
	}, nil
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
