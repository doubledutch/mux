package tests

import (
	"testing"

	"github.com/doubledutch/mux"
)

// Receiver tests Receiver with a mux.Pool
func Receiver(t *testing.T, pool mux.Pool) {
	strCh := make(chan string, 1)
	expected := "hello world"

	strR := NewReceiver(strCh)
	defer strR.Close()

	enc := pool.NewBufferEncoder()
	enc.Encode(expected)

	if err := strR.Receive(enc.Bytes()); err != nil {
		t.Fatal(err)
	}

	actual := <-strCh
	if actual != expected {
		t.Fatalf("actual '%s' != expected '%s'", actual, expected)
	}
}
