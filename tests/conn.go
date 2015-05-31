package tests

import (
	"net"
	"testing"
	"time"

	"github.com/doubledutch/lager"
	"github.com/doubledutch/mux"
)

func Lager() lager.Lager {
	return lager.NewLogLager(nil)
}

type NewConn func(conn net.Conn) (mux.Conn, error)

type NewConfigConn func(conn net.Conn, config *mux.Config) (mux.Conn, error)

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
		strR := mux.NewReceiver(recvr, pool)
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
	strR := mux.NewReceiver(recvr, pool)
	mConn.Receive(errType, strR)
	go mConn.Recv()

	if err := mConn.Send(logType, text); err != nil {
		t.Fatal(err)
	}

	<-recvr
	mConn.Shutdown()
	conn.Close()
}

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

func TimeoutSend(t *testing.T, pool mux.Pool, newConn NewConfigConn) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	timeout := 1 * time.Millisecond

	logType := uint8(1)
	errType := uint8(2)
	text := "hello world"

	// server
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		mConn, err := newConn(conn, &mux.Config{
			Timeout: timeout,
			Lager:   Lager(),
		})
		if err != nil {
			t.Fatal(err)
		}

		recvr := make(chan string, 1)

		strR := mux.NewReceiver(recvr, pool)
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

	mConn, err := newConn(conn, &mux.Config{
		Timeout: timeout,
		Lager:   Lager(),
	})
	if err != nil {
		t.Fatal(err)
	}

	recvr := make(chan string, 1)
	strR := mux.NewReceiver(recvr, pool)
	mConn.Receive(errType, strR)
	go mConn.Recv()

	time.Sleep(10 * timeout)
	if err := mConn.Send(logType, text); err != nil {
		t.Fatal(err)
	}

	<-recvr
	mConn.Shutdown()
}

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
