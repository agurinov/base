package io

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"testing"
)

func BenchmarkChaining(b *testing.B) {
	b.Run("first", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()

			input := bytes.NewBuffer([]byte{'g', 'o', 'l', 'a', 'n'})
			output := bytes.NewBuffer([]byte{})

			b.StartTimer()

			err := connect(
				input,
				output,
				exec.Command("ls"),
				// exec.Command("grep bench"),
				exec.Command("wc", "-l"),
			)

			b.StopTimer()

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("RESPONSE", output.String())
		}
	})
}
