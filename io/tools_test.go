package io

import (
	"io"
	// "bytes"
	"testing"
)

var input io.ReadCloser
var output io.WriteCloser

func TestPiping(t *testing.T) {
	t.Run("len(layers)==1", func(t *testing.T) {
		// layers for piping
		layers := []Pipeable{
			NewSocket("example.com"),
		}
		// check for errors
		if err := piping(input, output, layers...); err != nil {
			t.Error(err)
		}
		// check layers stdio
		if layers[0].(*Socket).stdin != input {
			t.Errorf("Unexpected stdin")
		}
		if layers[0].(*Socket).stdout != output {
			t.Errorf("Unexpected stdout")
		}
	})

	t.Run("len(layers)==2", func(t *testing.T) {
		// layers for piping
		layers := []Pipeable{
			NewSocket("example.com"),
			NewProcess("echo", "foobar"),
		}
		// check for errors
		if err := piping(input, output, layers...); err != nil {
			t.Error(err)
		}
		// check layers stdio
		if layers[0].(*Socket).stdin != input {
			t.Errorf("Unexpected stdin")
		}
		// if layers[0].(*Socket).stdout != layers[1].(*Process).stdin {
		// 	t.Errorf("Unexpected pipe")
		// }
		if layers[1].(*Process).stdout != output {
			t.Errorf("Unexpected stdout")
		}
	})
}
