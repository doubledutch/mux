package msgpack

import (
	"testing"

	"github.com/doubledutch/mux/tests"
)

func TestStringReceiver(t *testing.T) {
	tests.StringReceiver(t, new(Pool))
}

func TestSignalReceiver(t *testing.T) {
	tests.SignalReceiver(t, new(Pool))
}
