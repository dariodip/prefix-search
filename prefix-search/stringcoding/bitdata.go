package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"errors"
	"bytes"
)

// BitData type abstracts a string accessible as a sequence of bits;
// it also contains information about the number of bits.
type BitData struct {
	// List of bits representing some data.
	Bits bitarray.BitArray
	// Number of significant bits in the BitArray.
	Len uint64
}

// NewBitData returns a pointer to a new BitData structure at the specified size.
func NewBitData(ba bitarray.BitArray, len uint64) *BitData { // TODO private in the package ?
	return &BitData{ba, len}
}

// getDifferentSuffix, given another pointer to a BitData, returns a new
// BitData containing the suffix that is not equal between the two BitDatas.
// If something goes wrong, returns a nil pointer and an error.
func (s1 *BitData) getDifferentSuffix(s2 *BitData) (*BitData, error) {
	var(
		commonPrefixLen uint64  // length of the common prefix
		idx1 = s1.Len			// length in bit of the first "bitted" string
		idx2 = s2.Len			// length in bit of the second "bitted" string
	)

	for idx1>=0 && idx2>=0 {  // we must keep both the indexes in order to avoid out of bound
		bit1, e1 := s1.Bits.GetBit(idx1) // get bits in position idx1 (risp. idx2) on both strings
		bit2, e2 := s2.Bits.GetBit(idx2)
		if e1 != nil || e2 != nil {  // something has gone wrong
			return nil, errors.New("Cannot access bitarray in position: " + string(idx1))
		}
		if bit1 == bit2 {  // bits are still equal
			commonPrefixLen++
		} else {
			break  		   // bits are not still equal, we are in the different suffix
		}
		idx1--
		idx2--
	}

	var (
		suffixLen = s2.Len - commonPrefixLen + 1									// length of the different suffix
		differentSuffix = NewBitData(bitarray.NewBitArray(suffixLen), suffixLen)	// init a new BitData to keep suffix
	)
	for i:=uint64(0);i<differentSuffix.Len;i++ {
		if bit, err := s2.Bits.GetBit(i); err == nil {
			if bit {  // bit set to 1
				differentSuffix.Bits.SetBit(i)
			}
		} else {
			return nil, err  // something has gone wrong
		}
	}
	return differentSuffix, nil
}

// bitToByte returns a byte array in which each byte represents
// a character of the string at first stored as BitData.
// If something has gone wrong it returns a nil array an an error.
func (s1 *BitData) bitToByte() ([]byte, error) {
	var (
		lenInBytes = s1.Len / 8							// number of characters in the string
		finalBytes = make([]byte, lenInBytes)			// byte array containing characters of the encoded string
		currentByte = lenInBytes - 1  					// current byte at the iteration; initially the last one
		currentBit uint									// current bit at the iteration; initially 0
	)

	for i := uint64(0); i<s1.Len; i++  {
		bit, err := s1.Bits.GetBit(i)					// retrieve current bit
		if err != nil {									// an error occurred
			return nil, err
		}
		if bit {										// bit set to 1
			finalBytes[currentByte] |= 1 << currentBit	// let's set the correspective in the byte
		}
		if currentBit == 7 {							// last bit, let's switch to the other byte
			currentByte--
		}
		currentBit = (currentBit + 1) % 8 				// cyclic decrement
	}
	return finalBytes, nil
}

// bitToByte returns a decoded string given a BitData.
// If something has gone wrong it returns a nil array an an error.
func (s1 *BitData) bitToString() (string, error)  {
	bt, err := s1.bitToByte()
	if err!=nil {
		return "", err
	}
	return bytes.NewBuffer(bt).String(), nil
}