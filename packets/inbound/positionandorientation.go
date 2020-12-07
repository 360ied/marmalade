package inbound

import "bufio"

func ReadPositionAndOrientation(reader *bufio.Reader) (x, y, z uint16, yaw, pitch uint8, err error) {
	err = do(reader,
		assertPacketID(0x08),
		discard(1), // player id, always 255, referring to itself
		readShort(&x),
		readShort(&y),
		readShort(&z),
		readByte(&yaw),
		readByte(&pitch))
	return
}
