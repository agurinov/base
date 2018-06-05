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
}
