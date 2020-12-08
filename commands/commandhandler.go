package commands

import (
	"fmt"
	"strings"

	"marmalade/world"
)

var commandMap = map[string]func(*world.Player, []string){
	"ping": ping,
	"tp":   teleport,
}

func HandleCommand(player *world.Player, command string) {
	split := strings.Split(command, " ")
	fun, found := commandMap[split[0]]
	if !found {
		_ = player.Writer.SendMessage(fmt.Sprintf("[System] Unknown command `%v`", split[0]))
		return
	}
	fun(player, split[1:])
}

func ping(player *world.Player, _ []string) {
	_ = player.Writer.SendMessage("[System] Pong!")
}

func teleport(player *world.Player, args []string) {
	if !player.OP {
		_ = player.Writer.SendMessage("[System] You do not have the permissions to run this command.")
		return
	}
	switch len(args) {
	case 1:
		targetUsername := args[0]

		world.PlayersMu.Lock()
		defer world.PlayersMu.Unlock()

		// find target player
		for _, v := range world.Players {
			if v != nil && strings.EqualFold(targetUsername, v.Username) {
				// set player's position to new position
				player.Position = v.Position
				// update player's position
				// the player will broadcast this new position themselves, so no need to broadcast it for them
				_ = player.Writer.SendPositionAndOrientation(255, player.X, player.Y, player.Z, player.Yaw, player.Pitch)
				return
			}
		}
		_ = player.Writer.SendMessage("[System] Player not found!")
	default:
		_ = player.Writer.SendMessage("[System] Invalid number of arguments.")
	}
}
