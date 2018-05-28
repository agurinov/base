package io

import (
	"testing"
	"bytes"
	"os/exec"
	"fmt"
)

func BenchmarkChaining(b *testing.B) {
	// tableTests := []struct {
	// 	in  string // url for parsing
	// 	out net.IP // expected value of IP
	// }{
	// 	{"", nil},
	// 	{"tcp://golang.org", nil},
	// 	{"http://golang.org", nil},
	// 	{"tcp://192.168.99.101:2376", net.ParseIP("192.168.99.101")},
	// }
	b.Run("first", func(b *testing.B) {
		for i := 0; i < 1; i++ {
			input := bytes.NewBuffer([]byte{'g', 'o', 'l', 'a', 'n'})
			output := bytes.NewBuffer([]byte{})

			chaining(
				input,
				output,
				exec.Command("echo", "$@"),
			)

			fmt.Println("RESPONSE", output.String())
		}
	})
}
