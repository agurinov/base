package pipeline

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPipelineRun(t *testing.T) {
	t.Run("processes", func(t *testing.T) {
		input := readCloser{bytes.NewBuffer([]byte("foobar"))}
		output := writeCloser{bytes.NewBuffer([]byte{})}

		process1 := NewProcess("cat", "/dev/stdin")    // read 'foobar' from stdin
		process2 := NewProcess("rev")                  // reverse -> raboof
		process3 := NewProcess("grep", "-o", "raboof") // grep reversed (must be 1 match)
		process4 := NewProcess("wc", "-l")             // count matches

		layers := []Layer{process1, process2, process3, process4}

		pipeline := Pipeline{layers}

		if err := pipeline.Run(input, output); err != nil {
			t.Error(err)
		}

		if outputString := fmt.Sprint(output); outputString != "{1\n}" {
			t.Errorf("Expected %q, got %q", "{1\n}", outputString)
		}
	})

	t.Run("mix", func(t *testing.T) {
		input := readCloser{bytes.NewBuffer([]byte("HEAD / HTTP/1.0\r\n\r\n"))}
		output := writeCloser{bytes.NewBuffer([]byte{})}

		process := NewProcess("cat", "/dev/stdin") // read simple http request from process stdin
		socket := NewSocket("golang.org:80")       // and pass to golang.org via socket

		layers := []Layer{process, socket}

		pipeline := Pipeline{layers}

		if err := pipeline.Run(input, output); err != nil {
			t.Error(err)
		}

		t.Log(output)
		// if outputString := fmt.Sprint(output); outputString != "{1\n}" {
		// 	t.Errorf("Expected %q, got %q", "{1\n}", outputString)
		// }
	})
}
