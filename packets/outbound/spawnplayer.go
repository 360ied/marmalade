package outbound

func (w *AFCBW) SendSpawnPlayer(playerID uint8, playerName string, x, y, z uint16, yaw, pitch uint8) error {
	return w.do(writeByte(0x07),
		writeByte(playerID),
		writeString(playerName),
		writeShort(x),
		writeShort(y),
		writeShort(z),
		writeByte(yaw),
		writeByte(pitch))
}
