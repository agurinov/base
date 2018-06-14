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

	stdio
}

// connect binds all layers of the Pipeline using io.Pipe objects
// connect calls private api's piping method
func (p *Pipeline) prepare() error {
	// Stage 1. Pipeline prepare -> piping
	// convert Layer -> Able
	layers := make([]Able, len(p.layers))
	// creating []Able with same pointers as p.layers
	for i, layer := range p.layers {
		layers[i] = layer.(Able)
	}
	if err := piping(p.stdin, p.stdout, layers...); err != nil {
		return err
	}

	// Stage 2. Prepare all layers
	// invoke all layer's .prepare method
	for _, layer := range p.layers {
		if err := layer.prepare(); err != nil {
			return err
		}
	}

	return nil
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

	// run this layers (Execs) throw tool
	return run(layers...)
}

func (p *Pipeline) close() error {
	return nil
}

func (p *Pipeline) Run(input io.ReadCloser, output io.WriteCloser) error {
	// Piping this Able with input and output
	// save request and response to inner data for prepare and piping internal layers
	if err := piping(input, output, p); err != nil {
		return err
	}

	// run this Exec throw tool
	return run(p)
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
		case "tcp":
			p.layers = append(p.layers, NewTCPSocket(layer["address"].(string)))

		case "process":
			p.layers = append(p.layers, NewProcess(layer["cmd"].(string)))

		default:
			return errors.New("pipeline: Unknown layer type")
		}
	}

	return nil
}
