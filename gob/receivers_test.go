package gob

import (
	"testing"

	"github.com/doubledutch/mux"
)

func TestReceiver(t *testing.T) {
	strCh := make(chan string, 1)
	expected := "hello world"

	pool := new(Pool)

	strR := mux.NewReceiver(strCh, pool)
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
