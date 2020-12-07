package outbound

func (w *AFCBW) SendPositionAndOrientation(playerID uint8, x, y, z uint16, yaw, pitch uint8) error {
	return w.do(writeByte(0x08),
		writeByte(playerID),
		writeShort(x),
		writeShort(y),
		writeShort(z),
		writeByte(yaw),
		writeByte(pitch))
}
