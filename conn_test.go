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
	"net"
	"testing"
	"time"
)

func TestGobConnection(t *testing.T) {
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

		gConn := NewGobConn(conn)

		dec := NewDecoder()
		recvr := make(chan string, 1)

		strR := StringReceiver{
			dec: dec,
			ch:  recvr,
		}
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

	gConn := NewGobConn(conn)

	dec := NewDecoder()
	recvr := make(chan string, 1)

	strR := StringReceiver{
		dec: dec,
		ch:  recvr,
	}
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

	gConn := NewGobConn(conn)
	gConn.Shutdown()
	gConn.Shutdown()
}

func TestTimeoutSend(t *testing.T) {
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

		gConn := NewGobConn(conn)
		gConn.timeout = timeout

		recvr := make(chan string, 1)

		strR := NewStringReceiver(recvr)
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

	gConn := NewGobConn(conn)
	gConn.timeout = timeout

	recvr := make(chan string, 1)
	strR := NewStringReceiver(recvr)
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

		gConn := NewGobConn(conn)
		go gConn.Recv()
	}()

	// client
	conn, err := net.Dial("tcp", l.Addr().String())

	gConn := NewGobConn(conn)
	gConn.Send(LogType, "asdf")
	gConn.Send(LogType, "asdf")
	gConn.Shutdown()
}
