package tests

import (
	"errors"
	"net"
	"testing"

	"github.com/doubledutch/mux"
)

// NewServer defines a func for creating new servers
type NewServer func(conn net.Conn) (mux.Server, error)

// NewClient defines a func for creating new clients
type NewClient func(conn net.Conn) (mux.Client, error)

// HappyClientServer tests client server communication
func HappyClientServer(t *testing.T, pool mux.Pool, newServer NewServer, newClient NewClient) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	text := "hello world"

	// server
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		logCh := make(chan string, 1)
		logR := pool.NewReceiver(logCh)

		server, err := newServer(conn)
		if err != nil {
			t.Fatal(err)
		}
		server.Receive(mux.LogType, logR)
		go server.Recv()
		actual := <-logCh
		if actual != text {
			t.Fatalf("'%s' != '%s'", actual, text)
		}

		server.Done(nil)
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	client, err := newClient(conn)
	if err != nil {
		t.Fatal(err)
	}

	go client.Recv()

	if err := client.Send(mux.LogType, text); err != nil {
		t.Fatal(err)
	}

	if err := client.Wait(); err != nil {
		t.Fatal(err)
	}
}

// ClientServerErr tests client server communication during error
func ClientServerErr(t *testing.T, pool mux.Pool, newServer NewServer, newClient NewClient) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	expectedErr := errors.New("error")

	// server
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		server, err := newServer(conn)
		if err != nil {
			t.Fatal(err)
		}
		go server.Recv()

		server.Done(expectedErr)
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	client, err := newClient(conn)
	if err != nil {
		t.Fatal(err)
	}

	go client.Recv()

	if err := client.Wait(); err != nil && err.Error() != expectedErr.Error() {
		t.Fatal(err)
	}
}
