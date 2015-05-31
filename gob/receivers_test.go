package gob

import (
	"testing"

	"github.com/doubledutch/mux/tests"
)

func TestReceiver(t *testing.T) {
	tests.Receiver(t, new(Pool))
}
