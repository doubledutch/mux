package mux

import (
	"os"
	"testing"
)

func TestStringReceiver(t *testing.T) {
	strCh := make(chan string, 1)
	expected := "hello world"

	strR := NewStringReceiver(strCh)
	defer strR.Close()

	enc := NewEncoder()
	enc.Encode(expected)

	if err := strR.Receive(enc.Bytes()); err != nil {
		t.Fatal(err)
	}

	actual := <-strCh
	if actual != expected {
		t.Fatalf("actual '%s' != expected '%s'", actual, expected)
	}

}

func TestSignalReceiver(t *testing.T) {
	sigCh := make(chan os.Signal, 1)

	sigR := NewSignalReceiver(sigCh)
	defer sigR.Close()

	enc := NewEncoder()
	enc.Encode(os.Kill)

	if err := sigR.Receive(enc.Bytes()); err != nil {
		t.Fatal(err)
	}

	if <-sigCh != os.Kill {
		t.Fatal("expected os.Kill")
	}
}
