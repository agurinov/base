package pipeline

import (
// "fmt"
// "net"
// "testing"
)

const (
	socketFailYAML    = "foo: bar"
	socketSuccessYAML = "address: tcp://geoiphost/foo/bar"
)

// func TestSocketFromYAML(t *testing.T) {
// 	t.Run("fail", func(t *testing.T) {
// 		s, err := SocketFromYAML([]byte(socketFailYAML))
//
// 		t.Log(s)
// 		if err == nil {
// 			t.Error("Expected error, got \"nil\"")
// 		}
// 	})
//
// 	t.Run("success", func(t *testing.T) {
// 		s, err := SocketFromYAML([]byte(socketSuccessYAML))
// 		if err != nil {
// 			t.Error(err)
// 		}
//
// 		if s == nil {
// 			t.Errorf("Expected socket, got error: %q", err)
// 		}
// 		// t.Log(s)
// 	})
// }

// func TestSCRRR(t *testing.T) {
// 	tableTests := []struct {
// 		address string // source
// 		out     string // expected value of String() method
// 	}{
// 		{"tcp://foobar", ""},
// 		{"http://foobar", ""},
// 		{"https://foobar", ""},
// 		{"http://golang.org", ""},
// 		{"https://golang.org", ""},
// 		{"tcp://golang.org", ""},
// 		{"golang.org", ""},
// 	}
//
// 	for i, tt := range tableTests {
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			names, err := net.LookupAddr(tt.address)
// 			t.Log(names, err)
// 		})
// 	}
// }
