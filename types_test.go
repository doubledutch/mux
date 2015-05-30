package mux

import (
	"errors"
	"net"
	"testing"
)

func TestHappyClientServer(t *testing.T) {
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
		logR := NewStringReceiver(logCh)

		server, err := NewDefaultServer(conn)
		if err != nil {
			t.Fatal(err)
		}
		server.Receive(LogType, logR)
		go server.Recv()
		actual := <-logCh
		if actual != text {
			t.Fatalf("'%s' != '%s'", actual, text)
		}

		server.Done(nil)
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	client, err := NewDefaultClient(conn)
	if err != nil {
		t.Fatal(err)
	}

	go client.Recv()

	if err := client.Send(LogType, text); err != nil {
		t.Fatal(err)
	}

	if err := client.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestClientServerErr(t *testing.T) {
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

		server, err := NewDefaultServer(conn)
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
	client, err := NewDefaultClient(conn)
	if err != nil {
		t.Fatal(err)
	}

	go client.Recv()

	if err := client.Wait(); err != nil && err.Error() != expectedErr.Error() {
		t.Fatal(err)
	}
}
