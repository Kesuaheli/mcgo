package main

import (
	"bytes"
	"encoding/binary"
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

	name   string
	uuid   uuid.UUID
	config ClientConfiguration
}

type ClientConfiguration struct {
	lang               string
	viewDistance       int8
	chatMode           int32
	chatColors         bool
	displayedSkinParts uint8
	mainHand           int32
	textFiltering      bool
	inServerListing    bool
}

var switch_to_play = true

func (c Client) String() string {
	if c.name != "" {
		return c.name
	} else if c.uuid != (uuid.UUID{}) {
		return c.uuid.String()
	}
	return c.conn.RemoteAddr().String()
}

func (c *Client) Stop() {
	err := c.conn.Close()
	if err == nil {
		fmt.Printf("Stopped Client %s\n", c)
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

	packetNum, err := types.PopVarInt(&data)
	if err != nil {
		fmt.Printf("Error on packetID: %v", err)
		return
	}
	packetID := Packet(packetNum)

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
		switch packetID {
		case PACKETLOGINSTART:
			c.handleLoginStart(data)
			return
		case PACKETLOGINACKNOWLEDGED:
			c.handleLoginAcknowledgement()
			return
		}
	case STATECONFIGURATION:
		switch packetID {
		case PACKETCLIENTINFORMATION:
			c.handleClientInformation(data)
			return
		case PACKETPLUGINMESSAGE:
			c.handlePluginMessage(data)
			return
		case PACKETFINISHCONFIGURATION:
			c.state = STATEPLAYING
			fmt.Printf("Finished configuration for %s switched to playing state\n", c.name)
			c.sendLogin()
			c.sendPlayerPosition()
			return
		}
	case STATEPLAYING:
		switch packetID {
		case PACKETACCEPTTELEPORT:
			c.sendChunkCenter()
			c.sendChunkData()
		}
	}
	fmt.Printf("Unknown packet 0x%02x in state %s: dropping connection\npacket data: % 02x\nstring: %s\n\n", packetID, c.state, data, string(data))
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
	if n == 3 {
		n = STATELOGIN
	}
	c.state = CommunicationState(n)

	fmt.Printf("Shaked hands with %s - switched to state %s with version %d on %s:%d\n", c, c.state, c.protocolVersion, c.connectionAddress, c.connectionPort)
}

