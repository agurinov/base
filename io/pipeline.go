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
