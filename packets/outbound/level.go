package outbound

func (w *AFCBW) SendLevelInitialize() error {
	return w.do(writeByte(0x02))
}

// destroys parameter data
func (w *AFCBW) SendLevelDataChunk(length uint16, data []byte, percentComplete uint8) error {
	return w.do(writeByte(0x03),
		writeShort(length),
		writeBytes(classicByteArray(data)),
		writeByte(percentComplete))
}

func (w *AFCBW) SendLevelFinalize(xSize, ySize, zSize uint16) error {
	return w.do(writeByte(0x04),
		writeShort(xSize),
		writeShort(ySize),
		writeShort(zSize))
}
