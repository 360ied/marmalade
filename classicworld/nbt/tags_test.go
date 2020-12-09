package nbt

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"testing"
)

// fuzzing
// warning: uses a lot of ram
func TestReadCompound(t *testing.T) {
	source := bufio.NewReader(rand.Reader)
	for {
		m := Compound{}
		if err := readCompound(source, m); err != nil {
			log.Println(err)
		}
	}
}

func TestReadWorld(t *testing.T) {
	file, fileErr := os.Open("tests/nbttest.cw")
	if fileErr != nil {
		panic(fileErr)
	}
	gzipReader, gzipReaderErr := gzip.NewReader(file)
	if gzipReaderErr != nil {
		panic(gzipReaderErr)
	}
	fmt.Println(gzipReader.Header)
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(gzipReader); err != nil {
		panic(err)
	}
	bufReader := bufio.NewReader(buf)

	// fmt.Println(buf.Bytes()[:1000])

	c, name, err := Read(bufReader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
	fmt.Println(name)
}
