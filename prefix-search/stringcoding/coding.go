package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"prefix-search/prefix-search/bititerator"
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
	// LastIndex marks the last index in Strings (risp. Starts) arrays.
	LastIndex uint64
	// LastLengthsIndex marks the last index in the Lengths array.
	LastLengthsIndex uint64
}

// Create creates and returns a new Coding structure inserting the strings
// that are in the array of strings.
func New(strings []string, lenCalc func(uint, uint) uint) *Coding {
	fc := Coding{LastIndex:0} // TODO
	return &fc
}

// add adds the string s to the structure
func (c *Coding) add(s string) error {
	// TODO
	return nil
}

func getBitData(s string) (*BitData, error) {
	sBitLen := getLengthInBit(s)
	btdata := NewBitData(bitarray.NewBitArray(sBitLen), sBitLen)
	lastIndex := uint64(0)
	bitit := bititerator.NewStringToBitIterator(s)
	for bitit.HasNext() {
		bit, err := bitit.Next()
		if err != nil {
			return nil, err
		}
		if bit {
			btdata.Bits.SetBit(lastIndex)
		}
		lastIndex++
	}
	return btdata, nil
}

func getLengthInBit(s string) uint64 {
	return uint64(len([]byte(s)) * 8)
}

