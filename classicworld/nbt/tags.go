package nbt

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"

	"marmalade/helpers"
	"marmalade/packets"
)

const (
	tagEnd byte = iota
	tagByte
	tagShort
	tagInt
	tagLong
	tagFloat
	tagDouble
	tagByteArray
	tagString
	tagList
	tagCompound
	tagIntArray
	tagLongArray
)

func Read(reader *bufio.Reader) (Compound, string, error) {
	if err := assertTypeID(reader, tagCompound); err != nil {
		return nil, "", err
	}
	name, nameErr := readName(reader)
	if nameErr != nil {
		return nil, name, nameErr
	}
	c := Compound{}
	return c, name, readCompound(reader, c)
}

// readCompound reads the payload of a compound
// Note: This interprets the data values as unsigned
// If you want signed, just cast them to their signed equivalents
func readCompound(reader *bufio.Reader, out Compound) error {
	for {
		b, bErr := reader.ReadByte()
		if bErr != nil {
			return bErr
		}
		if b == tagEnd { // `tagEnd`s don't have names
			return nil
		}
		fmt.Println("b", b)
		name, nameErr := readName(reader)
		if nameErr != nil {
			return nameErr
		}
		fmt.Println("name", name)
		switch b {
		case tagByte:
			read, err := reader.ReadByte()
			if err != nil {
				return err
			}
			out[name] = read
		case tagShort:
			bs := [unsafe.Sizeof(uint16(0))]byte{}
			if _, err := reader.Read(bs[:]); err != nil {
				return err
			}
			out[name] = binary.BigEndian.Uint16(bs[:])
		case tagInt:
			bs := [unsafe.Sizeof(uint32(0))]byte{}
			if _, err := reader.Read(bs[:]); err != nil {
				return err
			}
			out[name] = binary.BigEndian.Uint32(bs[:])
		case tagLong:
			bs := [unsafe.Sizeof(uint64(0))]byte{}
			if _, err := reader.Read(bs[:]); err != nil {
				return err
			}
			out[name] = binary.BigEndian.Uint64(bs[:])
		case tagFloat:
			bs := [unsafe.Sizeof(float32(0))]byte{}
			if _, err := reader.Read(bs[:]); err != nil {
				return err
			}
			out[name] = math.Float32frombits(binary.BigEndian.Uint32(bs[:]))
		case tagDouble:
			bs := [unsafe.Sizeof(float64(0))]byte{}
			if _, err := reader.Read(bs[:]); err != nil {
				return err
			}
			out[name] = math.Float64frombits(binary.BigEndian.Uint64(bs[:]))
		case tagByteArray:
			prefixBS := [unsafe.Sizeof(int32(0))]byte{}
			if _, err := reader.Read(prefixBS[:]); err != nil {
				return err
			}
			length := int32(binary.BigEndian.Uint32(prefixBS[:]))
			if length < 0 { // cannot have a byte array of negative length
				return newInvalidLengthError(length)
			}
			// note: consider expending a bytes.Buffer here instead of allocating a new slice
			read, readErr := packets.ReadN(reader, int(length)) // use make instead of a bytes.Buffer as the data will be used after the function is done
			if readErr != nil {
				return readErr
			}
			out[name] = read
		case tagString:
			prefixBS := [unsafe.Sizeof(uint16(0))]byte{}
			if _, err := reader.Read(prefixBS[:]); err != nil {
				return err
			}
			length := int32(binary.BigEndian.Uint16(prefixBS[:]))
			if length < 0 {
				return newInvalidLengthError(length)
			}
			if err := func() error { // use a closure here so that the defer is run once we no longer need the buffer
				buf := helpers.GetBuffer()
				defer helpers.PutBuffer(buf)
				buf.Reset()
				buf.Write(helpers.GetZeroes(int(length))) // grow doesn't change the length of buf.Bytes()
				if _, err := reader.Read(buf.Bytes()); err != nil {
					return err
				}
				out[name] = buf.String() // this allocates a new string
				return nil
			}(); err != nil {
				return err
			}
		case tagList:
			return NotImplementedError // too problematic to implement, and the classic world format doesn't even use them anyways
		case tagCompound:
			c := Compound{}
			if err := readCompound(reader, c); err != nil {
				return err
			}
			out[name] = c
		case tagIntArray, tagLongArray:
			return NotImplementedError // classic world format doesn't use these
		default:
			return newInvalidTypeIDError(b)
		}
	}
}

// reads a tag name
func readName(reader *bufio.Reader) (string, error) {
	bs := [unsafe.Sizeof(uint16(0))]byte{}
	if _, err := reader.Read(bs[:]); err != nil {
		return "", err
	}
	length := int32(binary.BigEndian.Uint16(bs[:]))
	if length < 0 {
		return "", newInvalidLengthError(length)
	}

	buf := helpers.GetBuffer()
	defer helpers.PutBuffer(buf)
	buf.Reset()
	buf.Write(helpers.GetZeroes(int(length))) // grow doesn't change the length of buf.Bytes()

	if _, err := reader.Read(buf.Bytes()); err != nil {
		return "", err
	}

	return buf.String(), nil
}
