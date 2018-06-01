package io

import (
	"io"
)

type Pipeline struct {
	layers []PipeLayer
}

// connect binds all layers of the Pipeline using io.Pipe objects
func (p *Pipeline) connect(input io.ReadCloser, output io.WriteCloser) error {
	return piping(input, output, p.layers.([]Pipeable)...)
}

func (p *Pipeline) run() error {
	return run(p.layers.([]RunCloser)...)
}

func (p *Pipeline) Run(input io.ReadCloser, output io.WriteCloser) error {
	// Stage 1 - piping
	if err := p.connect(input, output); err != nil {
		return err
	}

	// Stage 2 - run
	return p.run()
}
