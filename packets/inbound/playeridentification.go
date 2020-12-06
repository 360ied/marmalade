package inbound

import "bufio"

func ReadPlayerIdentification(reader *bufio.Reader) (protocolVersion uint8, username, verificationKey string, err error) {
	err = do(reader,
		assertPacketID(0x00),
		readByte(&protocolVersion),
		readString(&username),
		readString(&verificationKey),
		discard(1))
	return
}
