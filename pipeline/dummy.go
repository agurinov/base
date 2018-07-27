package pipeline

import (
	"context"
	"errors"
	"io"
)

// FakeLayer is special `fake` Layer
type FakeLayer struct {
	countPrepare int
	countCheck   int
	countRun     int
	countClose   int

	mockFailPrepare bool
	mockFailCheck   bool
	mockFailRun     bool
	mockFailClose   bool
}

func (o *FakeLayer) copy() Layer {
	clone := *o
	return &clone
}

func (o *FakeLayer) prepare() error {
	o.countPrepare++

	if o.mockFailPrepare {
		return errors.New("prepare failed")
	}

	return nil
}
func (o *FakeLayer) check() error {
	o.countCheck++

	if o.mockFailCheck {
		return errors.New("check failed")
	}

	return nil
}
func (o *FakeLayer) run(ctx context.Context) error {
	o.countRun++

	if o.mockFailRun {
		return errors.New("run failed")
	}

	return nil
}
func (o *FakeLayer) close() error {
	o.countClose++

	if o.mockFailClose {
		return errors.New("close failed")
	}

	return nil
}
func (o *FakeLayer) setStdin(reader io.ReadCloser) {
}
func (o *FakeLayer) setStdout(writer io.WriteCloser) {
}
