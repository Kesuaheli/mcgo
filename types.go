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

	PACKETADDENTITY              Packet = 0x01
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

	PACKETSETPLAYERINVENTORY Packet = 0x65 // 0x64

	// A clientbound packet to send a message to the player.
	PACKETPLAYERCHAT       Packet = 0x3A
	PACKETPLAYERINFOUPDATE Packet = 0x3F
	PACKETSYSTEMCHAT       Packet = 0x72

	PACKETCLIENTTICKEND    Packet = 0x0c // 0x0b
	PACKETMOVEPLAYERPOS    Packet = 0x1d // 0x1d
	PACKETMOVEPLAYERPOSROT Packet = 0x1e // 0x1d

	// A serverbound packet when a player sends a chat message.
	PACKETPLAYERSENTMESSAGE Packet = 0x08 // 0x09
	// A serverbound packet to respond to the [PACKETKEEPALIVE] packet.
	PACKETCLIENTKEEPALIVE Packet = 0x1B
	PACKETPLAYERROTATED   Packet = 0x1F
	PACKETPLAYERINPUT     Packet = 0x2A
	PACKETPLAYERLOADED    Packet = 0x2B
	// A serverbound packet sent when the client clicks a text component with the minecraft:custom click action.
	PACKETCUSTOMCLICKACTION Packet = 0x41
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

type EntityType int32

const (
	TYPE_acacia_boat EntityType = iota
	TYPE_acacia_chest_boat
	TYPE_allay
	TYPE_area_effect_cloud
	TYPE_armadillo
	TYPE_armor_stand
	TYPE_arrow
	TYPE_axolotl
	TYPE_bamboo_chest_raft
	TYPE_bamboo_raft
	TYPE_bat
	TYPE_bee
	TYPE_birch_boat
	TYPE_birch_chest_boat
	TYPE_blaze
	TYPE_block_display
	TYPE_bogged
	TYPE_breeze
	TYPE_breeze_wind_charge
	TYPE_camel
	TYPE_cat
	TYPE_cave_spider
	TYPE_cherry_boat
	TYPE_cherry_chest_boat
	TYPE_chest_minecart
	TYPE_chicken
	TYPE_cod
	TYPE_command_block_minecart
	TYPE_cow
	TYPE_creaking
	TYPE_creeper
	TYPE_dark_oak_boat
	TYPE_dark_oak_chest_boat
	TYPE_dolphin
	TYPE_donkey
	TYPE_dragon_fireball
	TYPE_drowned
	TYPE_egg
	TYPE_elder_guardian
	TYPE_enderman
	TYPE_endermite
	TYPE_ender_dragon
	TYPE_ender_pearl
	TYPE_end_crystal
	TYPE_evoker
	TYPE_evoker_fangs
	TYPE_experience_bottle
	TYPE_experience_orb
	TYPE_eye_of_ender
	TYPE_falling_block
	TYPE_fireball
	TYPE_firework_rocket
	TYPE_fox
	TYPE_frog
	TYPE_furnace_minecart
	TYPE_ghast
	TYPE_happy_ghast
	TYPE_giant
	TYPE_glow_item_frame
	TYPE_glow_squid
	TYPE_goat
	TYPE_guardian
	TYPE_hoglin
	TYPE_hopper_minecart
	TYPE_horse
	TYPE_husk
	TYPE_illusioner
	TYPE_interaction
	TYPE_iron_golem
	TYPE_item
	TYPE_item_display
	TYPE_item_frame
	TYPE_jungle_boat
	TYPE_jungle_chest_boat
	TYPE_leash_knot
	TYPE_lightning_bolt
	TYPE_llama
	TYPE_llama_spit
	TYPE_magma_cube
	TYPE_mangrove_boat
	TYPE_mangrove_chest_boat
	TYPE_marker
	TYPE_minecart
	TYPE_mooshroom
	TYPE_mule
	TYPE_oak_boat
	TYPE_oak_chest_boat
	TYPE_ocelot
	TYPE_ominous_item_spawner
	TYPE_painting
	TYPE_pale_oak_boat
	TYPE_pale_oak_chest_boat
	TYPE_panda
	TYPE_parrot
	TYPE_phantom
	TYPE_pig
	TYPE_piglin
	TYPE_piglin_brute
	TYPE_pillager
	TYPE_polar_bear
	TYPE_splash_potion
	TYPE_lingering_potion
	TYPE_pufferfish
	TYPE_rabbit
	TYPE_ravager
	TYPE_salmon
	TYPE_sheep
	TYPE_shulker
	TYPE_shulker_bullet
	TYPE_silverfish
	TYPE_skeleton
	TYPE_skeleton_horse
	TYPE_slime
	TYPE_small_fireball
	TYPE_sniffer
	TYPE_snowball
	TYPE_snow_golem
	TYPE_spawner_minecart
	TYPE_spectral_arrow
	TYPE_spider
	TYPE_spruce_boat
	TYPE_spruce_chest_boat
	TYPE_squid
	TYPE_stray
	TYPE_strider
	TYPE_tadpole
	TYPE_text_display
	TYPE_tnt
	TYPE_tnt_minecart
	TYPE_trader_llama
	TYPE_trident
	TYPE_tropical_fish
	TYPE_turtle
	TYPE_vex
	TYPE_villager
	TYPE_vindicator
	TYPE_wandering_trader
	TYPE_warden
	TYPE_wind_charge
	TYPE_witch
	TYPE_wither
	TYPE_wither_skeleton
	TYPE_wither_skull
	TYPE_wolf
	TYPE_zoglin
	TYPE_zombie
	TYPE_zombie_horse
	TYPE_zombie_villager
	TYPE_zombified_piglin
	TYPE_player
	TYPE_fishing_bobber
)
