package pipeline

import (
	"io"
	"net"

	"gopkg.in/yaml.v2"
)

type Socket struct {
	Address string `yaml:"address,omitempty"`
	conn    net.Conn

	stdio
}

func SocketFromYAML(yml []byte) (*Socket, error) {
	var s Socket

	err := yaml.Unmarshal(yml, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func NewSocket(address string) *Socket {
	return &Socket{Address: address}
}

func (s *Socket) check() error {
	return nil
}

func (s *Socket) prepare() error {
	// TODO resolve address only
	if s.conn == nil {
		conn, err := net.Dial("tcp", s.Address)
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
