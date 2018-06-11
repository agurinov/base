package pipeline

import (
	"os/exec"

	"gopkg.in/yaml.v2"
)

type Process struct {
	name string
	arg  []string
	cmd  *exec.Cmd

	stdio
}

func ProcessFromYAML(yml []byte) (*Process, error) {
	var p Process

	err := yaml.Unmarshal(yml, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func NewProcess(name string, arg ...string) *Process {
	return &Process{name: name, arg: arg}
}

func (p *Process) check() error {
	return nil
}

func (p *Process) prepare() error {
	if p.cmd == nil {
		p.cmd = exec.Command(p.name, p.arg...)
	}

	p.cmd.Stdin = p.stdin
	p.cmd.Stdout = p.stdout

	if p.cmd.Process == nil {
		if err := p.cmd.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Process) Run() error {
	return p.cmd.Wait()
}

func (p *Process) Close() error {
	return p.closeStdio()
}
