package io

import (
	"io"
	"net"
)

type Socket struct {
	address string
	conn    *net.Conn

	stdin  io.ReadCloser
	stdout io.WriteCloser
}

func (s *Socket) start() (err error) {
	if s.conn == nil {
		if s.conn, err = net.Dial("tcp", s.address); err != nil {
			return err
		}
	}

	return nil
}
func (s *Socket) pipe(reader io.ReadCloser, writer io.WriteCloser) {
	s.stdin = reader
	s.stdout = writer
}
func (s *Socket) run() (err error) {
	// check socket exists
	if err = s.start(); err != nil {
		return err
	}

	// just write to open socket from stdin
	if _, err := io.Copy(s.conn, s.stdin); err != nil {
		return err
	}
	if s.stdin.Close(); err != nil {
		return err
	}

	// and receive data as response -> read from connection
	if _, err := io.Copy(s.stdout, s.conn); err != nil {
		return err
	}
	if s.stdout.Close(); err != nil {
		return err
	}

	return nil
}
func (s *Socket) close() (err error) {
	return s.conn.Close()
}
