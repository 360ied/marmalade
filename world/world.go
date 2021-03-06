package world

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"marmalade/classicworld/nbt"
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

	Blocks *ConcurrentSlice

	// Name string
	// UUID []byte

	XSize uint16
	YSize uint16
	ZSize uint16

	SpawnPos Position
)

func Initialize() {
	worldFile, worldFileErr := os.Open(config.WorldPath)
	if worldFileErr != nil {
		panic(worldFileErr)
	}
	defer func() { _ = worldFile.Close() }()
	// gzipR, gzipRErr := gzip.NewReader(worldFile)
	// if gzipRErr != nil {
	// 	panic(gzipRErr)
	// }
	// defer func() { _ = gzipR.Close() }()
	bufR := bufio.NewReader(worldFile)

	wNBT, _ /* wNBTName */, wNBTErr := nbt.Read(bufR)
	if wNBTErr != nil {
		panic(wNBTErr)
	}

	// Name = wNBT["Name"].(string)
	// UUID = wNBT["UUID"].([]byte)

	XSize = wNBT["X"].(uint16)
	YSize = wNBT["Y"].(uint16)
	ZSize = wNBT["Z"].(uint16)

	spawnNBT := wNBT["Spawn"].(nbt.Compound)
	SpawnPos.X = spawnNBT["X"].(uint16)
	SpawnPos.Y = spawnNBT["Y"].(uint16)
	SpawnPos.Z = spawnNBT["Z"].(uint16)
	SpawnPos.Yaw = spawnNBT["H"].(uint8) // Heading is another name for yaw
	SpawnPos.Pitch = spawnNBT["P"].(uint8)

	Blocks = NewConcurrentSlice(wNBT["BlockArray"].([]byte))

	log.Printf("[INFO] Loaded map %v", config.WorldPath)

	go func() {
		for {
			time.Sleep(config.WorldSaveDelay)
			if err := save(); err != nil {
				log.Printf("[ERROR] Failed to save world: %v", err)
			} else {
				log.Print("[INFO] Successfully saved world.")
			}
		}
	}()
}

func save() error {
	// Write-only, because we're not going to read or append anything
	// Create, because we want to create a new file if it didn't previously exist
	worldScFile, worldScFileErr := os.OpenFile(config.WorldScratchPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if worldScFileErr != nil {
		return worldScFileErr
	}
	defer func() { _ = worldScFile.Close() }()
	// buf0 := bufio.NewWriter(worldScFile)
	// gzipW := gzip.NewWriter(buf0) // compression unexpected eof for some reason
	bufW := bufio.NewWriter(worldScFile)

	snapshot := snapshotBufferPool.Get().([]byte)
	Blocks.Snapshot(snapshot)

	if err := nbt.DoWrite(bufW,
		nbt.WriteCompound("ClassicWorld"),
		nbt.WriteShort("X", XSize),
		nbt.WriteShort("Y", YSize),
		nbt.WriteShort("Z", ZSize),
		nbt.WriteCompound("Spawn"),
		nbt.WriteShort("X", SpawnPos.X),
		nbt.WriteShort("Y", SpawnPos.Y),
		nbt.WriteShort("Z", SpawnPos.Z),
		nbt.WriteByte("H", SpawnPos.Yaw),
		nbt.WriteByte("P", SpawnPos.Pitch),
		nbt.WriteEnd(), // Spawn
		nbt.WriteByteArray("BlockArray", snapshot),
		nbt.WriteCompound("Metadata"),
		nbt.WriteString("Made_With", "marmalade"),
		nbt.WriteEnd(), // Metadata
		nbt.WriteEnd(), // ClassicWorld
	); err != nil {
		return err
	}

	// Close writers before renaming the files
	// if err := buf0.Flush(); err != nil {
	// 	return err
	// }
	// if err := gzipW.Close(); err != nil {
	// 	return err
	// }
	if err := bufW.Flush(); err != nil {
		return err
	}
	if err := worldScFile.Close(); err != nil {
		return err
	}

	// First rename WorldPath to WorldTempPath instead of deletion in case renaming WorldScratchPath to WorldPath fails
	if err := os.Rename(config.WorldPath, config.WorldTempPath); err != nil {
		return err
	}
	if err := os.Rename(config.WorldScratchPath, config.WorldPath); err != nil {
		return err
	}

	// Now we can safely delete the old world
	return os.Remove(config.WorldTempPath)
}

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

	return w.SendLevelFinalize(XSize, YSize, ZSize)
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
	return int(y)*int(XSize)*int(ZSize) + int(z)*int(XSize) + int(x)
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

	buf := helpers.GetBuffer()
	defer helpers.PutBuffer(buf)
	buf.Reset()

	for _, v := range lines {
		split := strings.Split(v, " ")
		for _, vv := range split {
			for _, vvv := range helpers.PartitionString(vv, 63) { // 64 - len(' ') = 63
				if buf.Len()+len(vvv) > 64 {
					if err := msgBufFlush(player, buf); err != nil {
						return err
					}
					buf.Reset()
				}
				buf.WriteString(vvv)
			}
			if buf.Len() > 63 {
				if err := msgBufFlush(player, buf); err != nil {
					return err
				}
			}
			buf.WriteByte(' ')
		}
	}
	// flush rest
	return player.Writer.SendMessageBytes(buf.Bytes())
}

func msgBufFlush(player *Player, buf *bytes.Buffer) error {
	if err := player.Writer.SendMessageBytes(buf.Bytes()); err != nil {
		return err
	}
	buf.Reset()
	return nil
}
