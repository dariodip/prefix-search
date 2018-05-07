package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"errors"
)

type BitData struct {
	// List of bits representing some data
	Bits bitarray.BitArray
	// Number of significant bits in the BitArray
	Len uint64
}

func NewBitData(ba bitarray.BitArray, len uint64) *BitData {
	return &BitData{ba, len}
}


func (s1 *BitData) getDifferentSuffix(s2 *BitData) (*BitData, error) {
	commonPrefixLen := uint64(0)

	idx1:=s1.Len
	idx2:=s2.Len
	for idx1>=0 && idx2>=0 {
		bit1, e1 := s1.Bits.GetBit(idx1)
		bit2, e2 := s2.Bits.GetBit(idx2)
		if e1 != nil || e2 != nil {
			return nil, errors.New("Cannot access bitarray in position: " + string(idx1))
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
	differentSuffix := NewBitData(bitarray.NewBitArray(suffixLen), suffixLen)
	for i:=uint64(0);i<differentSuffix.Len;i++ {
		if bit, err := s2.Bits.GetBit(i); err == nil {
			if bit {
				differentSuffix.Bits.SetBit(i)
			}
		} else {
			return nil, err
		}
	}
	return differentSuffix, nil
}
