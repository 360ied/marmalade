package outbound

func (w *AFCBW) SendSetBlock(x, y, z uint16, blockType byte) error {
	return w.do(writeByte(0x06),
		writeShort(x),
		writeShort(y),
		writeShort(z),
		writeByte(blockType))
}
