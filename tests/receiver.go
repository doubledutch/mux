package tests

import (
	"os"
	"testing"

	"github.com/doubledutch/mux"
)

// StringReceiver tests Receiver with a mux.Pool
func StringReceiver(t *testing.T, pool mux.Pool) {
	strCh := make(chan string, 1)
	expected := "hello world"

	strR := pool.NewReceiver(strCh)
	defer strR.Close()

	enc := pool.NewBufferEncoder()
	if err := enc.Encode(&expected); err != nil {
		t.Fatal(err)
	}

	if err := strR.Receive(enc.Bytes()); err != nil {
		t.Fatal(err)
	}

	actual := <-strCh
	if actual != expected {
		t.Fatalf("actual '%s' != expected '%s'", actual, expected)
	}
}

// SignalReceiver tests NewSignalReceiver with a mux.Pool
func SignalReceiver(t *testing.T, pool mux.Pool) {
	sigCh := make(chan os.Signal, 1)
	expected := os.Kill

	sigR := mux.NewSignalReceiver(sigCh, pool)
	defer sigR.Close()

	enc := pool.NewBufferEncoder()
	if err := enc.Encode(expected); err != nil {
		t.Fatal(err)
	}

	if err := sigR.Receive(enc.Bytes()); err != nil {
		t.Fatal(err)
	}

	actual := <-sigCh
	if actual != expected {
		t.Fatalf("actual '%v' != expected '%v'", actual, expected)
	}
}
