// stringcoding package provides an easy way to deal
// with strings in a bit-to-bit fashion
// Each compressed/uncompressed string will be represented in the structure in right to left order (from the
// least significant bit to the most significant).
// i.e: The list of two strings [cia, cio] (without loss of generality we suppose to work with characters instead
//      bits), will be encoded as [aic|o] (the vertical bar here is simply used to divide the two encoded string)
package stringcoding

import (
	"fmt"
	bd "github.com/dariodip/prefix-search/prefix-search/bitdata"
	"github.com/golang-collections/go-datastructures/bitarray"
)

type Coding struct {
	// Strings consists of all the concatenated bit sequences
	// corresponding to the suffixes L[i] of S's strings.
	Strings *bd.BitData
	// Starts consists of a sequence of bits in which each bit
	// set to 1 marks the first bit of each of those suffixes
	// in the aforementioned array (Strings).
	Starts *bd.BitData
	// Lengths encodes in unary the length of the shared prefixes
	// between consecutive strings.
	Lengths *bd.BitData
	// LastString contains the last processed string as a sequence
	// of bit. It can be used for a more efficient processing of the
	// current string to deal with.
	LastString *bd.BitData
	// NextIndex marks the last index in Strings (risp. Starts) arrays.
	NextIndex uint64
	// NextLengthsIndex marks the last index in the Lengths array.
	NextLengthsIndex uint64
}

// Create creates and returns a new Coding structure inserting the strings
// that are in the array of strings.
func New(strings []string) *Coding {
	maxCapacity := bd.GetTotalBitCount(strings)
	maxLengthCapacity, err := getEliasGammaLength(strings)
	if err != nil {
		panic("Cannot find the required capacity")
	}
	fc := Coding{
		Strings:          bd.New(bitarray.NewBitArray(maxCapacity), 0),
		Starts:           bd.New(bitarray.NewBitArray(maxCapacity), 0),
		Lengths:          bd.New(bitarray.NewBitArray(maxLengthCapacity), 0),
		NextLengthsIndex: uint64(0),
	}
	return &fc
}

// setStartsWithOffset sets the bit in the Starts bitdata in order
// to state where the suffix in Strings starts.
func (c *Coding) setStartsWithOffset(differentSuffix *bd.BitData) error {
	if differentSuffix.Len == 0 {
		return nil // nothing to do here
	}

	c.Starts.AppendBit(true)

	c.Starts.Len += differentSuffix.Len - 1

	return nil
}

func (c *Coding) String() string {
	return fmt.Sprintf("type:%T, Strings: %v, Starts:%v, Lengths:%v, LastString:%v, NextIndex:%v, NextLengthsIndex:%v",
		c, c.Strings, c.Starts, c.Lengths, c.LastString, c.NextIndex, c.NextLengthsIndex)
}
