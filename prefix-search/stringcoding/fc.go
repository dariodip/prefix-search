package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"prefix-search/prefix-search/bititerator"
)

type FC struct {
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

// Create creates and returns a new FC structure inserting the strings
// that are in the array of strings.
func Create(strings []string) FC {
	fc := FC{LastIndex:0}
	return fc
}

func (fc *FC) Retrieval(u string, l uint) []string {
	ss := make([] string, 1)
	return ss
}

func (fc *FC) Select1(l uint) uint {
	return uint(1)
}

// add adds the string s to the structure
func (fc *FC) add(s string) error {
	bitit := bititerator.NewStringToBitIterator(s)
	for bitit.HasNext() {
		toSet, err := bitit.Next()
		if err != nil {
			return err
		}
		if toSet {
			fc.Strings.SetBit(fc.LastIndex)
			fc.LastIndex++
		}
	}
	return nil
}
