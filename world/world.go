package world

type (
	World struct {
		blocks *ConcurrentSlice
	}

	Player struct {
		Username   string
		X, Y, Z    uint16
		Yaw, Pitch uint8
		OP         bool

		writer *AutoFlushingConcurrentBufferedWriter
	}
)
