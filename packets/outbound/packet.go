package outbound

import (
	"bufio"
	"encoding/binary"
)

func writeByte(b byte) action {
	return func(writer *bufio.Writer) error {
		return writer.WriteByte(b)
	}
}

func writeShort(s uint16) action {
	return func(writer *bufio.Writer) error {
		buf := [2]byte{}
		binary.BigEndian.PutUint16(buf[:], s)
		_, err := writer.Write(buf[:])
		return err
	}
}

func writeString(s string) action {
	return func(writer *bufio.Writer) error {
		_, err := writer.WriteString(classicString(s))
		return err
	}
}

func writeBytes(b []byte) action {
	return func(writer *bufio.Writer) error {
		_, err := writer.Write(b)
		return err
	}
}
