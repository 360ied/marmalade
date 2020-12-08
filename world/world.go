package world

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"strings"
	"sync"

	"marmalade/config"
	"marmalade/helpers"
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

	Blocks = NewConcurrentSlice(config.WorldXSize * config.WorldYSize * config.WorldZSize)
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
			_ = v.Writer.SendDespawnPlayer(id)
		}
	}

	Players[id] = nil
}

var snapshotBufferPool = sync.Pool{New: func() interface{} { return make([]byte, Blocks.Len()) }}

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
		snapshot := snapshotBufferPool.Get().([]byte)
		defer snapshotBufferPool.Put(snapshot)
		Blocks.Snapshot(snapshot)
		_ = binary.Write(gzipW, binary.BigEndian, uint32(len(snapshot)))
		_, _ = gzipW.Write(snapshot)
		_ = gzipW.Close()
		_ = bufW.Flush()
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

func HandleSetBlock(x, y, z uint16, mode, blockType byte) {
	pos := position(x, y, z)
	if pos >= Blocks.Len() || pos < 0 { // bounds check
		return
	}
	if mode == 0x00 {
		blockType = 0x00
	}
	Blocks.Set(pos, blockType)

	PlayersMu.Lock()
	defer PlayersMu.Unlock()

	for _, v := range Players {
		if v != nil {
			_ = v.Writer.SendSetBlock(x, y, z, blockType)
		}
	}
}

// calculates the Blocks index from x, y, z
func position(x, y, z uint16) int {
	return int(y)*config.WorldXSize*config.WorldZSize + int(z)*config.WorldXSize + int(x)
}

func HandlePositionAndOrientation(player *Player, x, y, z uint16, yaw, pitch uint8) {
	PlayersMu.Lock()
	defer PlayersMu.Unlock()

	player.X = x
	player.Y = y
	player.Z = z
	player.Yaw = yaw
	player.Pitch = pitch

	for _, v := range Players {
		if v != nil {
			_ = v.Writer.SendPositionAndOrientation(player.ID, x, y, z, yaw, pitch)
		}
	}
}

func SpawnOtherPlayers(newPlayer *Player) {
	PlayersMu.Lock()
	defer PlayersMu.Unlock()

	// send player other players
	for _, v := range Players {
		if v != nil && v.ID != newPlayer.ID {
			// send player other players
			_ = newPlayer.Writer.SendSpawnPlayer(v.ID, v.Username, v.X, v.Y, v.Z, v.Yaw, v.Pitch)
			// send other players player
			_ = v.Writer.SendSpawnPlayer(newPlayer.ID, newPlayer.Username, newPlayer.X, newPlayer.Y, newPlayer.Z, newPlayer.Yaw, newPlayer.Pitch)
		}
	}
}

func BroadcastMessage(message string) {
	PlayersMu.Lock()
	defer PlayersMu.Unlock()
	for _, v := range Players {
		if v != nil {
			_ = v.Writer.SendMessageStr(message)
		}
	}
}

func SendLargeMessage(player *Player, message string) error {
	lines := strings.Split(message, "\n")

	buf := helpers.BufferPool.Get().(*bytes.Buffer)
	defer helpers.BufferPool.Put(buf)
	buf.Reset()

	for _, v := range lines {
		split := strings.Split(v, " ")
		for _, vv := range split {
			for _, vvv := range helpers.PartitionString(vv, 63) { // 64 - len(' ') = 63
				if buf.Len()+len(vvv) > 64 {
					// flush
					if err := player.Writer.SendMessageBytes(buf.Bytes()); err != nil {
						return err
					}
					buf.Reset()
				}
				buf.WriteString(vvv)
				buf.WriteByte(' ')
			}
		}
	}
	// flush rest
	return player.Writer.SendMessageBytes(buf.Bytes())
}
