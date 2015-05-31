package tests

import (
	"testing"

	"github.com/doubledutch/mux"
)

func Receiver(t *testing.T, pool mux.Pool) {
	strCh := make(chan string, 1)
	expected := "hello world"

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

// BenchmarkStringReceiver benchmarks a mux.Pool using a StringReceiver
func BenchmarkStringReceiver(b *testing.B, pool mux.Pool) {
	ch := make(chan string, 1)

	r := mux.NewStringReceiver(ch, pool)

	go func() {
		for _ = range ch {
		}
	}()

	for i := 0; i < b.N; i++ {
		r.Receive([]byte("hello"))
	}
}

// BenchmarkValueReceiver benchmarks a mux.Pool using a ValueReceiver
func BenchmarkValueReceiver(b *testing.B, pool mux.Pool) {
	ch := make(chan string, 1)

	r := mux.NewReceiver(ch, pool)

	go func() {
		for _ = range ch {
		}
	}()

	for i := 0; i < b.N; i++ {
		r.Receive([]byte("hello"))
	}
}
