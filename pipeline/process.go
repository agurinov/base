package pipeline

import (
	"errors"
	"os/exec"
	"strings"
)

type process struct {
	name string
	arg  []string

	cmd *exec.Cmd

	stdio
}

func NewProcess(cmd string) *process {
	command := strings.Split(cmd, " ")

	return &process{name: command[0], arg: command[1:]}
}

func (p *process) prepare() error {
	if p.cmd == nil {
		p.cmd = exec.Command(p.name, p.arg...)
	}

	p.cmd.Stdin = p.stdin
	p.cmd.Stdout = p.stdout

	// if p.cmd.Process == nil {
	// 	if err := p.cmd.Start(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// check method guarantees that the object can be launched at any time
// process is piped
func (p *process) check() error {
	// check layer piped
	if err := p.checkStdio(); err != nil {
		return errors.New("pipeline: Process not piped")
	}

	// check command is valid and ready
	if p.cmd == nil {
		return errors.New("pipeline: Process without exec.Cmd")
	}

	// process ready for run
	return nil
}

func (p *process) run() error {
	return p.cmd.Wait()
}

func (p *process) close() error {
	return p.closeStdio()
}
