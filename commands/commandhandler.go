package commands

import (
	"fmt"
	"strconv"
	"strings"

	"marmalade/world"
)

var commandMap = map[string]func(*world.Player, []string){
	"ping": ping,
	"tp":   teleport,
	"fill": fill,
}

func HandleCommand(player *world.Player, command string) {
	split := strings.Split(command, " ")
	fun, found := commandMap[split[0]]
	if !found {
		_ = player.Writer.SendMessageStr(fmt.Sprintf("[System] Unknown command `%v`", split[0]))
		return
	}
	fun(player, split[1:])
}

func ping(player *world.Player, _ []string) {
	_ = player.Writer.SendMessageStr("[System] Pong!")
}

func teleport(player *world.Player, args []string) {
	// if !player.OP {
	// 	_ = player.Writer.SendMessageStr("[System] You do not have the permissions to run this command.")
	// 	return
	// }
	switch len(args) {
	case 1:
		targetUsername := args[0]

		world.PlayersMu.Lock()
		defer world.PlayersMu.Unlock()

		// find target player
		for _, v := range world.Players {
			if v != nil && strings.EqualFold(targetUsername, v.Username) {
				// update player's position
				// the player will broadcast this new position themselves, so no need to broadcast it for them
				_ = player.Writer.SendPositionAndOrientation(255, player.X, player.Y, player.Z, player.Yaw, player.Pitch)
				_ = player.Writer.SendMessageStr("Done.")
				return
			}
		}
		_ = player.Writer.SendMessageStr("[System] Player not found!")
	default:
		_ = player.Writer.SendMessageStr("[System] Invalid number of arguments.")
	}
}

func fill(player *world.Player, args []string) {
	if len(args) != 7 {
		_ = player.Writer.SendMessageStr("[System] Too many or too few arguments! 7 required.")
	}

	// Decode arguments
	block, blockErr := strconv.Atoi(args[0])
	if blockErr != nil {
		_ = world.SendLargeMessage(player, fmt.Sprintf("[System] Failed to decode parameter block: %v", blockErr))
		return
	}
	lesserX, lesserXErr := strconv.Atoi(args[1])
	if lesserXErr != nil {
		_ = world.SendLargeMessage(player, fmt.Sprintf("[System] Failed to decode parameter lesserX: %v", lesserXErr))
		return
	}
	lesserY, lesserYErr := strconv.Atoi(args[2])
	if lesserYErr != nil {
		_ = world.SendLargeMessage(player, fmt.Sprintf("[System] Failed to decode parameter lesserY: %v", lesserYErr))
		return
	}
	lesserZ, lesserZErr := strconv.Atoi(args[3])
	if lesserZErr != nil {
		_ = world.SendLargeMessage(player, fmt.Sprintf("[System] Failed to decode parameter lesserZ: %v", lesserZErr))
		return
	}
	greaterX, greaterXErr := strconv.Atoi(args[4])
	if greaterXErr != nil {
		_ = world.SendLargeMessage(player, fmt.Sprintf("[System] Failed to decode parameter greaterX: %v", greaterXErr))
		return
	}
	greaterY, greaterYErr := strconv.Atoi(args[5])
	if greaterYErr != nil {
		_ = world.SendLargeMessage(player, fmt.Sprintf("[System] Failed to decode parameter greaterY: %v", greaterYErr))
		return
	}
	greaterZ, greaterZErr := strconv.Atoi(args[6])
	if greaterZErr != nil {
		_ = world.SendLargeMessage(player, fmt.Sprintf("[System] Failed to decode parameter greaterZ: %v", greaterZErr))
		return
	}

	for x := lesserX; x <= greaterX; x++ {
		for y := lesserY; y <= greaterY; y++ {
			for z := lesserZ; z <= greaterZ; z++ {
				world.HandleSetBlock(uint16(x), uint16(y), uint16(z), 1, byte(block))
			}
		}
	}

	_ = player.Writer.SendMessageStr("Done.")
}
