package world

type (
	World struct {
		blocks *ConcurrentSlice
	}

	Player struct {
		Username string
		Position
		OP bool

		writer *AutoFlushingConcurrentBufferedWriter
	}

	Position struct {
		X, Y, Z    uint16
		Yaw, Pitch uint8
	}
)
