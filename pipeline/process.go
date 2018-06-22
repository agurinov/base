package pipeline

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/boomfunc/log"
)

type process struct {
	name string
	arg  []string

	cmd *exec.Cmd

	stdio
}

func NewProcess(cmd string) *process {
	parts := strings.Split(cmd, " ")

	return &process{name: parts[0], arg: parts[1:]}
}

func (p *process) prepare() error {
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

// check method guarantees that the object can be launched at any time
// process is piped
func (p *process) check() error {
	// check layer piped
	if err := p.checkStdio(); err != nil {
		return errors.New("pipeline: Process not piped")
	}

	// check command is valid and ready
	if p.cmd == nil {
		return errors.New("pipeline: Process without underlying exec.Cmd")
	}

	// check cmd stdio
	if p.cmd.Stdin == nil || p.cmd.Stdout == nil {
		return errors.New("pipeline: Process's underlying exec.Cmd not piped")
	}

	// process ready for run
	return nil
}

func (p *process) run() error {
	return p.cmd.Wait()
}

func (p *process) close() error {
	log.Debug("process.close()", p.stdin, p.stdout, p.cmd)
	// LIFO reverse order of piping and prepare
	// UNPREPARE
	// reset the command
	// TODO look for better solution
	p.cmd = nil

	// UNPIPING
	// stdio close and reset
	return p.closeStdio()
}
