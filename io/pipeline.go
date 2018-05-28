package io

import (
	"io"
	// "io/ioutil"
	// "fmt"
	// "log"
	"os/exec"
	// "bytes"
)

type Pipeline []io.ReadWriter

type Layer struct {
	cmd *exec.Cmd
}

func (l *Layer) Write(p []byte) (n int, err error) {
	return l.cmd.Stdout.Write(p)
}
func (l *Layer) Read(p []byte) (n int, err error) {
	return l.cmd.Stdin.Read(p)
}

// https://gist.github.com/tyndyll/89fbb2c2273f83a074dc

func run(layers []*exec.Cmd, pipes []*io.PipeWriter) {
	if len(layers) > 1 {
		defer func() {
			pipes[0].Close()
			run(layers[1:], pipes[1:])
		}()
	}

	layers[0].Wait()
}

func connect(input io.Reader, output io.Writer, layers ...*exec.Cmd) {
	pipes := make([]*io.PipeWriter, len(layers)-1)

	// piping input and output
	// TODO link first layer (input)
	for i := 0; i < len(layers)-1; i++ {
		// intermediate pipe
		r, w := io.Pipe()
		layers[i].Stdout = w
		layers[i+1].Stdin = r // next element exact!
		pipes[i] = w          // save pipe for next loops
	}
	// link last layer (output)
	layers[len(layers)-1].Stdout = output

	// start the pipeline
	for _, layer := range layers {
		layer.Start()
	}

	// run execution and chaining
	run(layers, pipes)
}
