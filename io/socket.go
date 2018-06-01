package io

import (
	"io"
	"net"
)

type Socket struct {
	address string
	conn    net.Conn

	stdin  io.ReadCloser
	stdout io.WriteCloser
}

func NewSocket(address string) *Socket {
	return &Socket{address: address}
}

func (s *Socket) preRun() (err error) {
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
	// check socket exists
	// just write to open socket from stdin
	if _, err := io.Copy(s.conn, s.stdin); err != nil {
		return err
	}
	if err := s.stdin.Close(); err != nil {
		return err
	}

	// and receive data as response -> read from connection
	if _, err := io.Copy(s.stdout, s.conn); err != nil {
		return err
	}

	if err := s.stdout.Close(); err != nil {
		return err
	}

	return nil
}

// func (s *Socket) preRun() error {
// 	if s.addr == nil {
// 		// trying to resolve TCP address
// 		addr, err := net.ResolveTCPAddr("tcp", s.address)
// 		if err != nil {
// 			return err
// 		}
// 		s.addr = addr
// 	}
//
// 	fmt.Println("net.ResolveTCPAddr()", s.addr)
// 	return nil
// }

func (s *Socket) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *Socket) setStdin(reader io.ReadCloser) {
	s.stdin = reader
}
func (s *Socket) setStdout(writer io.WriteCloser) {
	s.stdout = writer
}
