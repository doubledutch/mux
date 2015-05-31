package json

import (
	"testing"

	"github.com/doubledutch/mux/tests"
)

func TestConnection(t *testing.T) {
	tests.Connection(t, new(Pool), NewDefaultConn)
}

func TestShutdown(t *testing.T) {
	tests.Shutdown(t, NewDefaultConn)
}

func TestTimeoutSend(t *testing.T) {
	tests.TimeoutSend(t, new(Pool), NewConn)
}

func TestDroppedMessages(t *testing.T) {
	tests.DroppedMessages(t, new(Pool), NewDefaultConn)
}
