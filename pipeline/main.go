package pipeline

import (
	"io"
)

// StdIO struct is base struct to something that can have input/output
// automatically implements pipeline.Able interface
// must be inherited in some objects like Socket and Process
type stdio struct {
	stdin  io.ReadCloser
	stdout io.WriteCloser
}

func (obj *stdio) setStdin(reader io.ReadCloser) {
	obj.stdin = reader
}
func (obj *stdio) setStdout(writer io.WriteCloser) {
	obj.stdout = writer
}
func (obj *stdio) closeStdio() error {
	// close standart input
	// for start layer run and write to stdout
	if err := obj.stdin.Close(); err != nil {
		return err
	}

	// close standart output
	// for next layer can complete read from their stdin
	if err := obj.stdout.Close(); err != nil {
		return err
	}

	return nil
}

// Able interface describes an object that can be associated with other objects by stdio
type Able interface {
	setStdin(reader io.ReadCloser)
	setStdout(writer io.WriteCloser)
}

// Exec interface describes objects that can be self checked and can be executable by Pipeline
type Exec interface {
	check() error
	prepare() error
	Run() error
	io.Closer
}

// Layer interface describes complex type of object that can be a part of Pipeline
type Layer interface {
	Able
	Exec
}
