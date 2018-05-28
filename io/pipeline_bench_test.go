package io

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
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
		for i := 0; i < b.N; i++ {
			b.StopTimer()

			input := bytes.NewBuffer([]byte{'g', 'o', 'l', 'a', 'n'})
			output := bytes.NewBuffer([]byte{})

			b.StartTimer()

			connect(
				input,
				output,
				exec.Command("ls"),
				// exec.Command("grep bench"),
				exec.Command("wc", "-l"),
			)

			b.StopTimer()

			fmt.Println("RESPONSE", output.String())
		}
	})
}
