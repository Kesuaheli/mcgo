package main

import (
	"fmt"
	"mcgo/world"
	"net"

	"github.com/google/uuid"
)

type Server struct {
	listener *net.TCPListener
	clients  []*Client

	entityCount int32
}

var World *world.World = world.NewWorld()

// StartServer starts a new server listening to the given port number. If port is 0, an available
// port number is automatically chosen.
func StartServer(port int) (server *Server, err error) {
	World.GetChunk(0, 0).GetChunkSection(0).FillWithBlocks(0, 0, 0, 15, 0, 15, "minecraft:cobblestone")
	World.GetChunk(0, 0).GetChunkSection(0).FillWithBlocks(0, 1, 0, 15, 2, 15, "minecraft:dirt")
	World.GetChunk(0, 0).GetChunkSection(0).FillWithBlocks(0, 3, 0, 15, 3, 15, "minecraft:grass_block")
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

func (s *Server) NewEntityID() (eID int32, UUID uuid.UUID) {
	s.entityCount++
	if s.entityCount < 1 {
		panic("entity count overflow")
	}
	return s.entityCount - 1, uuid.New()
}
