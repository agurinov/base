package io

import (
	"io"
)

// Pipeable interface describes an object that can be associated with other objects by stdio
type Pipeable interface {
	setStdin(reader io.ReadCloser)
	setStdout(writer io.WriteCloser)
}

// RunCloser interface describes objects that can be self checked and can be executable by Pipeline
type RunCloser interface {
	preRun() error
	Run() error
	check() error
	io.Closer
}

// PipelineLayer interface describes complex type of object that can be a part of Pipeline
type PipeLayer interface {
	Pipeable
	RunCloser
}
