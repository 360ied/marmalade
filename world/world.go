package world

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"io"
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
	Players   = [255]*Player{}
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
			v := v // capture loop variable
			go func() {
				_ = v.Writer.SendDespawnPlayer(id)
			}()
		}
	}

	Players[id] = nil
}

func SendWorld(w *outbound.AFCBW) error {
	if err := w.SendLevelInitialize(); err != nil {
		return err
	}

	pipeR, pipeW := io.Pipe()
	defer func() { _ = pipeR.Close() }()
	defer func() { _ = pipeW.Close() }()

	bufW := bufio.NewWriter(pipeW)
	gzipW := gzip.NewWriter(bufW)

	go func() {
		snapshot := World.Snapshot()
		_ = binary.Write(gzipW, binary.BigEndian, uint32(len(snapshot)))
		_, _ = gzipW.Write(snapshot)
		defer func() { _ = gzipW.Close() }()
		defer func() { _ = bufW.Flush() }()
	}()

	readBuf := make([]byte, 1024)
	for {
		n, err := pipeR.Read(readBuf)
		if err != nil {
			return err
		}
		if sErr := w.SendLevelDataChunk(uint16(n), readBuf, 50); sErr != nil {
			return sErr
		}

		if n < len(readBuf) {
			break
		}
	}

	return w.SendLevelFinalize(uint16(config.WorldXSize), uint16(config.WorldYSize), uint16(config.WorldZSize))
}
