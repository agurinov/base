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

func run(layers []*exec.Cmd, pipeWriters []*io.PipeWriter) (err error) {
	for i, layer := range layers {
		if err := layer.Wait(); err != nil {
			return err
		}

		if i < len(layers)-1 {
			if err := pipeWriters[i].Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func connect(input io.Reader, output io.Writer, layers ...*exec.Cmd) (err error) {
	//pipes count less than layers by one, because last layer no need in piping
	pipeWriters := make([]*io.PipeWriter, len(layers)-1)

	// piping input and output
	// TODO link first layer (input)
	for i := 0; i < len(layers)-1; i++ {
		// intermediate pipe
		r, w := io.Pipe()
		layers[i].Stdout = w
		layers[i+1].Stdin = r // next element exact!
		pipeWriters[i] = w    // save pipe for next loops
	}
	// link last layer (output)
	layers[len(layers)-1].Stdout = output

	// start the pipeline
	for _, layer := range layers {
		if err := layer.Start(); err != nil {
			return err
		}
	}

	// run execution and chaining
	return run(layers, pipeWriters)
}
