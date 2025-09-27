package types

type GameEvent uint8

const (
	GameEventNoRespawnBlockAvailable        GameEvent = 0 // Note: Displays message 'block.minecraft.spawn.not_valid' (You have no home bed or charged respawn anchor, or it was obstructed) to the player.
	GameEventBeginRaining                   GameEvent = 1
	GameEventEndRaining                     GameEvent = 2
	GameEventChangeGameMode                 GameEvent = 3 // 0: Survival, 1: Creative, 2: Adventure, 3: Spectator.
	GameEventWinGame                        GameEvent = 4 // 0: Just respawn player. 1: Roll the credits and respawn player. Note that 1 is only sent by vanilla server when player has not yet achieved advancement "The end?", else 0 is sent.
	GameEventDemoEvent                      GameEvent = 5 // 0: Show welcome to demo screen. 101: Tell movement controls. 102: Tell jump control. 103: Tell inventory control. 104: Tell that the demo is over and print a message about how to take a screenshot.
	GameEventArrowHitPlayer                 GameEvent = 6 // Note: Sent when any player is struck by an arrow.
	GameEventRainLevelChange                GameEvent = 7 // Note: Seems to change both sky color and lighting. Rain level ranging from 0 to 1.
	GameEventThunderLevelChange             GameEvent = 8 // Note: Seems to change both sky color and lighting (same as Rain level change, but doesn't start rain). It also requires rain to render by vanilla client. Thunder level ranging from 0 to 1.
	GameEventPlayPufferfishStingSound       GameEvent = 9
	GameEventPlayElderGuardianMobAppearance GameEvent = 10
	GameEventEnableRespawnScreen            GameEvent = 11 // 0: Enable respawn screen. 1: Immediately respawn (sent when the doImmediateRespawn gamerule changes).
	GameEventLimitedCrafting                GameEvent = 12 // 0: Disable limited crafting. 1: Enable limited crafting (sent when the doLimitedCrafting gamerule changes).
	GameEventStartWaitingForLevelChunks     GameEvent = 13 // Instructs the client to begin the waiting process for the level chunks. Sent by the server after the level is cleared on the client and is being re-sent (either during the first, or subsequent reconfigurations).
)
