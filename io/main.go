package io

import (
	"io"
)

type Pipeable interface {
	setStdin(reader io.ReadCloser)
	setStdout(writer io.WriteCloser)
}

type Checker interface {
	check() error
}

// TODO rename
type RunCloser interface {
	preRun() error
	Run() error
	// Checker
	io.Closer
}

type PipeExec interface {
	Pipeable
	RunCloser
}
