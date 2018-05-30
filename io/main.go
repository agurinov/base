package io

import (
	"io"
)

type Pipeable interface {
	setStdin(reader io.ReadCloser)
	setStdout(writer io.WriteCloser)
}

type RunCloser interface {
	run() (err error)
	io.Closer
}

type PipeExec interface {
	Pipeable
	RunCloser
}
