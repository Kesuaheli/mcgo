package main

import (
	"encoding/binary"
	"fmt"
	"mcgo/types"
	"os"

	"github.com/google/uuid"
)

func (c *Client) SendPlayerAddTabUpdate() bool {
	/*
		Packet ID: 0x3F
		State: Play
		Bound To: Client
		Resource: player_info_update

		+-----------+---------+-----------+----------------+---------------------------------------------------------------+
		| Field     | Name    | Type      | Notes                                                                 |
		+-----------+---------+-----------+---------------------------------------------------------------+
		| Actions   |         | EnumSet   | Determines what actions are present.                                 |
		| Players   | UUID    | UUID      | The player UUID.                                                     |
		|           | Player  | Prefixed  | Array of Player Actions.                                             |
		|           | Actions | Array     | The length of this array is determined by the number of Player       |
		|           |         |           | Actions that give a non-zero value when applying its mask to the     |
		|           |         |           | actions flag. For example: given decimal 5 (binary 00000101).        |
		|           |         |           | The masks 0x01 and 0x04 return non-zero, meaning the Player Actions  |
		|           |         |           | array would include two actions: Add Player and Update Game Mode.    |
		+-----------+---------+-----------+---------------------------------------------------------------+
	*/

	sender := NewSender(PACKETPLAYERINFOUPDATE)

	sender.Write([]byte{0x9}) // Action Add Player and Update Listed

	type Player struct {
		uuid    uuid.UUID
		name    string
		texture string
	}

	players := map[string]Player{
		"2b33fd89-56ca-4f77-beeb-44b95198014d": {name: "20Philipp13", uuid: uuid.MustParse("2b33fd89-56ca-4f77-beeb-44b95198014d"), texture: "ewogICJ0aW1lc3RhbXAiIDogMTcyMjA3NzM2Njk5OCwKICAicHJvZmlsZUlkIiA6ICIwNDJmNDdkNTlhM2M0Yzk4OWE1MGM3MWYzOGYzOGVkMCIsCiAgInByb2ZpbGVOYW1lIiA6ICJLZXN1YWhlbGkiLAogICJ0ZXh0dXJlcyIgOiB7CiAgICAiU0tJTiIgOiB7CiAgICAgICJ1cmwiIDogImh0dHA6Ly90ZXh0dXJlcy5taW5lY3JhZnQubmV0L3RleHR1cmUvODJmMzNiNmZlM2JiMzZiODg1NTkwZjM3NzZkYWUxNDRmODMyMzM2ZmM5NmJkNGJjYzcxYzUxYWI5ZjM1YmQyIgogICAgfSwKICAgICJDQVBFIiA6IHsKICAgICAgInVybCIgOiAiaHR0cDovL3RleHR1cmVzLm1pbmVjcmFmdC5uZXQvdGV4dHVyZS9hZmQ1NTNiMzkzNThhMjRlZGZlM2I4YTlhOTM5ZmE1ZmE0ZmFhNGQ5YTljM2Q2YWY4ZWFmYjM3N2ZhMDVjMmJiIgogICAgfQogIH0KfQ=="},
		"042f47d5-9a3c-4c98-9a50-c71f38f38ed0": {name: "Kesuaheli", uuid: uuid.MustParse("042f47d5-9a3c-4c98-9a50-c71f38f38ed0"), texture: "ewogICJ0aW1lc3RhbXAiIDogMTcyMjA3NzM2Njk5OCwKICAicHJvZmlsZUlkIiA6ICIwNDJmNDdkNTlhM2M0Yzk4OWE1MGM3MWYzOGYzOGVkMCIsCiAgInByb2ZpbGVOYW1lIiA6ICJLZXN1YWhlbGkiLAogICJ0ZXh0dXJlcyIgOiB7CiAgICAiU0tJTiIgOiB7CiAgICAgICJ1cmwiIDogImh0dHA6Ly90ZXh0dXJlcy5taW5lY3JhZnQubmV0L3RleHR1cmUvODJmMzNiNmZlM2JiMzZiODg1NTkwZjM3NzZkYWUxNDRmODMyMzM2ZmM5NmJkNGJjYzcxYzUxYWI5ZjM1YmQyIgogICAgfSwKICAgICJDQVBFIiA6IHsKICAgICAgInVybCIgOiAiaHR0cDovL3RleHR1cmVzLm1pbmVjcmFmdC5uZXQvdGV4dHVyZS9hZmQ1NTNiMzkzNThhMjRlZGZlM2I4YTlhOTM5ZmE1ZmE0ZmFhNGQ5YTljM2Q2YWY4ZWFmYjM3N2ZhMDVjMmJiIgogICAgfQogIH0KfQ=="},
	}

	connectedClients := []*Client{}

	for _, client := range c.server.clients {
		if client.connected {
			fmt.Println("Found connected client")
			connectedClients = append(connectedClients, client)
		} else {
			fmt.Println("Found disconnected client ", client.String())
		}
	}

	fmt.Println("Connected Client amount: ", len(connectedClients))
	fmt.Println("Full Client list amount: ", len(c.server.clients))

	types.WriteVarInt(sender, int32(len(connectedClients)))

	for _, client := range connectedClients {
		player := players[client.uuid.String()]
		binary.Write(sender, binary.BigEndian, player.uuid)
		types.WriteString(sender, player.name)

		types.WriteVarInt(sender, 1)
		for i := 0; i < 1; i++ {
			types.WriteString(sender, "textures")
			types.WriteString(sender, player.texture)
			sender.Write([]byte{0}) // signature
		}
		sender.Write([]byte{1}) // update listed boolean
	}

	os.WriteFile("updatePlayerInfo.dump", fmt.Appendf(nil, "0x% 02x\n", sender.Bytes()), 0600)

	return !sender.SendMultipleClients(c.server.clients, "TabUpdate")
}

func (c *Client) AddFakePlayer(UUID uuid.UUID, name string) {
	sender := NewSender(PACKETPLAYERINFOUPDATE)
	sender.Write([]byte{0x9}) // Action add player and update listed

	types.WriteVarInt(sender, 1) // 1 player

	binary.Write(sender, binary.BigEndian, UUID)

	// Action: add player
	types.WriteString(sender, name)
	types.WriteBoolean(sender, false) // no data/skin

	// Action: update listed
	types.WriteBoolean(sender, false)

	sender.Send(c, "add fake player info")
}
