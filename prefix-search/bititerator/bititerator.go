// Package bititerator provides functions and methods to create
// and deal with an iterator of bit given a string.
package bititerator

import (
	"errors"
)

type StringToBit struct {
	// string (still decoded) on which iterate.
	s string
	// current byte in the iterator.
	currentByte uint
	// string converted in a binary sequence on whose bits we iterate.
	encodedString []byte
	// current bit on the byte.
	currentBit uint
	// are we on the last bit?
	isLast bool
}

// StringToBitIterator creates an iterator that iterates upon the bits
// of the encoded version of a given string.
type StringToBitIterator interface {
	// Next returns true if the next bit is 1, false otherwise.
	Next() (bool, error)
	// HasNext returns true if the sequence has another bit.
	HasNext() bool
}

// NewStringToBitIterator creates a StringToBit iterator (of type StringToBitIterator)
// and returns it.
func NewStringToBitIterator(s string) StringToBitIterator {
	stb := StringToBit{s, stringByteLen(s) - 1, []byte(s), 0, false}
	var stbi StringToBitIterator
	stbi = &stb
	return stbi
}

// stringByteLen returns the size in bytes of the string. Supports unicode strings
func stringByteLen(s string) uint {
	return uint(len([]byte(s))) // Converts our string in a slice of bytes and returns length in bit of the string
}

// Next returns true if the next bit is 1, false otherwise.
func (bt *StringToBit) Next() (bool, error) {
	if !bt.HasNext() {
		return false, errors.New("no more bits")
	}
	toRet := uint(bt.encodedString[bt.currentByte])&(1<<bt.currentBit) == 1<<bt.currentBit
	if bt.currentBit == 7 { // we reached last bit, let's switch to the next byte
		if bt.currentByte == 0 {
			bt.isLast = true
		}
		bt.currentByte--
	}
	bt.currentBit = (bt.currentBit + 1) % 8 // cyclic increment
	return toRet, nil

}

// HasNext returns true if the sequence has another bit.
func (bt *StringToBit) HasNext() bool {
	return !bt.isLast
}
