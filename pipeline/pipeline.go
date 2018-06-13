package pipeline

import (
	"errors"
	"io"
)

// Pipeline is set of layers
// it is something like view, workflow for this request
// lifecycle:
// 		1. prepare (piping all layers -> set stdio for each layer)
// 		2. check (check that all layers piped and can be executed)
// 		3. run
// 		4. close (close stdio and clear all layer's sensitive data for reuse this pipeline)
type Pipeline struct {
	layers []Layer
}

// connect binds all layers of the Pipeline using io.Pipe objects
// connect calls private api's piping method
func (p *Pipeline) prepare(input io.ReadCloser, output io.WriteCloser) error {
	// convert Layer -> Able
	layers := make([]Able, len(p.layers))
	// creating []Able with same pointers as p.layers
	for i, layer := range p.layers {
		layers[i] = layer.(Able)
	}

	return piping(input, output, layers...)
}

// check checks all layers can be launched by .Run() at any moment
func (p *Pipeline) check() error {
	// invoke all layer's .check method
	for _, layer := range p.layers {
		if err := layer.check(); err != nil {
			return err
		}
	}

	// all layers is OK and ready for launching
	return nil
}

// run calls private api's run method
func (p *Pipeline) run() error {
	// convert Layer -> Exec
	layers := make([]Exec, len(p.layers))
	// creating []Exec with same pointers as p.layers
	for i, layer := range p.layers {
		layers[i] = layer.(Exec)
	}

	return run(layers...)
}

func (p *Pipeline) close() error {
	return nil
}

func (p *Pipeline) Run(input io.ReadCloser, output io.WriteCloser) error {
	// Stage 1 - piping
	if err := p.prepare(input, output); err != nil {
		return err
	}

	// Stage 2 - check
	if err := p.check(); err != nil {
		return err
	}

	// Stage 3 - run
	// before this stage there is nothing to clean, so we close after
	defer p.close()

	return p.run()
}

func (p *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// inner struct for accepting strings
	var pipeline []map[string]interface{}

	if err := unmarshal(&pipeline); err != nil {
		return err
	}

	// sequece successfully translated, create layers from data
	for _, layer := range pipeline {
		// get required 'type' key
		switch t := layer["type"]; t {
		case "socket":
			socket, err := NewSocket(layer["address"].(string))
			// socket creating errors
			if err != nil {
				return err
			}
			p.layers = append(p.layers, socket)

		case "process":
			process, err := NewProcess("echo FOOBAR")
			// process creating errors
			if err != nil {
				return err
			}
			p.layers = append(p.layers, process)

		default:
			return errors.New("pipeline: Unknown layer.Type")
		}
	}

	return nil
}
