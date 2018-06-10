// Package bitdata contains types and methods in order to deal with array of bits
package bitdata

import (
	"bytes"
	"fmt"
	"github.com/dariodip/prefix-search/prefix-search/bititerator"
	"github.com/golang-collections/go-datastructures/bitarray"
)

// BitData type abstracts a string accessible as a sequence of bits;
// it also contains information about the number of bits.
type BitData struct {
	// List of bits representing some data.
	bits bitarray.BitArray
	// Number of significant bits in the BitArray.
	Len uint64
}

// New returns a pointer to a new BitData structure of the specified size.
func New(ba bitarray.BitArray, len uint64) *BitData {
	return &BitData{ba, len}
}

// GetBitData , given a string 's', returns a pointer to a BitData
// encoding the string s. If something has gone wrong, returns
// a nil pointer and and error.
func GetBitData(s string) (*BitData, error) {
	var (
		sBitLen = GetLengthInBit(s)                     // length in bit of the string
		btdata  = New(bitarray.NewBitArray(sBitLen), 0) // empty BitData
		bitit   = bititerator.NewStringToBitIterator(s) // BitIterator on the string s
	)

	for bitit.HasNext() {
		bit, err := bitit.Next()
		if err != nil {
			panic(err)
		}
		btdata.AppendBit(bit)
	}
	return btdata, nil
}

// GetLengthInBit returns the length in bit of the string s.
func GetLengthInBit(s string) uint64 {
	return uint64(len([]byte(s)) * 8)
}

// GetBit returns true if the bit in position 'index' is 1, false otherwise.
// It returns an error if something has gone wrong
func (s1 *BitData) GetBit(index uint64) (bool, error) {
	if index >= s1.Len {
		return false, ErrIndexOutOfBound
	}
	return s1.bits.GetBit(index)
}

// AppendBits appends the bits of a BitData (s2) onto the s1 BitData.
func (s1 *BitData) AppendBits(s2 *BitData) error {
	for i := uint64(0); i < s2.Len; i++ { // for each bit in the current string
		bit, err := s2.GetBit(i) // get the i-th bit in s2
		if err != nil {
			return err
		}
		err = s1.AppendBit(bit) // then append it to s1
		if err != nil {
			return err
		}
	}
	return nil
}

// AppendBit appends a bit given as a bool to the s1 BitData.
func (s1 *BitData) AppendBit(bit bool) error {
	if bit { // if the bit to append is 1
		err := s1.bits.SetBit(s1.Len) // append it to the next unmarked bit
		if err != nil {
			panic(err)
		}
	}
	s1.Len++ // always increment length
	return nil
}

// SetBit sets a bit in the BitData if the index is not out of bound.
// It won't resize the structure so Len will be as before.
func (s1 *BitData) SetBit(index uint64) error {
	if index >= s1.Len {
		return ErrIndexOutOfBound
	}
	err := s1.bits.SetBit(index)
	if err != nil {
		return err
	}
	return nil
}

// ClearBit reset a bit in the BitData if the index is not out of bound.
// It won't resize the structure so Len will be as before.
func (s1 *BitData) ClearBit(index uint64) error {
	if index >= s1.Len {
		return ErrIndexOutOfBound
	}
	err := s1.bits.ClearBit(index)
	if err != nil {
		return err
	}
	return nil
}

