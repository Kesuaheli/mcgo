package main

import (
	"bytes"
	"fmt"
	"mcgo/types"
)

type Sender struct {
	buf *bytes.Buffer
}

func NewSender(packet Packet) *Sender {
	s := &Sender{}
	s.buf = &bytes.Buffer{}
	s.buf.WriteByte(byte(packet))
	return s
}

// Send writes the packet to the client.
func (s *Sender) Send(c *Client, name string) (err error) {
	b := &bytes.Buffer{}
	types.WriteVarInt(b, int32(s.buf.Len()))
	b.Write(s.buf.Bytes())

	_, err = c.conn.Write(b.Bytes())
	if err != nil {
		fmt.Printf("Failed to send packet %s: %v\n", name, err)
		c.Disconnect(name+" failed", "red")
	}
	return err
}

func (s *Sender) SendMultipleClients(clients []*Client, name string) (hadError bool) {

	for _, c := range clients {
		if !c.connected {
			continue
		}
		_, err := c.conn.Write(s.Bytes())
		if err != nil {
			fmt.Printf("Failed to send packet %s: %v\n", name, err)
			c.Disconnect(name+" failed", "red")
			hadError = true
		}
	}

	return hadError
}

// Write implements [io.Writer]
func (s *Sender) Write(p []byte) (n int, err error) {
	return s.buf.Write(p)
}

// Debugging

func (s *Sender) Bytes() (data []byte) {
	b := &bytes.Buffer{}
	types.WriteVarInt(b, int32(s.buf.Len()))
	b.Write(s.buf.Bytes())

	return b.Bytes()
}
