package packets

import "io"

// note: fails for n < 0
func ReadN(reader io.Reader, n int) ([]byte, error) {
	arr := make([]byte, n)
	_, err := io.ReadFull(reader, arr)
	return arr, err
}
