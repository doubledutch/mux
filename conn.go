package mux

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Conn defines a mux connection
type Conn interface {
	Receive(t uint8, r Receiver)
	Send(t uint8, e interface{}) error
	Recv()
	Shutdown()
	IsShutdown() chan struct{}
}

// Encoder is used to encode values to bytes
type Encoder struct {
	*bytes.Buffer
	*gob.Encoder
}

// NewEncoder creates a encoder
func NewEncoder() *Encoder {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	return &Encoder{
		Buffer:  buf,
		Encoder: enc,
	}
}

// Decoder is used to decode bytes to values
type Decoder struct {
	*bytes.Buffer
	*gob.Decoder
}

// NewDecoder creates a decoder
func NewDecoder() *Decoder {
	buf := new(bytes.Buffer)
	dec := gob.NewDecoder(buf)

	return &Decoder{
		Buffer:  buf,
		Decoder: dec,
	}
}

// Frame represents transport
type Frame struct {
	Type uint8
	Data []byte
}

// GobConn wraps net.Conn which communicates using gob
type GobConn struct {
	// store the net.Conn to SetDeadlines
	conn net.Conn

	// used to encode data into frames
	sendEnc  *Encoder
	sendLock sync.Mutex

	// encode and decode conn
	enc *gob.Encoder
	dec *gob.Decoder

	// Store receivers for Frames
	Receivers map[uint8]Receiver

	// allow of users and ourselves to listen for shutdown
	ShutdownCh chan struct{}
	isShutdown bool

	// timeout for receiving frames
	timeout time.Duration

	logger *log.Logger
}

// NewDefaultGobConn returns a new GobConn using net.Conn and DefaultConfig
func NewDefaultGobConn(conn net.Conn) (*GobConn, error) {
	return NewGobConn(conn, DefaultConfig())
}

// NewGobConn creates a new GobConn using the specified conn and config
func NewGobConn(conn net.Conn, config *Config) (*GobConn, error) {
	if err := config.Verify(); err != nil {
		return nil, err
	}

	return &GobConn{
		conn: conn,

		sendEnc:  NewEncoder(),
		sendLock: sync.Mutex{},

		dec: gob.NewDecoder(conn),
		enc: gob.NewEncoder(conn),

		Receivers:  make(map[uint8]Receiver),
		ShutdownCh: make(chan struct{}),

		timeout: config.Timeout,
		logger:  log.New(config.LogOutput, "", log.LstdFlags),
	}, nil
}

// Send encodes a frame on conn using t and e
func (c *GobConn) Send(t uint8, e interface{}) error {
	// Single threaded through here
	c.sendLock.Lock()
	c.sendEnc.Encode(e)

	d := make([]byte, c.sendEnc.Len())
	copy(d, c.sendEnc.Bytes())
	c.sendEnc.Reset()
	c.sendLock.Unlock()

	f := Frame{
		Type: t,
		Data: d,
	}
	c.logger.Printf("[DEBUG] Sending frame: %v\n", f)

	return c.enc.Encode(f)
}

// Receiver defines an interface for receiving
type Receiver interface {
	Receive(b []byte) error
	Close() error
}

// Receive registers a receiver to receive t
func (c *GobConn) Receive(t uint8, r Receiver) {
	c.Receivers[t] = r
	c.logger.Printf("[DEBUG] Added receiver type %d\n", t)
}

// Recv listens for frames and sends them to a receiver
func (c *GobConn) Recv() {
	for {
		var frame Frame
		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		err := c.dec.Decode(&frame)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed") || strings.Contains(err.Error(), "reset by peer") {
				// This is the expected way for us to return
				return
			}
			if err, ok := err.(*net.OpError); ok && err.Timeout() {
				select {
				case <-c.ShutdownCh:
					return
				default: // Keep listening
					continue
				}
			} else {
				// Unexpected error
				return
			}
		}
		c.logger.Printf("[DEBUG] Received frame: %v\n", frame)
		r, ok := c.Receivers[frame.Type]
		if !ok {
			c.logger.Printf("[WARN] dropping frame %d\n", frame.Type)
			continue
		}
		err = r.Receive(frame.Data)
		if err != nil {
			// The handler returns an error.. what now?
		}
	}
}

// IsShutdown provides a way to listen for this connection to shutdown
func (c *GobConn) IsShutdown() chan struct{} {
	return c.ShutdownCh
}

// Shutdown closes the gob connection
func (c *GobConn) Shutdown() {
	if c.isShutdown {
		return
	}
	c.logger.Println("[INFO] Shutting down")
	c.isShutdown = true
	// Notify that we're shutdown
	close(c.ShutdownCh)

	// Let receivers clean themselves up
	for _, h := range c.Receivers {
		h.Close()
	}

	// We're done with conn
	c.conn.Close()
}
