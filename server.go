package main

import (
	"fmt"
	"net"
)

type Server struct {
	listener *net.TCPListener
	clients  []*Client
}

// StartServer starts a new server listening to the given port number. If port is 0, an available
// port number is automatically chosen.
func StartServer(port int) (server *Server, err error) {
	server = &Server{}
	server.listener, err = net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Listening on port %d\n", port)
	go server.accept()
	return server, nil
}

func (s *Server) Stop() {
	for _, c := range s.clients {
		c.Stop()
	}
	s.listener.Close()
}

func (s *Server) accept() {
	for {
		conn, err := s.listener.AcceptTCP()
		if conn == nil {
			fmt.Printf("Stopped Listener\n")
			break
		}
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			conn.Close()
			continue
		}
		s.NewClient(conn)
	}
}

func (s *Server) NewClient(conn *net.TCPConn) {
	if conn == nil {
		return
	}

	c := Client{conn: conn, state: STATEHANDSHAKE, server: s}
	go c.listen()
	s.clients = append(s.clients, &c)
}
