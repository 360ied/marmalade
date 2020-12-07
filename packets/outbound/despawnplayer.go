package outbound

func (w *AFCBW) SendDespawnPlayer(playerID uint8) error {
	return w.do(writeByte(0x0c),
		writeByte(playerID))
}
