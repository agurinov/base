package pipeline

import (
	"errors"
	"io"
	"net"
)

type tcp struct {
	address string

	addr *net.TCPAddr
	conn *net.TCPConn

	stdio
}

func NewTCPSocket(address string) *tcp {
	return &tcp{address: address}
}

func (s *tcp) prepare() error {
	if s.addr == nil {
		var err error

		// try to resolve remote address
		// TODO too much time when fake host
		if s.addr, err = net.ResolveTCPAddr("tcp", s.address); err != nil {
			return err
		}
	}

	return nil
}

// check method guarantees that the object can be launched at any time
// tcp socket is piped
// remote address resolvable
func (s *tcp) check() error {
	// check layer piped
	if err := s.checkStdio(); err != nil {
		return err
	}

	// check tcp socket have real address
	if s.addr == nil {
		return errors.New("pipeline: TCP socket without address")
	}

	// TODO check empty input

	// tcp socket ready for run
	return nil
}

// BUG: when tcp socket anywhere as stdin and stdout is the same conn -> blocking when io.Copy
func (s *tcp) run() error {
	var err error

	// establish tcp socket
	if s.conn, err = net.DialTCP("tcp", nil, s.addr); err != nil {
		return err
	}

	// just write to open tcp socket from stdin
	// completes when previous layers stdout closed
	// TODO handle empty input here
	if _, err = io.Copy(s.conn, s.stdin); err != nil {
		return err
	}

	// and receive data as response -> read from connection
	if _, err = io.Copy(s.stdout, s.conn); err != nil {
		return err
	}

	return nil
}

func (s *tcp) close() (err error) {
	defer func() {
		s.conn = nil
	}()

	defer func() {
		if err != nil {
			s.closeStdio()
		} else {
			err = s.closeStdio()
		}
	}()

	// close connection
	err = s.conn.Close()
	return
}
