package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"prefix-search/prefix-search/bititerator"
	"errors"
	"fmt"
)

type Coding struct {
	// Strings consists of all the concatenated bit sequences
	// corresponding to the suffixes L[i] of S's strings.
	Strings bitarray.BitArray
	// Starts consists of a sequence of bits in which each bit
	// set to 1 marks the first bit of each of those suffixes
	// in the aforementioned array (Strings).
	Starts bitarray.BitArray
	// Lengths encodes in unary the length of the shared prefixes
	// between consecutive strings.
	Lengths bitarray.BitArray
	// LastString contains the last processed string as a sequence
	// of bit. It can be used for a more efficient processing of the
	// current string to deal with.
	LastString bitarray.BitArray
	// LastIndex marks the last index in Strings (risp. Starts) arrays.
	LastIndex uint64
	// LastLengthsIndex marks the last index in the Lengths array.
	LastLengthsIndex uint64
}

// Create creates and returns a new Coding structure inserting the strings
// that are in the array of strings.
func New(strings []string, lenCalc func(uint, uint) uint) Coding {
	fc := Coding{LastIndex:0}
	return fc
}

// add adds the string s to the structure
func (c *Coding) add(s string) error {

	return nil
}

func getBitArray(s string) (bitarray.BitArray, uint64, error) {
	btarr := bitarray.NewBitArray(getLengthInBit(s))
	lastIndex := uint64(0)
	bitit := bititerator.NewStringToBitIterator(s)
	for bitit.HasNext() {
		bit, err := bitit.Next()
		if err != nil {
			return nil, uint64(0), err
		}
		if bit {
			btarr.SetBit(lastIndex)
		}
		lastIndex++
	}
	return btarr, lastIndex + 1, nil
}

func getDifferentSuffix(s1 bitarray.BitArray, s2 bitarray.BitArray, l1 uint64, l2 uint64) (bitarray.BitArray, uint64, error) {
	minLen := uint64(0)
	if l1 < l2 {
		minLen = l1
	} else {
		minLen = l2
	}
	commonPrefixLen := uint64(0)
	for i:=uint64(0);i<minLen;i++ {
		bit1, e1 := s1.GetBit(i)
		bit2, e2 := s2.GetBit(i)
		if e1 != nil || e2 != nil {
			return nil, uint64(0),  errors.New("Cannot access bitarray in position: " + string(i))
		}
		if bit1 == bit2 {
			commonPrefixLen++
		} else {
			break
		}
	}
	suffixLen := l2 - commonPrefixLen
	fmt.Println(suffixLen)
	differentSuffix := bitarray.NewBitArray(suffixLen)
	for i:=commonPrefixLen;i<l2;i++ {
		if bit, err := s2.GetBit(i); err != nil {
			if bit {
				differentSuffix.SetBit(i)
			}
		} else {
			return nil, uint64(0), err
		}
	}
	return differentSuffix, suffixLen, nil
}

func getLengthInBit(s string) uint64 {
	return uint64(len([]byte(s)) * 8)
}

