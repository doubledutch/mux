package json

import (
	"testing"

	"github.com/doubledutch/mux/tests"
)

func TestHappyClientServer(t *testing.T) {
	tests.HappyClientServer(t, new(Pool), NewDefaultServer, NewDefaultClient)
}

func TestClientServerErr(t *testing.T) {
	tests.ClientServerErr(t, new(Pool), NewDefaultServer, NewDefaultClient)
}
