package benchmarks

import (
	"testing"

	"github.com/doubledutch/mux"
)

// StringReceiver benchmarks a mux.Pool using a StringReceiver
func StringReceiver(b *testing.B, pool mux.Pool) {
	ch := make(chan string, 1)
	done := make(chan struct{})

	r := mux.NewStringReceiver(ch, pool)

	go func() {
		for _ = range ch {
		}
		close(done)
	}()

	var err error
	enc := pool.NewBufferEncoder()
	str := "hello world"
	for i := 0; i < b.N; i++ {
		enc.Encode(&str)
		if err = r.Receive(enc.Bytes()); err != nil {
			b.Fatal(err)
		}
		enc.Reset()
	}
	r.Close()

	<-done
}

// ValueReceiver benchmarks a mux.Pool using a ValueReceiver
func ValueReceiver(b *testing.B, pool mux.Pool) {
	ch := make(chan string, 1)
	done := make(chan struct{})

	r := mux.NewReceiver(ch, pool)

	go func() {
		for _ = range ch {
		}
		close(done)
	}()

	var err error
	enc := pool.NewBufferEncoder()
	str := "hello world"
	for i := 0; i < b.N; i++ {
		enc.Encode(&str)
		if err = r.Receive(enc.Bytes()); err != nil {
			b.Fatal(err)
		}
		enc.Reset()
	}
	r.Close()

	<-done
}
