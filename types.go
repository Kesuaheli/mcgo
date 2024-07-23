package main

type CommunicationState int32

const (
	STATEHANDSHAKE = iota
	STATESTATUS
	STATELOGIN
	STATEPLAYING
)

type Packet int32

const (
	// A serverbound packet in the handshaking phase. Used to switch to the target state
	PACKETHANDSHAKE = 0x00

	// A clientbound packet in the status phase. Response of the server to a PACKETSTATUSREQUEST
	// packet.
	PACKETSTATUSRESPONSE = 0x00
	// A clientbound packet in the status phase. Pong response of the server to a PACKETPINGREQUEST
	// packet.
	PACKETPINGRESPONSE = 0x01
	// A serverbound packet in the status phase. Describes the request from the client to get status
	// information about the server. Should be answered with a PACKETSTATUSRESPONSE.
	PACKETSTATUSREQUEST = 0x00
	// A serverbound packet in the status phase. A ping of the client. Should be answered with a
	// pong (PACKETPINGRESPONSE).
	PACKETPINGREQUEST = 0x01
)
