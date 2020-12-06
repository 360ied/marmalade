package packets

import "io"

func ReadN(reader io.Reader, n int) ([]byte, error) {
	arr := make([]byte, n)
	_, err := io.ReadFull(reader, arr)
	return arr, err
}
