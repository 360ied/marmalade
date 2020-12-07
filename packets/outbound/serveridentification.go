package outbound

func (w *AFCBW) SendServerIdentification(serverName, serverMOTD string, isOP bool) error {
	return w.do(writeByte(0x00), // packet id
		writeByte(0x07), // protocol version
		writeString(serverName),
		writeString(serverMOTD),
		writeByte(opByte(isOP)))
}