func (c *Client) handleStatusRequest() {
	status := StatusResponse{
		Version: StatusResponseVersion{
			Name:     "Hello there",
			Protocol: c.protocolVersion,
		},
		Players: StatusResponsePlayers{
			Max:    69,
			Online: 1337,
			Sample: []StatusResponsePlayer{
				{
					Name: "Dont look at me!",
					ID:   uuid.New(),
				},
				{
					Name: "ja",
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
	buf.WriteByte(byte(PACKETSTATUSRESPONSE))
	types.WriteStringData(buf, data)

	err = types.WriteVarInt(c.conn, int32(buf.Len()))
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

	fmt.Printf("Send status response to %s\n", c)
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

	_, err = c.conn.Write([]byte{byte(PACKETPINGRESPONSE)})
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

	fmt.Printf("Send ping response to %s: %d\n", c, pingID)
}

func (c *Client) handleLoginStart(data []byte) {
	var err error
	c.name, err = types.PopString(&data)
	if err != nil {
		fmt.Printf("Error reading player name data: %v\n", err)
		c.Stop()
		return
	}
	c.uuid, err = types.PopUUID(&data)
	if err != nil {
		fmt.Printf("Error reading player uuid data: %v\n", err)
		c.Stop()
		return
	}

	fmt.Printf("%s (%s) wants to login to the server.\n", c.name, c.uuid)
	//c.Disconnect("Kommscht hier net rein!", "light_purple")
	c.sendLoginSuccess()
}

func (c *Client) Disconnect(reason, color string) {
	buf := &bytes.Buffer{}
	buf.WriteByte(byte(PACKETDISCONNECT))
	types.WriteString(buf, fmt.Sprintf("{\"text\":\"%s\",\"color\":\"%s\"}", reason, color))

	err := types.WriteVarInt(c.conn, int32(buf.Len()))
	if err != nil {
		fmt.Printf("[ERROR] Failed to disconnect client properly: %v\n", err)
		return
	}
	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("[ERROR] Failed to disconnect client properly: %v\n", err)
		return
	}

	fmt.Printf("Wrote Disconnect message to %s: %s\n", c, string(buf.Bytes()))
}

func (c *Client) sendLoginSuccess() {
	buf := &bytes.Buffer{}
	buf.WriteByte(byte(PACKETLOGINSUCCESS))
	buf.Write(c.uuid[:])
	types.WriteString(buf, c.name)
	const properties = 1
	types.WriteVarInt(buf, properties)
	for i := 0; i < properties; i++ {
		types.WriteString(buf, "textures")
		types.WriteString(buf, "ewogICJ0aW1lc3RhbXAiIDogMTcyMjA3NzM2Njk5OCwKICAicHJvZmlsZUlkIiA6ICIwNDJmNDdkNTlhM2M0Yzk4OWE1MGM3MWYzOGYzOGVkMCIsCiAgInByb2ZpbGVOYW1lIiA6ICJLZXN1YWhlbGkiLAogICJ0ZXh0dXJlcyIgOiB7CiAgICAiU0tJTiIgOiB7CiAgICAgICJ1cmwiIDogImh0dHA6Ly90ZXh0dXJlcy5taW5lY3JhZnQubmV0L3RleHR1cmUvODJmMzNiNmZlM2JiMzZiODg1NTkwZjM3NzZkYWUxNDRmODMyMzM2ZmM5NmJkNGJjYzcxYzUxYWI5ZjM1YmQyIgogICAgfSwKICAgICJDQVBFIiA6IHsKICAgICAgInVybCIgOiAiaHR0cDovL3RleHR1cmVzLm1pbmVjcmFmdC5uZXQvdGV4dHVyZS9hZmQ1NTNiMzkzNThhMjRlZGZlM2I4YTlhOTM5ZmE1ZmE0ZmFhNGQ5YTljM2Q2YWY4ZWFmYjM3N2ZhMDVjMmJiIgogICAgfQogIH0KfQ==")
		buf.WriteByte(0) // signature
	}

	err := types.WriteVarInt(c.conn, int32(buf.Len()))
	if err != nil {
		fmt.Printf("Failed to write send login success packet: %v\n", err)
		c.Disconnect("Login failed", "red")
		return
	}
	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Failed to write send login success packet: %v\n", err)
		c.Disconnect("Login failed", "red")
		return
	}

	fmt.Printf("%s logged in\n", c)
}

func (c *Client) handleLoginAcknowledgement() {
	c.state = STATECONFIGURATION
}

func (c *Client) handlePluginMessage(data []byte) {
	channel, err := types.PopString(&data)
	if err != nil {
		fmt.Printf("failed to read plugin message: %v\n", err)
		c.Stop()
		return
	}

	fmt.Printf("Got plugin message %s: % 02x (%s)\n", channel, data, string(data))
}

func (c *Client) handleClientInformation(data []byte) {
	var err error
	c.config.lang, err = types.PopString(&data)
	if err != nil {
		fmt.Printf("Failed to read client lang: %v\n", err)
		c.Stop()
		return
	}

	c.config.viewDistance = int8(data[0])
	data = data[1:]

	c.config.chatMode, err = types.PopVarInt(&data)
	if err != nil {
		fmt.Printf("Failed to read client chat mode: %v\n", err)
		c.Stop()
		return
	}

	c.config.chatColors = data[0] == 1
	c.config.displayedSkinParts = data[1]
	data = data[2:]

	c.config.mainHand, err = types.PopVarInt(&data)
	if err != nil {
		fmt.Printf("Failed to read client main hand: %v\n", err)
		c.Stop()
		return
	}

	c.config.textFiltering = data[0] == 1
	c.config.inServerListing = data[1] == 1

	fmt.Printf("Client config: %+v\n", c.config)

	if switch_to_play {
		_, err = c.conn.Write([]byte{0x01, byte(PACKETFINISHCONFIGURATION)})
		if err != nil {
			fmt.Printf("Failed to switch %s to play state: %v\n", c, err)
			c.Stop()
			return
		}
		c.conn.Write([]byte{0x00, 0x00, 0x00})
	} else {
		// we want to transfer
		c.sendTranfer()
		switch_to_play = true
	}

}

func (c *Client) sendTranfer() {
	buf := &bytes.Buffer{}
	switch c.state {
	case STATECONFIGURATION:
		buf.WriteByte(byte(PACKETTRANSFERCONFIGURATION))
	case STATEPLAYING:
		buf.WriteByte(byte(PACKETTRANSFERPLAYING))
	}

	types.WriteString(buf, "localhost")
	types.WriteVarInt(buf, 25566)

	err := types.WriteVarInt(c.conn, int32(buf.Len()))
	if err != nil {
		fmt.Printf("Failed to write send transfer packet: %v\n", err)
		c.Disconnect("Transfer failed", "red")
		return
	}
	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Failed to write send transfer packet: %v\n", err)
		c.Disconnect("Transfer failed", "red")
		return
	}

	fmt.Printf("%s transfer package send\n", c)
}

func (c *Client) sendLogin() {
	buf := &bytes.Buffer{}

	buf.WriteByte(byte(PACKETLOGINPLAYING))

	// Write Player Entity Identity
	binary.Write(buf, binary.BigEndian, int32(1))

	// Send Dimension Data
	types.WriteVarInt(buf, 1) // 1 dimension
	types.WriteString(buf, "minecraft:overworld")

	// Max Players
	types.WriteVarInt(buf, 1) // 1 Player

	// View Distance
	types.WriteVarInt(buf, 2) // min 2, max 32

	// Simulation Distance
	types.WriteVarInt(buf, 1) // lets try it

	// Reduced Debug Info
	types.WriteBoolean(buf, false)

	// Enable Respawn screen
	types.WriteBoolean(buf, true)

	// Do limited Crafting
	types.WriteBoolean(buf, false)

	// Dimension Type
	types.WriteVarInt(buf, 0) // hopefully overworld

	// Dimension Name
	types.WriteString(buf, "minecraft:overworld")

	// Hashed Seed
	types.WriteLong(buf, 0x12345678)

	// Game Mode
	buf.WriteByte(0) // 0: Survival, 1: Creative, 2: Adventure, 3: Spectator.

	// Previous Game Mode
	buf.WriteByte(1) // All bytes set = -1

	// Is Debug
	types.WriteBoolean(buf, true) // switch off later

	// Is Flat
	types.WriteBoolean(buf, false)

	// Has Death location
	types.WriteBoolean(buf, false)

	// Space to implement Death location here!

	// Portal cooldown
	types.WriteVarInt(buf, 20)

	types.WriteVarInt(buf, 0)

	types.WriteBoolean(buf, false)

	fmt.Printf("Login Data: % 2x\n", buf.Bytes())

	err := types.WriteVarInt(c.conn, int32(buf.Len()))
	if err != nil {
		fmt.Printf("Failed to write send login(play) packet: %v\n", err)
		c.Disconnect("login(play) failed", "red")
		return
	}
	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Failed to write send login(play) packet: %v\n", err)
		c.Disconnect("login(play) failed", "red")
		return
	}

}

