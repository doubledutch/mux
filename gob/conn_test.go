package gob

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

func TestConnection(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	pool := new(Pool)

	logType := uint8(1)
	errType := uint8(2)
	text := "hello world"

	// server
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		gConn, err := NewDefaultConn(conn)
		if err != nil {
			t.Fatal(err)
		}

		recvr := make(chan string, 1)
		strR := mux.NewReceiver(recvr, pool)
		gConn.Receive(logType, strR)
		go gConn.Recv()
		actual := <-recvr
		if actual != text {
			t.Fatalf("'%s' != '%s'", actual, text)
		}

		gConn.Send(errType, "")
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	gConn, err := NewDefaultConn(conn)
	if err != nil {
		t.Fatal(err)
	}

	recvr := make(chan string, 1)
	strR := mux.NewReceiver(recvr, pool)
	gConn.Receive(errType, strR)
	go gConn.Recv()

	if err := gConn.Send(logType, text); err != nil {
		t.Fatal(err)
	}

	<-recvr
	gConn.Shutdown()
	conn.Close()
}

func TestShutdown(t *testing.T) {
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

	gConn, err := NewDefaultConn(conn)
	if err != nil {
		t.Fatal(err)
	}
	gConn.Shutdown()
}

func TestTimeoutSend(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	pool := new(Pool)

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

		gConn, err := NewConn(conn, &mux.Config{
			Timeout: timeout,
			Lager:   Lager(),
		})
		if err != nil {
			t.Fatal(err)
		}

		recvr := make(chan string, 1)

		strR := mux.NewReceiver(recvr, pool)
		gConn.Receive(logType, strR)
		go gConn.Recv()
		actual := <-recvr
		if actual != text {
			t.Fatalf("'%s' != '%s'", actual, text)
		}

		gConn.Send(errType, "")
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	gConn, err := NewConn(conn, &mux.Config{
		Timeout: timeout,
		Lager:   Lager(),
	})
	if err != nil {
		t.Fatal(err)
	}

	recvr := make(chan string, 1)
	strR := mux.NewReceiver(recvr, pool)
	gConn.Receive(errType, strR)
	go gConn.Recv()

	time.Sleep(10 * timeout)
	if err := gConn.Send(logType, text); err != nil {
		t.Fatal(err)
	}

	<-recvr
	gConn.Shutdown()
}

func TestDroppedMessages(t *testing.T) {
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
		gConn, err := NewDefaultConn(conn)
		if err != nil {
			t.Fatal(err)
		}
		go gConn.Recv()
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	gConn, err := NewDefaultConn(conn)
	if err != nil {
		t.Fatal(err)
	}
	gConn.Send(mux.LogType, "asdf")
	gConn.Send(mux.LogType, "asdf")
	gConn.Shutdown()
}
