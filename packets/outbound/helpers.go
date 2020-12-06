package outbound

import (
	"strings"

	"marmalade/helpers"
)

// pads a string with spaces and trims it a length of 64
func classicString(s string) string {
	return (s + strings.Repeat("\x20", helpers.MaxInt(64-len(s), 0)))[:64]
}

func opByte(b bool) byte {
	if b {
		return 64
	}
	return 0
}
