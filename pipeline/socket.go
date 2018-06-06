package pipeline

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

func (s *Socket) check() error { return nil }

func (s *Socket) prepare() (err error) {
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
	// close standart input
	// for start layer run and write to stdout
	if err := s.stdin.Close(); err != nil {
		return err
	}

	// close standart output
	// for next layer can complete read from their stdin
	if err := s.stdout.Close(); err != nil {
		return err
	}

	// close connection
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
