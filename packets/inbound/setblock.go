package inbound

import "bufio"

func ReadSetBlock(reader *bufio.Reader) (x, y, z uint16, mode, blockType byte, err error) {
	err = do(reader,
		assertPacketID(0x05),
		readShort(&x),
		readShort(&y),
		readShort(&z),
		readByte(&mode),
		readByte(&blockType))
	return
}
