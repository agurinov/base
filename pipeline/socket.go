package pipeline

import (
	"errors"
	"io"
	"net"
)

type socket struct {
	address string

	addr *net.TCPAddr
	conn *net.TCPConn

	stdio
}

func NewSocket(address string) *socket {
	return &socket{address: address}
}

func (s *socket) prepare() error {
	// try to resolve remote address
	if s.addr, err = net.ResolveTCPAddr("tcp", s.address); err != nil {
		return err
	}

	return nil
}

// check method guarantees that the object can be launched at any time
// socket is piped
// remote address resolvable
func (s *socket) check() error {
	// check layer piped
	if err := s.checkStdio(); err != nil {
		return errors.New("pipeline: Socket not piped")
	}

	// schek socket have real address
	if s.addr == nil {
		return errors.New("pipeline: Socket without address")
	}

	// socket ready for run
	return nil
}

func (s *socket) run() error {
	// establish socket
	if s.conn == nil {
		if s.conn, err = net.DialTCP("tcp", nil, s.addr); err != nil {
			return err
		}
	}

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

func (s *socket) close() error {
	if err := s.closeStdio(); err != nil {
		return err
	}

	// close connection
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}
