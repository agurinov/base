package pipeline

import (
	"testing"
)

const (
	socketFailYAML    = "foo: bar"
	socketSuccessYAML = "address: tcp://geoiphost/foo/bar"
)

func TestSocketFromYAML(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		s, err := SocketFromYAML([]byte(socketFailYAML))

		t.Log(s)
		if err == nil {
			t.Error("Expected error, got \"nil\"")
		}
	})

	t.Run("success", func(t *testing.T) {
		s, err := SocketFromYAML([]byte(socketSuccessYAML))
		if err != nil {
			t.Error(err)
		}

		if s == nil {
			t.Errorf("Expected socket, got error: %q", err)
		}
		// t.Log(s)
	})
}
