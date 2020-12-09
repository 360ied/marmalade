package nbt

import (
	"bufio"
	"errors"
	"fmt"
)

var (
	WrongTypeIDError    = errors.New("marmalade_nbt: wrong type id")
	InvalidTypeIDError  = errors.New("marmalade_nbt: invalid type id")
	InvalidLengthError  = errors.New("marmalade_nbt: invalid length")
	NotImplementedError = errors.New("marmalade_nbt: not implemented")
)

func newWrongTypeIDError(expected, got byte) error {
	return fmt.Errorf("%w, expected %v, but got %v", WrongTypeIDError, expected, got)
}

func assertTypeID(reader *bufio.Reader, expect byte) error {
	got, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if got != expect {
		return newWrongTypeIDError(expect, got)
	} else {
		return nil
	}
}

func newInvalidTypeIDError(got byte) error {
	return fmt.Errorf("%w, got %v", InvalidTypeIDError, got)
}

func newInvalidLengthError(l interface{}) error {
	return fmt.Errorf("%w, got %v", InvalidLengthError, l)
}
