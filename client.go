package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mcgo/types"
	"net"

	"github.com/google/uuid"
)

type Client struct {
	conn  *net.TCPConn
	state CommunicationState

	protocolVersion   int32
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

func (c *Client) parsePackage(length int32) {
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
			c.handleStatusRequest()
			return
		case PACKETPINGREQUEST:
			c.handlePingRequest(data)
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

func (c *Client) handleStatusRequest() {
	status := StatusResponse{
		Version: StatusResponseVersion{
			Name:     "Hello there",
			Protocol: c.protocolVersion - 1,
		},
		Players: StatusResponsePlayers{
			Max:    69,
			Online: 1337,
			Sample: []StatusResponsePlayer{
				{
					Name: "Dont look at me!",
					ID:   uuid.New(),
				},
			},
		},
		Description: types.TextComponent(fmt.Sprintf("You pinged: §2%s:%d§r\nYour IP is: §6%s§r", c.connectionAddress, c.connectionPort, c.conn.RemoteAddr())),
	}

	data, err := json.Marshal(status)
	if err != nil {
		fmt.Printf("Failed to marshal status response: %v\n", err)
		c.Stop()
		return
	}

	buf := &bytes.Buffer{}
	buf.WriteByte(PACKETSTATUSRESPONSE)
	err = types.WriteStringData(buf, data)

	types.WriteVarInt(c.conn, int32(buf.Len()))
	if err != nil {
		fmt.Printf("Failed to send status response length %d: %v\n", buf.Len(), err)
		c.Stop()
		return
	}

	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Failed to send status response: %v\n", err)
		c.Stop()
		return
	}

	fmt.Printf("Send status response to %s\n", c.conn.RemoteAddr())
}

func (c *Client) handlePingRequest(data []byte) {
	pingID, err := types.PopLong(&data)
	if err != nil {
		fmt.Printf("Error reading ping data: %v\n", err)
		c.Stop()
		return
	}

	err = types.WriteVarInt(c.conn, 9)
	if err != nil {
		fmt.Printf("Failed to send ping response (packet length): %v\n", err)
		c.Stop()
		return
	}

	_, err = c.conn.Write([]byte{PACKETPINGRESPONSE})
	if err != nil {
		fmt.Printf("Failed to send ping response (packet id): %v\n", err)
		c.Stop()
		return
	}

	err = types.WriteLong(c.conn, pingID)
	if err != nil {
		fmt.Printf("Failed to send ping response (ping id): %v\n", err)
		c.Stop()
		return
	}

	fmt.Printf("Send ping response to %s: %d\n", c.conn.RemoteAddr(), pingID)
}
