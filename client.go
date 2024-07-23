package main

import (
	"errors"
	"fmt"
	"io"
	"mcgo/types"
	"net"
)

type Client struct {
	conn  *net.TCPConn
	state CommunicationState

	protocolVersion   uint32
	connectionAddress string
	connectionPort    uint16
}

func (c *Client) Stop() {
	addr := c.conn.RemoteAddr()
	err := c.conn.Close()
	if err == nil {
		fmt.Printf("Stopped Client %s\n", addr)
	}
}

func (c *Client) listen() {
	defer c.Stop()
	for {
		length, err := types.ReadVarInt(c.conn)
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			break
		} else if err != nil {
			fmt.Printf("Error on reading package length: %v\n", err)
			break
		}
		c.parsePackage(length)
	}
}

func (c *Client) parsePackage(length uint32) {
	if length == 0 {
		return
	}

	data, err := types.Read(int(length), c.conn)
	if err != nil {
		fmt.Printf("Error on reading %d bytes of package data: %v", length, err)
		return
	}
	if int(length) != len(data) {
		fmt.Printf("Error tried to read %d bytes of packet data, but got %d bytes\ndata: % 02x\nstring: %s\n\n", length, len(data), data, data)
		c.Stop()
		return
	}

	packetID, err := types.PopVarInt(&data)
	if err != nil {
		fmt.Printf("Error on packetID: %v", err)
		return
	}

	switch c.state {
	case STATEHANDSHAKE:
		if packetID == PACKETHANDSHAKE {
			c.handleHandshake(data)
			return
		}
	case STATESTATUS:
		switch packetID {
		case PACKETSTATUSREQUEST:
			fmt.Printf("/*TODO status request*/\n")
			c.Stop()
			return
		case PACKETPINGREQUEST:
			fmt.Printf("/*TODO ping request*/\n")
			c.Stop()
			return
		}
	case STATELOGIN:
		fmt.Printf("Client %s tried to login, but isnt supportet yet!\n", c.conn.RemoteAddr())
		c.Stop()
		return
	}
	fmt.Printf("Unknown packet 0x%02x in state 0x%02x: dropping connection\npacket data: % 02x\nstring: %s\n\n", packetID, c.state, data, string(data))
	c.Stop()
}

func (c *Client) handleHandshake(data []byte) {
	var err error
	c.protocolVersion, err = types.PopVarInt(&data)
	if err != nil {
		fmt.Printf("Error reading protocol version: %v\n", err)
		c.Stop()
		return
	}

	c.connectionAddress, err = types.PopString(&data)
	if err != nil {
		fmt.Printf("Error reading server address: %v\n", err)
		c.Stop()
		return
	}

	c.connectionPort, err = types.PopUShort(&data)
	if err != nil {
		fmt.Printf("Error reading server port: %v\n", err)
		c.Stop()
		return
	}

	n, err := types.PopVarInt(&data)
	if err != nil {
		fmt.Printf("Error reading next state: %v\n", err)
		c.Stop()
		return
	}
	c.state = CommunicationState(n)

	fmt.Printf("Shaked hands with %s - switched to state 0x%02x with version %d on %s:%d\n", c.conn.RemoteAddr(), c.state, c.protocolVersion, c.connectionAddress, c.connectionPort)
}
