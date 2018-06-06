package pipeline

import (
	"io"
)

// Able interface describes an object that can be associated with other objects by stdio
type Able interface {
	setStdin(reader io.ReadCloser)
	setStdout(writer io.WriteCloser)
}

// Exec interface describes objects that can be self checked and can be executable by Pipeline
type Exec interface {
	check() error
	preRun() error
	Run() error
	io.Closer
}

// Layer interface describes complex type of object that can be a part of Pipeline
type Layer interface {
	Able
	Exec
}
