package helpers

import (
	"bytes"
	"sync"
)

var (
	bufferPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
	zeroes     []byte
	zeroesMu   sync.Mutex
)

// Fast, thread-safe way to get a slice of zero bytes
// DO NOT MODIFY WHAT IS RETURNED OR IT WILL CONTAMINATE THE BUFFER
func GetZeroes(n int) []byte {
	zeroesMu.Lock()
	defer zeroesMu.Unlock()

	// If zeroes is not large enough, allocate a larger slice
	if len(zeroes) < n {
		zeroes = make([]byte, MaxInt(n, len(zeroes)*2))
	}

	return zeroes[:n]
}

// Type-safe buffer get
func GetBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// Type-safe buffer put
func PutBuffer(buf *bytes.Buffer) {
	bufferPool.Put(buf)
}

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
