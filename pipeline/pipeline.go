package pipeline

import (
	"errors"
	"io"
)

type Pipeline struct {
	layers []Layer
}

// connect binds all layers of the Pipeline using io.Pipe objects
// connect calls private api's piping method
func (p *Pipeline) connect(input io.ReadCloser, output io.WriteCloser) error {
	// convert Layer -> Able
	layers := make([]Able, len(p.layers))
	// creating []Able with same pointers as p.layers
	for i, layer := range p.layers {
		layers[i] = layer.(Able)
	}

	return piping(input, output, layers...)
}

// run calls private api's piping method
func (p *Pipeline) run() error {
	// convert Layer -> Exec
	layers := make([]Exec, len(p.layers))
	// creating []Exec with same pointers as p.layers
	for i, layer := range p.layers {
		layers[i] = layer.(Exec)
	}

	return run(layers...)
}

func (p *Pipeline) Run(input io.ReadCloser, output io.WriteCloser) error {
	// Stage 1 - piping
	if err := p.connect(input, output); err != nil {
		return err
	}

	// Stage 2 - run
	return p.run()
}

func (p *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// inner struct for accepting strings
	var pipeline []struct{ Type string }

	if err := unmarshal(&pipeline); err != nil {
		return err
	}

	// sequece successfully translated, create layers from data
	for _, layer := range pipeline {
		switch layer.Type {
		case "socket":
			p.layers = append(p.layers, NewSocket("golang.org"))
		case "process":
			p.layers = append(p.layers, NewProcess("echo", "FOOBAR"))
		default:
			return errors.New("pipeline: Unknown layer.Type")
		}
	}

	return nil
}
