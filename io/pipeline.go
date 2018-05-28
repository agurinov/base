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

func run(layers []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	// look at for loop
	if len(layers) > 1 {
		defer func() {
			if err = pipes[0].Close(); err != nil {
				return err
			}

			if err = run(layers[1:], pipes[1:]); err != nil {
				return err
			}
		}()
	}

	return layers[0].Wait()
}

func connect(input io.Reader, output io.Writer, layers ...*exec.Cmd) (err error) {
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
		if err = layer.Start(); err != nil {
			return err
		}
	}

	// run execution and chaining
	return run(layers, pipes)
}
