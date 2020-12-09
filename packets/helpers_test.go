package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestReadN(t *testing.T) {
	fmt.Println("1")
	_, _ = ReadN(bytes.NewBufferString("test"), 1)
	fmt.Println("0")
	_, _ = ReadN(bytes.NewBufferString("test"), 0)
	fmt.Println("-1")
	_, _ = ReadN(bytes.NewBufferString("test"), -1)
}
