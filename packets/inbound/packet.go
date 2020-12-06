package inbound

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"marmalade/packets"
)

type action = func(*bufio.Reader) error

func do(reader *bufio.Reader, actions ...action) error {
	for _, v := range actions {
		if err := v(reader); err != nil {
			return err
		}
	}
	return nil
}

var PacketIDAssertionFailure = errors.New("failed to assert packet id")

func assertPacketID(id byte) action {
	return func(reader *bufio.Reader) error {
		read, readErr := reader.ReadByte()
		if readErr != nil {
			return readErr
		}
		if read != id {
			return fmt.Errorf("%w, expected packet id %v, but instead got %v", PacketIDAssertionFailure, id, read)
		}
		return nil
	}
}

func readByte(b *byte) action {
	return func(reader *bufio.Reader) error {
		read, readErr := reader.ReadByte()
		if readErr != nil {
			return readErr
		}
		*b = read
		return nil
	}
}

func readShort(s *uint16) action {
	return func(reader *bufio.Reader) error {
		read, readErr := packets.ReadN(reader, 2)
		if readErr != nil {
			return readErr
		}
		*s = binary.BigEndian.Uint16(read)
		return nil
	}
}

func readString(s *string) action {
	return func(reader *bufio.Reader) error {
		read, readErr := packets.ReadN(reader, 64)
		if readErr != nil {
			return readErr
		}
		*s = strings.TrimSpace(string(read))
		return nil
	}
}

func discard(n int) action {
	return func(reader *bufio.Reader) error {
		if _, err := reader.Discard(n); err != nil {
			return err
		}
		return nil
	}
}
