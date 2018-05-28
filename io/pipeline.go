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

// start the pipeline
// if err := start(layers); err != nil {
// 	return err
// }

// run execution and chaining
// return run(layers, pipeWriters)

type Layer struct {
	cmd *exec.Cmd
}

func (l *Layer) Write(p []byte) (n int, err error) {
	return l.cmd.Stdout.Write(p)
}
func (l *Layer) Read(p []byte) (n int, err error) {
	return l.cmd.Stdin.Read(p)
}
