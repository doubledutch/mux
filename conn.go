package mux

// NetConn wraps net.Conn which communicates using gob
import (
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/doubledutch/lager"
)

// Conn defines a mux connection
type Conn interface {
	// Receive registers a receiver to receive t
	Receive(t uint8, r Receiver)
	// Send encodes a frame on conn using t and e
	Send(t uint8, e interface{}) error
	// Recv listens for frames and sends them to a receiver
	Recv()
	// Pool returns the pool used by the Conn
	Pool() Pool
	// Shutdown closes the gob connection
	Shutdown()
	// IsShutdown provides a way to listen for this connection to shutdown
	IsShutdown() chan struct{}
}

// Receiver defines an interface for receiving
type Receiver interface {
	Receive(b []byte) error
	Close() error
}

// conn is a Conn using net.Conn for communication
type conn struct {
	// store the net.Conn to SetDeadlines
	conn net.Conn

	// used to encode data into frames
	sendEnc  BufferEncoder
	sendLock sync.Mutex

	// encode and decode conn
	enc Encoder
	dec Decoder

	// Store receivers for Frames
	Receivers map[uint8]Receiver

	// allow of users and ourselves to listen for shutdown
	ShutdownCh chan struct{}
	isShutdown bool

	// timeout for receiving frames
	timeout time.Duration

	lgr lager.Lager

	pool Pool
}

// Frame represents transport
type Frame struct {
	Type uint8
	Data []byte
}

// NewConn creates a new NetConn using the specified conn and config
func NewConn(netConn net.Conn, pool Pool, config *Config) (Conn, error) {
	if err := config.Verify(); err != nil {
		return nil, err
	}

	return &conn{
		conn: netConn,

		sendEnc:  pool.NewBufferEncoder(),
		sendLock: sync.Mutex{},

		dec: pool.NewDecoder(netConn),
		enc: pool.NewEncoder(netConn),

		Receivers:  make(map[uint8]Receiver),
		ShutdownCh: make(chan struct{}),

		timeout: config.Timeout,
		lgr:     config.Lager,
		pool:    pool,
	}, nil
}

// Send encodes a frame on conn using t and e
func (c *conn) Send(t uint8, e interface{}) error {
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
	c.lgr.Debugf("Sending frame: %v\n", f)

	return c.enc.Encode(f)
}

// Receive registers a receiver to receive t
func (c *conn) Receive(t uint8, r Receiver) {
	c.Receivers[t] = r
	c.lgr.Debugf("Added receiver type %d\n", t)
}

// Recv listens for frames and sends them to a receiver
func (c *conn) Recv() {
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
		c.lgr.Debugf("Received frame: %v\n", frame)
		r, ok := c.Receivers[frame.Type]
		if !ok {
			c.lgr.Warnf("dropping frame %d\n", frame.Type)
			continue
		}
		err = r.Receive(frame.Data)
		if err != nil {
			// The handler returns an error.. what now?
		}
	}
}

// IsShutdown provides a way to listen for this connection to shutdown
func (c *conn) IsShutdown() chan struct{} {
	return c.ShutdownCh
}

// Pool returns the pool used by the Conn
func (c *conn) Pool() Pool {
	return c.pool
}

// Shutdown closes the gob connection
func (c *conn) Shutdown() {
	if c.isShutdown {
		return
	}
	c.lgr.Infof("Shutting down")
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
