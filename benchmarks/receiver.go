package benchmarks

import (
	"testing"

	"github.com/doubledutch/mux"
)

// StringReceiver benchmarks a mux.Pool using a StringReceiver
func StringReceiver(b *testing.B, pool mux.Pool) {
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

// ValueReceiver benchmarks a mux.Pool using a ValueReceiver
func ValueReceiver(b *testing.B, pool mux.Pool) {
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
