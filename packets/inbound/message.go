package inbound

import "bufio"

func ReadMessage(reader *bufio.Reader) (message string, err error) {
	err = do(reader,
		assertPacketID(0x0d),
		discard(1), // (player id?), always 255
		readString(&message))
	return
}
