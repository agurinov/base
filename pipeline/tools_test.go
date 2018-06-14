package pipeline

import (
	// "bytes"
	"io"
	// "testing"
)

type readCloser struct {
	io.Reader
}

func (readCloser) Close() error { return nil }

type writeCloser struct {
	io.Writer
}

func (writeCloser) Close() error { return nil }

// func TestPiping(t *testing.T) {
// 	t.Run("len(layers)==1", func(t *testing.T) {
// 		// layers for piping
// 		layers := []Pipeable{
// 			NewSocket("example.com"),
// 		}
// 		// check for errors
// if err := piping(input, output, layers...); err != nil {
// 	t.Error(err)
// }
// 		// check layers stdio
// 		if layers[0].(*Socket).stdin != input {
// 			t.Errorf("Unexpected stdin")
// 		}
// 		if layers[0].(*Socket).stdout != output {
// 			t.Errorf("Unexpected stdout")
// 		}
// 	})
//
// 	t.Run("len(layers)==2", func(t *testing.T) {
// 		// layers for piping
// 		layers := []Pipeable{
// 			NewSocket("example.com"),
// 			NewProcess("echo", "foobar"),
// 		}
// 		// check for errors
// 		if err := piping(input, output, layers...); err != nil {
// 			t.Error(err)
// 		}
// 		// check layers stdio
// 		if layers[0].(*Socket).stdin != input {
// 			t.Errorf("Unexpected stdin")
// 		}
// 		// if layers[0].(*Socket).stdout != layers[1].(*Process).stdin {
// 		// 	t.Errorf("Unexpected pipe")
// 		// }
// 		if layers[1].(*Process).stdout != output {
// 			t.Errorf("Unexpected stdout")
// 		}
// 	})
// }

// func TestRun(t *testing.T) {
// 	t.Run("processes", func(t *testing.T) {
// 		input := readCloser{bytes.NewBuffer([]byte("foobar"))}
// 		output := writeCloser{bytes.NewBuffer([]byte{})}
//
// 		process1 := NewProcess("cat", "/dev/stdin")    // read 'foobar' from stdin
// 		process2 := NewProcess("rev")                  // reverse -> raboof
// 		process3 := NewProcess("grep", "-o", "raboof") // grep reversed (must be 1 match)
// 		process4 := NewProcess("wc", "-l")             // count matches
//
// 		layers1 := []Able{process1, process2, process3, process4}
// 		layers2 := []Exec{process1, process2, process3, process4}
//
// 		if err := piping(input, output, layers1...); err != nil {
// 			t.Error(err)
// 		}
// 		if err := run(layers2...); err != nil {
// 			t.Error(err)
// 		}
// 		// TODO
// 		t.Log(output)
// 		// if output.(*Buffer) != "1" {
// 		// 	t.Errorf("Expected %q, got %q", "1", string(output))
// 		// }
// 	})
//
// 	t.Run("sockets", func(t *testing.T) {
// 		input := readCloser{bytes.NewBuffer([]byte("HEAD / HTTP/1.0\r\n\r\n"))}
// 		output := writeCloser{bytes.NewBuffer([]byte{})}
//
// 		// process := NewProcess("cat", "/dev/stdin") // read simple http request from stdin
// 		// process := NewProcess("echo", "HEAD / HTTP/1.0\r\n\r\n") // read simple http request from stdin
// 		socket := NewTCSocket("golang.org:80") // and pass to golang.org via socket
//
// 		layers1 := []Able{socket}
// 		layers2 := []Exec{socket}
//
// 		if err := piping(input, output, layers1...); err != nil {
// 			t.Error(err)
// 		}
// 		if err := run(layers2...); err != nil {
// 			t.Error(err)
// 		}
// 		// TODO
// 		t.Log(output)
// 		// if output.(*Buffer) != "1" {
// 		// 	t.Errorf("Expected %q, got %q", "1", string(output))
// 		// }
// 	})
// }