func (c *Client) sendPlayerPosition() {
	buf := &bytes.Buffer{}

	// PACKETPLAYERPOSITION
	buf.WriteByte(byte(PACKETPLAYERPOSITION))

	// Teleport ID
	types.WriteVarInt(buf, 0xdead)

	// X Y Z
	types.WriteDouble(buf, 0)
	types.WriteDouble(buf, 1)
	types.WriteDouble(buf, 0)

	// Velocity X Y Z
	types.WriteDouble(buf, 0)
	types.WriteDouble(buf, 0)
	types.WriteDouble(buf, 0)

	// Yaw, Pitch
	types.WriteFloat(buf, 0)
	types.WriteFloat(buf, 0)

	// Teleport Flags
	buf.WriteByte(0) // relative pos
	buf.WriteByte(0) // relative velocity
	buf.WriteByte(0) // Rotate velocity ..
	buf.WriteByte(0) // unused

	err := types.WriteVarInt(c.conn, int32(buf.Len()))
	if err != nil {
		fmt.Printf("Failed to write send teleport packet: %v\n", err)
		c.Disconnect("teleport failed", "red")
		return
	}
	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Failed to write send teleport packet: %v\n", err)
		c.Disconnect("teleport failed", "red")
		return
	}
}

func (c *Client) sendChunkCenter() {
	buf := &bytes.Buffer{}
	buf.WriteByte(byte(PACKETCHUNKCENTER))
	// Chunk X=0 Y=0
	types.WriteVarInt(buf, 0)
	types.WriteVarInt(buf, 0)

	err := types.WriteVarInt(c.conn, int32(buf.Len()))
	if err != nil {
		fmt.Printf("Failed to write send chunk center packet: %v\n", err)
		c.Disconnect("Loading World failed", "red")
		return
	}
	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Failed to write send chunk center packet: %v\n", err)
		c.Disconnect("Loading World failed", "red")
		return
	}
}

func (c *Client) sendChunkData() {
	buf := &bytes.Buffer{}
	buf.WriteByte(byte(PACKETCHUNKDATA))
	// Chunk X=0 Y=0
	types.WriteVarInt(buf, 0)
	types.WriteVarInt(buf, 0)

	// Chunk Data
	// prefixed array hightmap
	types.WriteVarInt(buf, 0)
	// prefixed array data
	types.WriteVarInt(buf, 0)
	// prefixed array block entities
	types.WriteVarInt(buf, 0)

	// Light Data
	// bitset sky light
	types.WriteVarInt(buf, 0)
	// bitset block light
	types.WriteVarInt(buf, 0)
	// bitset empty sky light
	types.WriteVarInt(buf, 0)
	// bitset empty block light
	types.WriteVarInt(buf, 0)
	// bitset sky light arrays
	types.WriteVarInt(buf, 0)
	// bitset sky light arrays
	types.WriteVarInt(buf, 0)

	err := types.WriteVarInt(c.conn, int32(buf.Len()))
	if err != nil {
		fmt.Printf("Failed to write send chunk center packet: %v\n", err)
		c.Disconnect("Loading World failed", "red")
		return
	}
	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Failed to write send chunk center packet: %v\n", err)
		c.Disconnect("Loading World failed", "red")
		return
	}
}
