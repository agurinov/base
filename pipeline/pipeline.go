package pipeline

import (
	"io"
)

type Pipeline struct {
	layers []Layer
}

// connect binds all layers of the Pipeline using io.Pipe objects
// connect calls private api's piping method
func (p *Pipeline) connect(input io.ReadCloser, output io.WriteCloser) error {
	return piping(input, output, p.layers.([]Able)...)
}

// run calls private api's piping method
func (p *Pipeline) run() error {
	return run(p.layers.([]Exec)...)
}

func (p *Pipeline) Run(input io.ReadCloser, output io.WriteCloser) error {
	// Stage 1 - piping
	if err := p.connect(input, output); err != nil {
		return err
	}

	// Stage 2 - run
	return p.run()
}
