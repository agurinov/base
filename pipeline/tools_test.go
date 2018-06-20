package pipeline

import (
	"bytes"
	"reflect"
	"testing"
)

func TestToCloser(t *testing.T) {
	noCloser := bytes.NewBuffer([]byte{})
	nativeReadCloser := toReadCloser(noCloser)
	nativeWriteCloser := toWriteCloser(noCloser)

	t.Run("Read", func(t *testing.T) {
		t.Run("native", func(t *testing.T) {
			oldPtr := reflect.ValueOf(nativeReadCloser).Pointer()
			newPtr := reflect.ValueOf(toReadCloser(nativeReadCloser)).Pointer()

			// native ReadCloser will return without any injections -> same pointer
			if oldPtr != newPtr {
				t.Error("Unexpected pointer")
			}
		})

		t.Run("obtained", func(t *testing.T) {
			oldPtr := reflect.ValueOf(noCloser).Pointer()
			newPtr := reflect.ValueOf(toReadCloser(noCloser)).Pointer()

			// no ReadCloser will be returned with injection of .close() method (just return nil)
			// -> different pointer
			if oldPtr == newPtr {
				t.Error("Unexpected pointer")
			}
		})
	})

	t.Run("Write", func(t *testing.T) {
		t.Run("native", func(t *testing.T) {
			oldPtr := reflect.ValueOf(nativeWriteCloser).Pointer()
			newPtr := reflect.ValueOf(toWriteCloser(nativeWriteCloser)).Pointer()

			// native WriteCloser will return without any injections -> same pointer
			if oldPtr != newPtr {
				t.Error("Unexpected pointer")
			}
		})

		t.Run("obtained", func(t *testing.T) {
			oldPtr := reflect.ValueOf(noCloser).Pointer()
			newPtr := reflect.ValueOf(toWriteCloser(noCloser)).Pointer()

			// no WriteCloser will be returned with injection of .close() method (just return nil)
			// -> different pointer
			if oldPtr == newPtr {
				t.Error("Unexpected pointer")
			}
		})
	})
}

func TestPiping(t *testing.T) {
	input := toReadCloser(bytes.NewBuffer([]byte{}))
	output := toWriteCloser(bytes.NewBuffer([]byte{}))
	inputPtr := reflect.ValueOf(input).Pointer()
	outputPtr := reflect.ValueOf(output).Pointer()

	t.Run("tcp", func(t *testing.T) {
		if inputPtr == outputPtr {
			t.Fatal("Unexpected same pointers for input and output")
		}

		t.Run("len==1", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewTCPSocket("example.com:80"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdinPtr := reflect.ValueOf(layers[0].(*tcp).stdin).Pointer()
			stdoutPtr := reflect.ValueOf(layers[0].(*tcp).stdout).Pointer()

			if stdinPtr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			if stdoutPtr != outputPtr {
				t.Fatal("layers[0]: unexpected stdout")
			}
		})

		t.Run("len==2", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewTCPSocket("example.com:80"),
				NewTCPSocket("domain.com:22"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdin1Ptr := reflect.ValueOf(layers[0].(*tcp).stdin).Pointer()
			// stdout1Ptr := reflect.ValueOf(layers[0].(*tcp).stdout).Pointer()
			// stdin2Ptr := reflect.ValueOf(layers[1].(*tcp).stdin).Pointer()
			stdout2Ptr := reflect.ValueOf(layers[1].(*tcp).stdout).Pointer()

			if stdin1Ptr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			// if stdout1Ptr != stdin2Ptr {
			// 	t.Fatal("layers[0]: unexpected stdout")
			// }
			//
			// if stdin2Ptr != stdout1Ptr {
			// 	t.Fatal("layers[1]: unexpected stdin")
			// }
			if stdout2Ptr != outputPtr {
				t.Fatal("layers[1]: unexpected stdout")
			}
		})
	})

	t.Run("process", func(t *testing.T) {
		if inputPtr == outputPtr {
			t.Fatal("Unexpected same pointers for input and output")
		}

		t.Run("len==1", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewProcess("pwd"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdinPtr := reflect.ValueOf(layers[0].(*process).stdin).Pointer()
			stdoutPtr := reflect.ValueOf(layers[0].(*process).stdout).Pointer()

			if stdinPtr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			if stdoutPtr != outputPtr {
				t.Fatal("layers[0]: unexpected stdout")
			}
		})

		t.Run("len==2", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewProcess("pwd"),
				NewProcess("rev"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdin1Ptr := reflect.ValueOf(layers[0].(*process).stdin).Pointer()
			// stdout1Ptr := reflect.ValueOf(layers[0].(*process).stdout).Pointer()
			// stdin2Ptr := reflect.ValueOf(layers[1].(*process).stdin).Pointer()
			stdout2Ptr := reflect.ValueOf(layers[1].(*process).stdout).Pointer()

			if stdin1Ptr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			// if stdout1Ptr != stdin2Ptr {
			// 	t.Fatal("Unexpected stdout for [0] layer")
			// }
			//
			// if stdin2Ptr != stdout1Ptr {
			// 	t.Fatal("Unexpected stdin for [1] layer")
			// }
			if stdout2Ptr != outputPtr {
				t.Fatal("layers[1]: unexpected stdout")
			}
		})
	})

	t.Run("mix", func(t *testing.T) {
		t.Run("len==3", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewProcess("pwd"),
				NewTCPSocket("example.com:80"),
				NewProcess("rev"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdin1Ptr := reflect.ValueOf(layers[0].(*process).stdin).Pointer()
			// stdout1Ptr := reflect.ValueOf(layers[0].(*process).stdout).Pointer()
			//
			// stdin2Ptr := reflect.ValueOf(layers[1].(*tcp).stdin).Pointer()
			// stdout2Ptr := reflect.ValueOf(layers[1].(*tcp).stdout).Pointer()
			//
			// stdin3Ptr := reflect.ValueOf(layers[2].(*process).stdin).Pointer()
			stdout3Ptr := reflect.ValueOf(layers[2].(*process).stdout).Pointer()

			if stdin1Ptr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			// if stdout1Ptr != stdin2Ptr {
			// 	t.Fatal("Unexpected stdout for [0] layer")
			// }
			//
			// if stdin2Ptr != stdout1Ptr {
			// 	t.Fatal("Unexpected stdin for [1] layer")
			// }
			if stdout3Ptr != outputPtr {
				t.Fatal("layers[2]: unexpected stdout")
			}
		})
	})
}

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
