package stringcoding

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
	"github.com/golang-collections/go-datastructures/bitarray"
)

// TODO doc
func TestBitData_GetBit(t *testing.T) {
	var (
		assert = assert.New(t)
		checkBit = []bool{false, true, false}
		expectedLength uint64
		bd = NewBitData(bitarray.NewBitArray(8), 0)

)
	assert.Equal(expectedLength, bd.Len, "Initially the BitData is empty")
	for _, bit := range checkBit {
		err := bd.AppendBit(bit)
		assert.Nil(err, "Error should be nil")
		expectedLength++
		assert.Equal(expectedLength, bd.Len, "Length should be %d", expectedLength)
	}
	assert.Equal(uint64(3), expectedLength, "ExpectedLegth should be 3")

	for i, checkBit := range checkBit {
		bit, err := bd.GetBit(uint64(i))
		assert.Nil(err, "Error should be nil")
		assert.Equal(checkBit, bit, "bit should be %s", checkBit)
	}
}

// TODO doc
func TestBitData_AppendBits(t *testing.T) {
	var (
		assert = assert.New(t)
		checkBit = []bool{false, true, false}
		expectedLength uint64
		bd = NewBitData(bitarray.NewBitArray(8), 0)
		bd2 = NewBitData(bitarray.NewBitArray(8), 0)
	)
	assert.Equal(expectedLength, bd.Len, "Initially the BitData is empty")
	for _, bit := range checkBit {
		err := bd.AppendBit(bit)
		assert.Nil(err, "Error should be nil")
		expectedLength++
		assert.Equal(expectedLength, bd.Len, "Length should be %d", expectedLength)
	}

	err := bd2.AppendBits(bd)
	assert.Nil(err, "Error should be nil")
	assert.Equal(bd.Len, bd2.Len, "Initially the BitData is empty")
	for i, checkBit := range checkBit {
		bit, err := bd2.GetBit(uint64(i))
		assert.Nil(err, "Error should be nil")
		assert.Equal(checkBit, bit, "bit should be %s", checkBit)
	}
}

// Unit test in order to check out if the method GetDifferentSuffix
// works on two string having the same length
func TestGetDifferentSuffixWithSameLength(t *testing.T) {
	var (
		assert = assert.New(t)
	)
	t.Logf("Test GetDifferentSuffixWithSameLength started! \n")
	s1, l1 := "cia", getLengthInBit("cia")
	s2, l2 := "cic", getLengthInBit("cic")
	assert.NotEqual(s1, s2, "Strings should be not equal")
	assert.Equal(l1, l2, "But their length should be equal")

	b1, e1 := getBitData(s1)
	b2, e2 := getBitData(s2)
	assert.Equal(l1, b1.Len, "String (1) length should be as expected")
	assert.Equal(l2, b2.Len, "String (2) length should be as expected")
	assert.Nil(e1, "Error on converting first string")
	assert.Nil(e2, "Error on converting second string")
	assert.Equal(b1.Len, b2.Len, "Strings have the same length")
	assert.NotEqual(b1.bits, b2.bits, "Bitarrays are different")

	expectedSuffix := []bool{true, true}  // The different suffix should be "11"
	receivedSuffix, err := b1.GetDifferentSuffix(b2)
	assert.Nil(err, "Unexpected error")
	assert.Equal(uint64(2), receivedSuffix.Len, "Prefix should be of 2 bits")
	for i := uint64(0);i<receivedSuffix.Len;i++ {
		receivedSuffixBit, err := receivedSuffix.bits.GetBit(uint64(i))
		assert.Nil(err, "An error should not be caught")
		assert.Equal(expectedSuffix[i], receivedSuffixBit, "Suffix is as expected")
	}
}

// Unit test in order to check out if the method GetDifferentSuffix
// works on two string having different length
func TestGetDifferentSuffixWithDifferentLength(t *testing.T) {
	var (
		assert = assert.New(t)
	)
	t.Logf("Test GetDifferentSuffixWithSameLength started! \n")
	s1, l1 := "cia", getLengthInBit("cia")
	s2, l2 := "cica", getLengthInBit("cica")
	assert.NotEqual(s1, s2, "Strings should be not equal")
	assert.NotEqual(l1, l2, "Their length should not be equal")

	b1, e1 := getBitData(s1)
	b2, e2 := getBitData(s2)
	assert.Nil(e1, "Error on converting first string")
	assert.Nil(e2, "Error on converting second string")
	assert.NotEqual(b1.Len, b2.Len, "Strings have not the same length")
	assert.NotEqual(b1.bits, b2.bits, "Bitarrays are different")

	// The different suffix should be the reverse of "11 01100001"
	expectedSuffix := []bool{true, false, false, false, false, true, true, false, true, true}
	receivedSuffix, err := b1.GetDifferentSuffix(b2)
	assert.Nil(err, "Unexpected error")
	assert.Equal(uint64(10), receivedSuffix.Len, "Prefix should be of 10 bits")
	for i := uint64(0);i<receivedSuffix.Len;i++ {
		receivedSuffixBit, err := receivedSuffix.bits.GetBit(uint64(i))
		assert.Nil(err, "An error should not be occurred")
		assert.Equal(expectedSuffix[i], receivedSuffixBit, "Suffix should be as expected")
	}
}

// Unit test in order to check out if the method bitToByte
// works on a string of 4 characters
func TestBitToByte(t *testing.T) {
	assert := assert.New(t)
	const s = "ciao"
	b, e := getBitData(s)
	assert.Nil(e)
	assert.Equal(uint64(4*8), b.Len)
	bt, e := b.bitToByte()
	assert.Nil(e)
	resultString := bytes.NewBuffer(bt).String()
	assert.Equal(s, resultString)
}

// Unit test in order to check out if the method bitToString
// works on a string of 4 characters
func TestBitToString(t *testing.T) {
	assert := assert.New(t)
	const s = "ciao"
	b, e := getBitData(s)
	assert.Nil(e)
	assert.Equal(uint64(4*8), b.Len)
	sCheck, e := b.bitToString()
	assert.Nil(e)
	assert.Equal(s, sCheck)
}
