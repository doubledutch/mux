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
}

// NewGobConn returns a new gob connection using a ReadWriter
func NewGobConn(conn net.Conn) *GobConn {
	return &GobConn{
		conn: conn,

		sendEnc:  NewEncoder(),
		sendLock: sync.Mutex{},

		dec: gob.NewDecoder(conn),
		enc: gob.NewEncoder(conn),

		Receivers:  make(map[uint8]Receiver),
		ShutdownCh: make(chan struct{}),

		timeout: 100 * time.Millisecond,
	}
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
		r, ok := c.Receivers[frame.Type]
		if !ok {
			log.Printf("[WARN] dropping frame %d\n", frame.Type)
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
