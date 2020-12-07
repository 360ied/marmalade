package world

import (
	"sync"

	"marmalade/config"
	"marmalade/packets/outbound"
)

type (
	Player struct {
		Username string
		Position

		ID uint8
		OP bool

		Writer *outbound.AFCBW
	}

	Position struct {
		X, Y, Z    uint16
		Yaw, Pitch uint8
	}
)

var (
	Players   = [256]*Player{}
	PlayersMu = new(sync.Mutex)

	World = NewConcurrentSlice(config.WorldXSize * config.WorldYSize * config.WorldZSize)
)

// returns true if there is space to put another player
func AddPlayer(player *Player) bool {
	PlayersMu.Lock()
	defer PlayersMu.Unlock()
	for i, v := range Players {
		if v == nil {
			Players[i] = player
			player.ID = uint8(i)
			return true
		}
	}
	return false
}

func RemovePlayer(id uint8) {
	PlayersMu.Lock()
	defer PlayersMu.Unlock()

	for _, v := range Players {
		if v != nil {
			go func() {
				_ = v.Writer.SendDespawnPlayer(id)
			}()
		}
	}

	Players[id] = nil
}
