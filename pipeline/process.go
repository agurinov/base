package pipeline

import (
	"context"
	// "errors"
	"os/exec"
	"strings"

	"github.com/boomfunc/base/tools"
)

type process struct {
	cmd string
	// cmd exec.Cmd

	stdio
}

func NewProcess(cmd string) *process {
	return &process{cmd: cmd}
}

func (p *process) copy() Layer {
	clone := *p
	return &clone
}

func (p *process) prepare() error {
	// p.cmd.Stdin = p.stdin
	// p.cmd.Stdout = p.stdout
	//
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
		return err
	}

	// // check command is valid and ready
	// if p.cmd == nil {
	// 	return errors.New("pipeline: Process without underlying exec.Cmd")
	// }

	// check cmd stdio
	// if p.cmd.Stdin == nil || p.cmd.Stdout == nil {
	// 	return errors.New("pipeline: Process's underlying exec.Cmd not piped")
	// }

	// process ready for run
	return nil
}

func (p *process) run(ctx context.Context) error {
	// fill templates from ctx
	parts := strings.Split(
		tools.StringFromCtx(ctx, p.cmd), " ",
	)

	cmd := exec.Command(parts[0], parts[1:]...)

	cmd.Stdin = p.stdin
	cmd.Stdout = p.stdout

	if cmd.Process == nil {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	return cmd.Wait()
}

func (p *process) close() error {
	// defer func() {
	// 	p.cmd = nil
	// }()

	return p.closeStdio()

	// defer func() {
	// 	if err != nil {
	// 		p.closeStdio()
	// 	} else {
	// 		err = p.closeStdio()
	// 	}
	// }()
	//
	// return
}
