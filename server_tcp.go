package main

import (
	"io"
	"net"
	"os"
)

const (
	defaultLogFormat = "pretty"
	defaultPort      = "8080"
)

type TCPServer struct {
	Port string
}

func (s *TCPServer) Init() {
	s.Port = os.Getenv("FLIES_PORT")
}

func (s *TCPServer) Listen(out, errOut io.Writer) error {
	l, err := net.Listen("tcp", s.getAddr())
	if err != nil {
		s.writeError(errOut, err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			s.writeError(errOut, err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			io.Copy(out, c)
			c.Close()
		}(conn)
	}
	return err
}

func (s *TCPServer) writeError(out io.Writer, err error) error {
	out.Write([]byte(err.Error()))
	out.Write([]byte("\n"))
	return nil
}

func (s *TCPServer) getAddr() string {
	port := s.Port
	if port == "" {
		port = defaultPort
	}
	return net.JoinHostPort("", port)
}
