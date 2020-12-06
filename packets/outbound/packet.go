package outbound

import (
	"bufio"
)

func writeByte(b byte) action {
	return func(writer *bufio.Writer) error {
		return writer.WriteByte(b)
	}
}

func writeString(s string) action {
	return func(writer *bufio.Writer) error {
		_, err := writer.WriteString(classicString(s))
		return err
	}
}