// GetDifferentSuffix ,given another pointer to a BitData, returns a new
// BitData containing the suffix that is not equal between the two BitDatas.
// If something goes wrong, returns a nil pointer and an error.
func (s1 *BitData) GetDifferentSuffix(s2 *BitData) (*BitData, error) {
	var (
		commonPrefixLen uint64       // length of the common prefix
		idx1            = s1.Len - 1 // last bit of the first "bitted" string
		idx2            = s2.Len - 1 // last bit of the second "bitted" string
	)
	if s1.bits == nil {
		return nil, ErrNotInitBitData
	}
	if s2.bits == nil || s2.Len == uint64(0) { // trying to find common prefix between a string and a void one
		differentSuffix := New(bitarray.NewBitArray(s1.Len), 0)
		for i := uint64(0); i < s1.Len; i++ { // we must copy the first array!
			bit, errGet := s1.GetBit(i)
			if errGet != nil {
				return nil, errGet
			}
			errAppend := differentSuffix.AppendBit(bit)
			if errAppend != nil {
				return nil, errAppend
			}
		}
		return differentSuffix, nil
	}
	// we must keep both the indexes in order to avoid out of bound
	for {
		bit1, e1 := s1.GetBit(idx1) // get bits in position idx1 (risp. idx2) on both strings
		bit2, e2 := s2.GetBit(idx2)
		if e1 != nil || e2 != nil { // something has gone wrong
			return nil, &ErrInvalidPosition{idx1}
		}
		if bit1 == bit2 { // bits are still equal
			commonPrefixLen++
		} else {
			break // bits are not still equal, we are in the different suffix
		}
		if idx1 == 0 || idx2 == 0 {
			break // no more bit to scan
		}
		idx1--
		idx2--
	}

	var (
		suffixLen       = s2.Len - commonPrefixLen                // length of the different suffix
		differentSuffix = New(bitarray.NewBitArray(suffixLen), 0) // init a new BitData to keep suffix
	)
	for i := uint64(0); i < suffixLen; i++ {
		if bit, err := s2.GetBit(i); err == nil {
			differentSuffix.AppendBit(bit)
		} else {
			return nil, err
		}
	}
	return differentSuffix, nil
}

// GetDifferentPrefix ,given another pointer to a BitData, returns a new
// BitData containing the prefix that is not equal between the two BitDatas.
// If something goes wrong, returns a nil pointer and an error.
func (s1 *BitData) GetDifferentPrefix(s2 *BitData) (*BitData, error) {
	var (
		commonSuffixLen uint64      // length of the common prefix
		idx1            = uint64(0) // first bit of the first "bitted" string
		idx2            = uint64(0) // first bit of the second "bitted" string
	)
	if s1.bits == nil {
		return nil, ErrNotInitBitData
	}
	if s2.bits == nil || s2.Len == uint64(0) { // trying to find common suffix between a string and a void one
		differentPrefix := New(bitarray.NewBitArray(s1.Len), 0)
		for i := uint64(0); i < s1.Len; i++ { // we must copy the first array!
			bit, errGet := s1.GetBit(i)
			if errGet != nil {
				return nil, errGet
			}
			errAppend := differentPrefix.AppendBit(bit)
			if errAppend != nil {
				return nil, errAppend
			}
		}
		return s1, nil
	}
	// we must keep both the indexes in order to avoid out of bound
	for idx1 != s1.Len && idx2 != s2.Len {
		bit1, e1 := s1.GetBit(idx1) // get bits in position idx1 (risp. idx2) on both strings
		bit2, e2 := s2.GetBit(idx2)
		if e1 != nil || e2 != nil { // something has gone wrong
			return nil, &ErrInvalidPosition{idx1}
		}
		if bit1 == bit2 { // bits are still equal
			commonSuffixLen++
		} else {
			break // bits are not still equal, we are in the different suffix
		}
		idx1++
		idx2++
	}

	var (
		prefixLen       = s2.Len - commonSuffixLen                // length of the different suffix
		differentPrefix = New(bitarray.NewBitArray(prefixLen), 0) // init a new BitData to keep suffix
	)
	for i := uint64(0); i < prefixLen; i++ {
		if bit, err := s2.GetBit(commonSuffixLen + i); err == nil {
			differentPrefix.AppendBit(bit)
		} else {
			return nil, err
		}
	}
	return differentPrefix, nil
}

// BitToByte returns a byte array in which each byte represents
// a character of the string at first stored as BitData.
// If something has gone wrong it returns a nil array an an error.
func (s1 *BitData) BitToByte() ([]byte, error) {
	if s1.Len%8 != 0 {
		return nil, ErrInvalidString
	}
	var (
		lenInBytes  = s1.Len / 8               // number of characters in the string
		finalBytes  = make([]byte, lenInBytes) // byte array containing characters of the encoded string
		currentByte = lenInBytes - 1           // current byte at the iteration; initially the last one
		currentBit  uint                       // current bit at the iteration; initially 0
	)
	for i := uint64(0); i < s1.Len; i++ {
		bit, err := s1.GetBit(i) // retrieve current bit
		if err != nil {          // an error occurred
			return nil, err
		}
		if bit { // bit set to 1
			finalBytes[currentByte] |= 1 << currentBit // let's set the correspective in the byte
		}
		if currentBit == 7 { // last bit, let's switch to the other byte
			currentByte--
		}
		currentBit = (currentBit + 1) % 8 // cyclic decrement
	}
	return finalBytes, nil
}

