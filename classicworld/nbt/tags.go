package nbt

import (
	"bufio"
	"encoding/binary"
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
		name, nameErr := readName(reader)
		if nameErr != nil {
			return nameErr
		}
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

// Helper function for easier error handling
func DoWrite(writer *bufio.Writer, actions ...helpers.Action) error {
	for _, v := range actions {
		if err := v(writer); err != nil {
			return err
		}
	}
	return nil
}

func WriteEnd() helpers.Action {
	return func(writer *bufio.Writer) error {
		return writer.WriteByte(tagEnd)
	}
}

func WriteByte(name string, b byte) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagByte); err != nil {
			return err
		}
		if err := writeName(writer, name); err != nil {
			return err
		}
		return writer.WriteByte(b)
	}
}

func WriteShort(name string, s uint16) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagShort); err != nil {
			return err
		}
		if err := writeName(writer, name); err != nil {
			return err
		}
		bs := [unsafe.Sizeof(s)]byte{}
		binary.BigEndian.PutUint16(bs[:], s)
		_, err := writer.Write(bs[:])
		return err
	}
}

func WriteInt(name string, i uint32) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagInt); err != nil {
			return err
		}
		if err := writeName(writer, name); err != nil {
			return err
		}
		bs := [unsafe.Sizeof(i)]byte{}
		binary.BigEndian.PutUint32(bs[:], i)
		_, err := writer.Write(bs[:])
		return err
	}
}

func WriteLong(name string, l uint64) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagLong); err != nil {
			return err
		}
		if err := writeName(writer, name); err != nil {
			return err
		}
		bs := [unsafe.Sizeof(l)]byte{}
		binary.BigEndian.PutUint64(bs[:], l)
		_, err := writer.Write(bs[:])
		return err
	}
}

func WriteFloat(name string, f float32) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagFloat); err != nil {
			return err
		}
		if err := writeName(writer, name); err != nil {
			return err
		}
		bs := [unsafe.Sizeof(uint32(0))]byte{}
		binary.BigEndian.PutUint32(bs[:], math.Float32bits(f))
		_, err := writer.Write(bs[:])
		return err
	}
}

func WriteDouble(name string, d float64) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagDouble); err != nil {
			return err
		}
		if err := writeName(writer, name); err != nil {
			return err
		}
		bs := [unsafe.Sizeof(uint64(0))]byte{}
		binary.BigEndian.PutUint64(bs[:], math.Float64bits(d))
		_, err := writer.Write(bs[:])
		return err
	}
}

func WriteByteArray(name string, b []byte) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagByteArray); err != nil {
			return err
		}
		if err := writeName(writer, name); err != nil {
			return err
		}
		if err := binary.Write(writer, binary.BigEndian, uint32(len(b))); err != nil {
			return err
		}
		_, err := writer.Write(b)
		return err
	}
}

func WriteString(name, s string) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagString); err != nil {
			return err
		}
		if err := writeName(writer, name); err != nil {
			return err
		}
		if err := binary.Write(writer, binary.BigEndian, uint16(len(s))); err != nil {
			return err
		}
		_, err := writer.WriteString(s)
		return err
	}
}

// func WriteList

func WriteCompound(name string) helpers.Action {
	return func(writer *bufio.Writer) error {
		if err := writer.WriteByte(tagCompound); err != nil {
			return err
		}
		return writeName(writer, name)
	}
}

// func WriteIntArray

// func WriteLongArray

func writeName(writer *bufio.Writer, name string) error {
	nameBS := [unsafe.Sizeof(uint16(0))]byte{}
	binary.BigEndian.PutUint16(nameBS[:], uint16(len(name)))
	if _, err := writer.Write(nameBS[:]); err != nil {
		return err
	}
	if _, err := writer.WriteString(name); err != nil {
		return err
	}
	return nil
}
