package gob

import (
	"testing"

	"github.com/doubledutch/mux"
)

func BenchmarkStringReceiver(b *testing.B) {
	ch := make(chan string, 1)

	r := mux.NewStringReceiver(ch, new(Pool))

	go func() {
		for _ = range ch {
		}
	}()

	for i := 0; i < b.N; i++ {
		r.Receive([]byte("hello"))
	}
}

func BenchmarkValueReceiver(b *testing.B) {
	ch := make(chan string, 1)

	r := mux.NewReceiver(ch, new(Pool))

	go func() {
		for _ = range ch {
		}
	}()

	for i := 0; i < b.N; i++ {
		r.Receive([]byte("hello"))
	}
}
