package gob

import (
	"testing"

	"github.com/doubledutch/mux/benchmarks"
)

func BenchmarkStringReceiver(b *testing.B) {
	benchmarks.StringReceiver(b, new(Pool))
}

func BenchmarkValueReceiver(b *testing.B) {
	benchmarks.ValueReceiver(b, new(Pool))
}
