package outbound

import (
	"bufio"
	"encoding/binary"

	"marmalade/helpers"
)

func writeByte(b byte) helpers.Action {
	return func(writer *bufio.Writer) error {
		return writer.WriteByte(b)
	}
}

func writeShort(s uint16) helpers.Action {
	return func(writer *bufio.Writer) error {
		buf := [2]byte{}
		binary.BigEndian.PutUint16(buf[:], s)
		_, err := writer.Write(buf[:])
		return err
	}
}

func writeString(s string) helpers.Action {
	return func(writer *bufio.Writer) error {
		_, err := writer.WriteString(classicString(s))
		return err
	}
}

func writeBytes(b []byte) helpers.Action {
	return func(writer *bufio.Writer) error {
		_, err := writer.Write(b)
		return err
	}
}
