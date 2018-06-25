package pipeline

import (
	"errors"
	"io"
)

var ErrNoStdin = errors.New("pipeline: Able without stdin (Not piped)")
var ErrNoStdout = errors.New("pipeline: Able without stdout (Not piped)")

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
		return ErrNoStdin
	}

	if obj.stdout == nil {
		return ErrNoStdout
	}

	return nil
}
func (obj *stdio) closeStdio() (err error) {

	defer func() {
		obj.stdin = nil
		obj.stdout = nil
	}()

	// close standart input
	// for start layer run and write to stdout
	if obj.stdin != nil {
		err = obj.stdin.Close()
	} else {
		err = ErrNoStdin
	}

	// close standart output
	// for next layer can complete read from their stdin
	if obj.stdout != nil {
		err = obj.stdout.Close()
	} else {
		err = ErrNoStdout
	}

	return
}
