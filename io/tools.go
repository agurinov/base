package io

import (
	"io"
	"os/exec"
)

// piping establishes pipe connections between IO processes (Layer)
// the first layer accepts as stdin the input buffer
// the last layer puts into output buffer his stdout
func piping(input io.Reader, output io.Writer, layers ...*exec.Cmd) (err error) {
	//pipes count less than layers by one, because last layer no need in piping
	pipeWriters := make([]*io.PipeWriter, len(layers)-1)

	// piping input and output
	// TODO link first layer (input)

	// piping intermediate layers
	for i := 0; i < len(layers)-1; i++ {
		// intermediate pipe
		r, w := io.Pipe()
		layers[i].Stdout = w
		layers[i+1].Stdin = r // next element exact!
		pipeWriters[i] = w    // save pipe for next loops
	}

	// link last layer (output)
	layers[len(layers)-1].Stdout = output

	return nil
}

// start invokes the layer's .Start() method in the order of the queue
func start(layers []*exec.Cmd) (err error) {
	// start the pipeline
	for _, layer := range layers {
		if err := layer.Start(); err != nil {
			return err
		}
	}

	return nil
}

// run causes processes in turn
// waiting for the previous stdout layer and picks it into the next one
func run(layers []*exec.Cmd, pipeWriters []*io.PipeWriter) (err error) {
	// original idea:
	// https://gist.github.com/tyndyll/89fbb2c2273f83a074dc
	for i, layer := range layers {
		if err := layer.Wait(); err != nil {
			return err
		}

		// if next layers in queue exists -> close pipe
		if i < len(layers)-1 {
			if err := pipeWriters[i].Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
