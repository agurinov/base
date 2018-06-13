package pipeline

import (
	"errors"
	"io"
)

// stdio struct is base struct to something that can have input/output
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
func (obj *stdio) checkStdio() error {
	if obj.stdin == nil {
		return errors.New("pipeline: Able without stdin (Not piped)")
	}

	if obj.stdout == nil {
		return errors.New("pipeline: Able without stdout (Not piped)")
	}

	return nil
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