// BitToString returns a decoded string given a BitData.
// If something has gone wrong it returns a nil string and an error.
func (s1 *BitData) BitToString() (string, error) {
	bt, err := s1.BitToByte()
	if err != nil {
		return "", err
	}
	return bytes.NewBuffer(bt).String(), nil
}

// BitToTrimmedString returns a decoded and trimmed string given a BitData.
// If something has gone wrong it returns a nil string and an error.
func (s1 *BitData) BitToTrimmedString() (string, error) {
	var (
		bt, err = s1.BitToByte()
	)
	if len(bt) == 0 {
		return "", nil
	}
	for bt[0] == byte(0) {
		bt = bt[1:]
	}
	r := len(bt) - 1
	for bt[r] == byte(0) {
		bt = bt[:r]
		r = len(bt) - 1
	}
	if err != nil {
		return "", err
	}

	return bytes.NewBuffer(bt).String(), nil
}

// BitToStringOfLengthL returns a decoded string of length l bits given a BitData.
// If something has gone wrong it returns a nil string and an error.
func (s1 *BitData) BitToStringOfLengthL(l uint64) (string, error) {
	var (
		bytesCount = l / 8 // l as bytes
		bt, err    = s1.BitToByte()
	)
	if err != nil {
		return "", err
	}
	return bytes.NewBuffer(bt[:bytesCount]).String(), nil
}

// GetFirstLBits return the first l bits of a BitData
func (s1 *BitData) GetFirstLBits(l uint64) (*BitData, error) {
	var (
		newBd = New(bitarray.NewBitArray(l), l)
	)

	for i := uint64(0); i < l; i++ {
		lastBit, err := s1.GetBit(s1.Len - 1 - i)
		if err != nil {
			return nil, err
		}
		if lastBit {
			newBd.SetBit(newBd.Len - 1 - i)
		}
	}
	return newBd, nil

}

// GetTotalBitCount , given a slice of string, returns the
// total count of bits for each string in the slice.
func GetTotalBitCount(strings []string) uint64 {
	var totalBitLen uint64
	for _, s := range strings {
		totalBitLen += GetLengthInBit(s)
	}

	return totalBitLen
}

// Select1 (B,i) with 1 <= i <= n returns the position in B of the i-th occurrence of 1.
func (s1 *BitData) Select1(i uint64) (uint64, error) {
	var (
		onesCount uint64 // number 1s found
	)

	// assertion check on the index
	if err := checkIndex(s1, i); err != nil {
		return uint64(0), err
	}

	// let's iterate on the array
	for j := uint64(0); j < s1.Len; j++ {
		if bit, err := s1.GetBit(j); err == nil {
			if bit {
				onesCount++
			}
		} else { // error == nil
			return uint64(0), err
		}
		if onesCount == i { // found the i-th occurrence of 1
			return j, nil
		}
	}
	return uint64(0), ErrLessThanIOnes
}

// Rank1 (B,i) returns the number of 1s in the prefix B[1...i] aka B[0...i-1].
func (s1 *BitData) Rank1(i uint64) (uint64, error) {
	var (
		onesCount uint64 // number 1s found
	)

	// assertion check on the index
	if err := checkIndex(s1, i); err != nil {
		return uint64(0), err
	}

	// let's iterate on the array
	for j := uint64(0); j < i; j++ {
		if bit, err := s1.GetBit(j); err == nil {
			if bit {
				onesCount++
			}
		} else { // error == nil
			return uint64(0), err
		}
	}
	return onesCount, nil
}

func checkIndex(s1 *BitData, i uint64) error {
	// invalid i check
	if i > s1.Len {
		return ErrInvalidI
	}
	if i == uint64(0) {
		return ErrZeroI
	}
	return nil
}

func (s1 *BitData) String() string {
	s := ""
	i := s1.Len - 1
	for {
		if bit, err := s1.GetBit(i); err != nil {
			return err.Error() + " in String()"
		} else {
			if bit {
				s += "1"
			} else {
				s += "0"
			}
		}

		if i == 0 {
			break
		}
		i--
	}
	return fmt.Sprintf("type: %T, bits:%v, Len:%v, readableBitData:%s", s1, s1.bits, s1.Len, s)
}
