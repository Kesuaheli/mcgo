package main

import (
	"mcgo/types"

	"github.com/google/uuid"
)

type CommunicationState int32

const (
	STATEHANDSHAKE = iota
	STATESTATUS
	STATELOGIN
	STATECONFIGURATION
	STATEPLAYING
)

func (state CommunicationState) String() string {
	switch state {
	case STATEHANDSHAKE:
		return "HANDSHAKE"
	case STATESTATUS:
		return "STATUS"
	case STATELOGIN:
		return "LOGIN"
	case STATECONFIGURATION:
		return "CONFIGURATION"
	case STATEPLAYING:
		return "PLAYING"
	default:
		return "UNKNOWN_STATE"
	}
}

type Packet int32

const (
	// A serverbound packet in the handshaking phase. Used to switch to the target state
	PACKETHANDSHAKE Packet = 0x00

	// A clientbound packet in the status phase. Response of the server to a PACKETSTATUSREQUEST
	// packet.
	PACKETSTATUSRESPONSE Packet = 0x00
	// A clientbound packet in the status phase. Pong response of the server to a PACKETPINGREQUEST
	// packet.
	PACKETPINGRESPONSE Packet = 0x01
	// A serverbound packet in the status phase. Describes the request from the client to get status
	// information about the server. Should be answered with a PACKETSTATUSRESPONSE.
	PACKETSTATUSREQUEST Packet = 0x00
	// A serverbound packet in the status phase. A ping of the client. Should be answered with a
	// pong (PACKETPINGRESPONSE).
	PACKETPINGREQUEST Packet = 0x01

	// A clientbound packet in the login phase.
	PACKETDISCONNECT Packet = 0x00
	// A clientbound packet in the login phase.
	PACKETENCRYPTIONREQUEST Packet = 0x01
	// A clientbound packet in the login phase.
	PACKETLOGINSUCCESS Packet = 0x02
	// A clientbound packet in the login phase.
	PACKETSETCOMPRESSEION Packet = 0x03
	// A clientbound packet in the login phase.
	PACKETLOGINPLUGINREQUEST Packet = 0x04
	// A clientbound packet in the login phase.
	PACKETCOOKIEREQUEST Packet = 0x05

	// A serverbound packet in the login phase.
	PACKETLOGINSTART Packet = 0x00
	// A serverbound packet in the login phase.
	PACKETENCRYPTIONRESPONSE Packet = 0x01
	// A serverbound packet in the login phase.
	PACKETLOGINPLUGINRESPONSE Packet = 0x02
	// A serverbound packet in the login phase.
	PACKETLOGINACKNOWLEDGED Packet = 0x03
	// A serverbound packet in the login phase.
	PACKETCOOKIERESPONSE Packet = 0x04

	// A clientbound packet in the configuration phase.
	PACKETFINISHCONFIGURATION Packet = 0x03

	// A serverbound packet in the configuration phase.
	PACKETCLIENTINFORMATION Packet = 0x00
	// A serverbound packet in the configuration phase.
	PACKETPLUGINMESSAGE Packet = 0x02

	PACKETSELECTKNOWNPACKSSERVER Packet = 0x07

	PACKETTRANSFERCONFIGURATION Packet = 0x0B
	PACKETTRANSFERPLAYING       Packet = 0x7A

	PACKETREGISTRYDATA Packet = 0x07

	// Offset in the packet registry: 1.21.8 // 1.21.4

	PACKETBLOCKUPDATE Packet = 0x08 // 0x23
	PACKETGAMEEVENT   Packet = 0x22 // 0x23

	PACKETLOGINPLAYING   Packet = 0x2B // 0x2C
	PACKETPLAYERPOSITION Packet = 0x41 // 0x42
	PACKETACCEPTTELEPORT Packet = 0x00

	PACKETCHUNKCENTER Packet = 0x57 // 0x58
	PACKETKEEPALIVE   Packet = 0x26 // 0x27
	PACKETCHUNKDATA   Packet = 0x27 // 0x28

	PACKETCLIENTTICKEND    Packet = 0x0c // 0x0b
	PACKETMOVEPLAYERPOS    Packet = 0x1d // 0x1d
	PACKETMOVEPLAYERPOSROT Packet = 0x1e // 0x1d
)

type StatusResponse struct {
	Version            StatusResponseVersion `json:"version"`
	Players            StatusResponsePlayers `json:"players,omitempty"`
	Description        types.TextComponent   `json:"description,omitempty"`
	Favicon            string                `json:"favicon,omitempty"`
	EnforcesSecureChat bool                  `json:"enforcesSecureChat"`
	PreviewsChat       bool                  `json:"previewsChat"`
}

type StatusResponseVersion struct {
	Name     string `json:"name"`
	Protocol int32  `json:"protocol"`
}

type StatusResponsePlayers struct {
	Max    int                    `json:"max"`
	Online int                    `json:"online"`
	Sample []StatusResponsePlayer `json:"sample,omitempty"`
}

type StatusResponsePlayer struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}
