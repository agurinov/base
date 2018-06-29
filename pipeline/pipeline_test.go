package pipeline

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPipelineRun(t *testing.T) {
	t.Run("processes", func(t *testing.T) {
		input := bytes.NewBuffer([]byte("foobar"))
		output := bytes.NewBuffer([]byte{})

		process1 := NewProcess("cat /dev/stdin") // read 'foobar' from stdin
		process2 := NewProcess("rev")            // reverse -> raboof
		process3 := NewProcess("grep -o raboof") // grep reversed (must be 1 match)
		process4 := NewProcess("wc -l")          // count matches

		pipeline := New(process1, process2, process3, process4)

		if err := pipeline.Run(input, output); err != nil {
			t.Fatal(err)
		}

		if outputString := fmt.Sprint(output); outputString != "1\n" {
			t.Fatalf("Expected %q, got %q", "{1\n}", outputString)
		}
	})

	t.Run("mix", func(t *testing.T) {
		input := bytes.NewBuffer([]byte("HEAD / HTTP/1.0\r\n\r\n"))
		output := bytes.NewBuffer([]byte{})

		process := NewProcess("cat /dev/stdin") // read simple http request from process stdin
		socket := NewTCPSocket("golang.org:80") // and pass to golang.org via socket

		pipeline := New(process, socket)

		if err := pipeline.Run(input, output); err != nil {
			t.Fatal(err)
		}

		// t.Log(output)
		// if outputString := fmt.Sprint(output); outputString != "{1\n}" {
		// 	t.Fatalf("Expected %q, got %q", "{1\n}", outputString)
		// }
	})
}
