package io

import (
	"io"
	"os/exec"
)

type Process struct {
	name string
	arg  []string
	cmd  *exec.Cmd

	stdin  io.ReadCloser
	stdout io.WriteCloser
}

func NewProcess(name string, arg ...string) *Process {
	return &Process{name: name, arg: arg}
}

func (p *Process) preRun() error {
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
	if err := p.cmd.Wait(); err != nil {
		return err
	}

	if err := p.stdin.Close(); err != nil {
		return err
	}

	if err := p.stdout.Close(); err != nil {
		return err
	}

	return nil
}

func (p *Process) postRun() error {
	return nil
}

func (p *Process) Close() error {
	return nil
}

func (p *Process) setStdin(reader io.ReadCloser) {
	p.stdin = reader
}

func (p *Process) setStdout(writer io.WriteCloser) {
	p.stdout = writer
}
