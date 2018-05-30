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

func (c *Process) start() (err error) {
	if c.cmd == nil {
		c.cmd = exec.Command(c.name, c.arg...)
	}

	c.cmd.Stdin = c.stdin
	c.cmd.Stdout = c.stdout

	if c.cmd.Process == nil {
		if err = c.cmd.Start(); err != nil {
			return err
		}
	}

	return nil
}
func (c *Process) setStdin(reader io.ReadCloser) {
	c.stdin = reader
}
func (c *Process) setStdout(writer io.WriteCloser) {
	c.stdout = writer
}
func (c *Process) run() (err error) {
	// check cmd associated and exists
	if err = c.start(); err != nil {
		return err
	}

	return nil
}
func (c *Process) Close() (err error) {
	return c.cmd.Wait()
}
