// stringcoding package provides an easy way to deal
// with strings in a bit-to-bit fashion
package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"prefix-search/prefix-search/bititerator"
	"errors"
)


type Coding struct {
	// Strings consists of all the concatenated bit sequences
	// corresponding to the suffixes L[i] of S's strings.
	Strings *BitData
	// Starts consists of a sequence of bits in which each bit
	// set to 1 marks the first bit of each of those suffixes
	// in the aforementioned array (Strings).
	Starts *BitData
	// Lengths encodes in unary the length of the shared prefixes
	// between consecutive strings.
	Lengths *BitData
	// LastString contains the last processed string as a sequence
	// of bit. It can be used for a more efficient processing of the
	// current string to deal with.
	LastString *BitData
	// NextIndex marks the last index in Strings (risp. Starts) arrays.
	NextIndex uint64
	// NextLengthsIndex marks the last index in the Lengths array.
	NextLengthsIndex uint64
}

// Create creates and returns a new Coding structure inserting the strings
// that are in the array of strings.
func New(strings []string, lenCalc func(uint, uint) uint) *Coding {
	maxCapacity:=getTotalBitCount(strings)
	maxLengthCapacity:=maxCapacity+uint64(len(strings)-1)
	fc := Coding{
		Strings: NewBitData(bitarray.NewBitArray(maxCapacity), 0),
		Starts:	NewBitData(bitarray.NewBitArray(maxCapacity), 0),
		Lengths: NewBitData(bitarray.NewBitArray(maxLengthCapacity), 1),
		NextLengthsIndex: uint64(1),
	} // TODO insert
	return &fc
}

// add adds the string s to the structure
func (c *Coding) add(s string) error {
	if c.LastString == nil {  // first string case

	}
	return nil
}

func (c *Coding) enqueueBitData(bd BitData) error {
	return nil
}

func (c *Coding) addUnaryLenght(n uint64) error {
	for i:=uint64(0);i<n;i++ {
		if err := c.Lengths.AppendBit(true); err != nil {
			return err
		}
		c.NextLengthsIndex++
	}
	if err := c.Lengths.AppendBit(false); err != nil {
		return err
	}
	c.NextLengthsIndex++
	return nil
}

// Given an index, returns the idx-th value of the unary array
func (c *Coding) unaryToInt(idx uint64) (uint64, error) {
	if bit, err := c.Lengths.GetBit(idx); err == nil {
		if bit && idx != uint64(0) {
			return 0, errors.New("index should point to a 0")
		}
	} else {
		return 0, err
	}

	var val uint64
	current:=idx
	for {
		current++
		if bit, err := c.Lengths.GetBit(current); err == nil {
			if bit {
				val++
			} else {
				break
			}
		} else {
			return 0, err
		}
	}

	return val, nil
}

// Given a string 's', getBitData returns a pointer to a BitData
// encoding the string s. If something has gone wrong, returns
// a nil pointer and and error.
func getBitData(s string) (*BitData, error) {
	var (
		sBitLen = getLengthInBit(s)									// length in bit of the string
		btdata = NewBitData(bitarray.NewBitArray(sBitLen), 0)	// empty BitData
		bitit = bititerator.NewStringToBitIterator(s)				// BitIterator on the string s
	)

	for bitit.HasNext() {
		bit, err := bitit.Next()
		if err != nil {
			panic(err)
			return nil, err // something has gone wrong
		}
		btdata.AppendBit(bit)
	}
	return btdata, nil
}

// Returns the length in bit of the string s.
func getLengthInBit(s string) uint64 {
	return uint64(len([]byte(s)) * 8)
}

func getTotalBitCount(strings []string) uint64 {
	var totalBitLen uint64
	for _, s:=range strings {
		totalBitLen+=getLengthInBit(s)
	}

	return totalBitLen
}

