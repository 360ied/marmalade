package world

import "marmalade/packets/outbound"

type (
	World struct {
		blocks *ConcurrentSlice
	}

	Player struct {
		Username string
		Position
		OP bool

		writer *outbound.AFCBW
	}

	Position struct {
		X, Y, Z    uint16
		Yaw, Pitch uint8
	}
)
