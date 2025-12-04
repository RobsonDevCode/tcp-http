package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"tcp-http/Internal/contracts/response"
)

type Server struct {
	closed bool
}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{closed: false}
	go func() {
		err := runServer(server, listener)
		if err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	return server, nil
}

func runConnection(_ *Server, connection io.ReadWriteCloser) error {
	defer connection.Close()
	headers := response.WithDefaultHeaders(0)

	if err := response.WriteStatusLine(connection, response.StatusOK); err != nil {
		return errors.Join(fmt.Errorf("error writing status line"), err)
	}

	if err := response.WriteHeaders(connection, headers); err != nil {
		return errors.Join(fmt.Errorf("error writing headers", err))
	}

	return nil
}
func runServer(s *Server, listener net.Listener) error {
	go func() {
		for {
			conn, err := listener.Accept()
			if s.closed {
				return
			}

			if err != nil {
				return
			}

			go func() {
				err := runConnection(s, conn)
				if err != nil {
					log.Fatalf("connection failed")
				}
			}()
		}
	}()

	return nil
}
func (s *Server) Close() error {
	s.closed = true
	return nil
}
