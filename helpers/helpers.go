package helpers

import (
	"bytes"
	"sync"
)

var BufferPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// ex: DivIntCeil(10, 3) = 4
func DivIntCeil(x, y int) int {
	div := x / y
	if div*y < x {
		return div + 1
	} else {
		return div
	}
}

func PartitionString(s string, size int) []string {
	buf := make([]string, DivIntCeil(len(s), size))
	for i := range buf {
		buf[i] = s[i*size : MinInt((i+1)*size, len(s))]
	}
	return buf
}
