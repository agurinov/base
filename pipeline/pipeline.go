package pipeline

import (
	"errors"
	"io"
)

type Pipeline []Layer

func New(layers ...Layer) *Pipeline {
	p := Pipeline(layers)
	return &p
}

func (p Pipeline) Run(input io.Reader, output io.Writer) error {
	// Convert io.Reader and io.Writer to io.ReadCloser and io.WriteCloser
	inputCloser := toReadCloser(input)
	outputCloser := toWriteCloser(output)

	ables := make([]Able, len(p))
	execs := make([]Exec, len(p))

	for i, layer := range p {
		newLayer := layer.(Cloneable).copy()

		ables[i] = newLayer.(Able)
		execs[i] = newLayer.(Exec)
	}

	if err := piping(inputCloser, outputCloser, ables...); err != nil {
		return err
	}

	return run(execs...)
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
			*p = append(*p, NewTCPSocket(layer["address"].(string)))

		case "process":
			*p = append(*p, NewProcess(layer["cmd"].(string)))

		default:
			return errors.New("pipeline: Unknown layer type")
		}
	}

	return nil
}
