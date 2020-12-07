package outbound

func (w *AFCBW) SendDisconnectPlayer(reason string) error {
	return w.do(writeByte(0x0e),
		writeString(reason))
}
