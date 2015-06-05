package tests

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/doubledutch/lager"
	"github.com/doubledutch/mux"
)

// Lager returns a new test Lager
func Lager() lager.Lager {
	return lager.NewLogLager(&lager.LogConfig{
		Levels: lager.LevelsFromString(os.Getenv("LOG_LEVELS")),
		Output: os.Stderr,
	})
}

// NewConn defines a func for creating a mux.Conn with default configuration
type NewConn func(conn net.Conn) (mux.Conn, error)

// NewConfigConn defines a func for creating a new Conn with configuration
type NewConfigConn func(conn net.Conn, config *mux.Config) (mux.Conn, error)

// Connection is a basic Conn test
func Connection(t *testing.T, pool mux.Pool, newConn NewConn) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	logType := uint8(1)
	errType := uint8(2)
	text := "hello world"

	// server
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		mConn, err := newConn(conn)
		if err != nil {
			t.Fatal(err)
		}

		recvr := make(chan string, 1)
		strR := pool.NewReceiver(recvr)
		mConn.Receive(logType, strR)
		go mConn.Recv()
		actual := <-recvr
		if actual != text {
			t.Fatalf("'%s' != '%s'", actual, text)
		}

		mConn.Send(errType, "")
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	mConn, err := newConn(conn)
	if err != nil {
		t.Fatal(err)
	}

	recvr := make(chan string, 1)
	strR := pool.NewReceiver(recvr)
	mConn.Receive(errType, strR)
	go mConn.Recv()

	if err := mConn.Send(logType, text); err != nil {
		t.Fatal(err)
	}

	<-recvr
	mConn.Shutdown()
	conn.Close()
}

// Shutdown tests shutting down a mux.Conn
func Shutdown(t *testing.T, newConn NewConn) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		conn.Close()
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	mConn, err := newConn(conn)
	if err != nil {
		t.Fatal(err)
	}
	mConn.Shutdown()
}

// TimeoutSend tests
func TimeoutSend(t *testing.T, pool mux.Pool, newConn NewConfigConn) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	timeout := 10 * time.Millisecond

	config := &mux.Config{
		Timeout: timeout,
		Lager:   Lager(),
	}

	logType := uint8(1)
	errType := uint8(2)
	text := "hello world"

	// server
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		mConn, err := newConn(conn, config)
		if err != nil {
			t.Fatal(err)
		}

		recvr := make(chan string, 1)

		strR := pool.NewReceiver(recvr)
		mConn.Receive(logType, strR)
		go mConn.Recv()

		actual := <-recvr
		if actual != text {
			t.Fatalf("'%s' != '%s'", actual, text)
		}

		mConn.Send(errType, "")
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	mConn, err := newConn(conn, config)
	if err != nil {
		t.Fatal(err)
	}

	recvr := make(chan string, 1)
	strR := pool.NewReceiver(recvr)
	mConn.Receive(errType, strR)
	go mConn.Recv()

	time.Sleep(10 * timeout)
	if err := mConn.Send(logType, text); err != nil {
		t.Fatal(err)
	}

	<-recvr
	mConn.Shutdown()
}

// DroppedMessages tests a mux.Conn for handling dropped messages
func DroppedMessages(t *testing.T, pool mux.Pool, newConn NewConn) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	// server
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		mConn, err := newConn(conn)
		if err != nil {
			t.Fatal(err)
		}
		go mConn.Recv()
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	mConn, err := newConn(conn)
	if err != nil {
		t.Fatal(err)
	}
	mConn.Send(mux.LogType, "asdf")
	mConn.Send(mux.LogType, "asdf")
	mConn.Shutdown()
}
