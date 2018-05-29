package io

import (
	"io"
	// "io/ioutil"
	// "fmt"
	// "log"
	"os/exec"
	// "bytes"
)

type Pipeline []*Layer

func (p *Pipeline) Run() (err error) {
	if err := piping(); err != nil {
		return err
	}

	if err := start(); err != nil {
		return err
	}

	if err := run(); err != nil {
		return err
	}
}

// start the pipeline
// if err := start(layers); err != nil {
// 	return err
// }

// run execution and chaining
// return run(layers, pipeWriters)

// io.ReadWriter
// type Layer struct {
// 	cmd *exec.Cmd
// }
//
// func (l *Layer) Write(p []byte) (n int, err error) {
// 	return l.cmd.Stdout.Write(p)
// }
// func (l *Layer) Read(p []byte) (n int, err error) {
// 	return l.cmd.Stdin.Read(p)
// }
