package outbound

func (w *AFCBW) SendMessage(message string) error {
	return w.do(writeByte(0x0d),
		writeByte(0x00), // unused byte? I'm not exactly sure what this is for, though it might represent player ID or message type
		writeString(message))
}
