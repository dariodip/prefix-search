package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"prefix-search/prefix-search/bititerator"
	"errors"
)

type BitData struct {
	// List of bits representing some data
	Bits bitarray.BitArray
	// Number of significant bits in the BitArray
	Len uint64
}

type Coding struct {
	// Strings consists of all the concatenated bit sequences
	// corresponding to the suffixes L[i] of S's strings.
	Strings BitData
	// Starts consists of a sequence of bits in which each bit
	// set to 1 marks the first bit of each of those suffixes
	// in the aforementioned array (Strings).
	Starts BitData
	// Lengths encodes in unary the length of the shared prefixes
	// between consecutive strings.
	Lengths BitData
	// LastString contains the last processed string as a sequence
	// of bit. It can be used for a more efficient processing of the
	// current string to deal with.
	LastString BitData
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

func getBitArray(s string) (BitData, error) {
	btarr := BitData{bitarray.NewBitArray(getLengthInBit(s)), getLengthInBit(s)}
	lastIndex := uint64(0)
	bitit := bititerator.NewStringToBitIterator(s)
	for bitit.HasNext() {
		bit, err := bitit.Next()
		if err != nil {
			// Should we return a nil value? In that case the return value of the function must be (*BitData, error)
			return BitData{}, err
		}
		if bit {
			btarr.Bits.SetBit(lastIndex)
		}
		lastIndex++
	}

	return btarr, nil
}

func getDifferentSuffix(s1 BitData, s2 BitData) (BitData, error) {
	commonPrefixLen := uint64(0)

	idx1:=s1.Len
	idx2:=s2.Len
	for idx1>=0 && idx2>=0 {
		bit1, e1 := s1.Bits.GetBit(idx1)
		bit2, e2 := s2.Bits.GetBit(idx2)
		if e1 != nil || e2 != nil {
			return BitData{}, errors.New("Cannot access bitarray in position: " + string(idx1))
		}
		if bit1 == bit2 {
			commonPrefixLen++
		} else {
			break
		}
		idx1--
		idx2--
	}

	suffixLen := s2.Len - commonPrefixLen + 1
	differentSuffix := BitData{bitarray.NewBitArray(suffixLen), suffixLen}
	for i:=uint64(0);i<differentSuffix.Len;i++ {
		if bit, err := s2.Bits.GetBit(i); err == nil {
			if bit {
				differentSuffix.Bits.SetBit(i)
			}
		} else {
			return BitData{}, err
		}
	}
	return differentSuffix, nil
}

func getLengthInBit(s string) uint64 {
	return uint64(len([]byte(s)) * 8)
}

