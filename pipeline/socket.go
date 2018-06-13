package pipeline

import (
	"io"
	"net"
)

type Socket struct {
	address string
	conn    net.Conn

	stdio
}


func NewSocket(address string) *Socket {
	return &Socket{address: address}
}

func (s *Socket) check() error {
	return nil
}

func (s *Socket) prepare() error {
	// TODO resolve address only
	if s.conn == nil {
		conn, err := net.Dial("tcp", s.address)
		if err != nil {
			return err
		}
		s.conn = conn
	}

	return nil
}

func (s *Socket) Run() error {
	// just write to open socket from stdin
	// completes when previous layers stdout closed
	if _, err := io.Copy(s.conn, s.stdin); err != nil {
		return err
	}

	// and receive data as response -> read from connection
	if _, err := io.Copy(s.stdout, s.conn); err != nil {
		return err
	}

	return nil
}

func (s *Socket) Close() error {
	if err := s.closeStdio(); err != nil {
		return err
	}

	// close connection
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}
