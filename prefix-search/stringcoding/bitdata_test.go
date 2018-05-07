package stringcoding

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetDifferentSuffixWithSameLength(t *testing.T) {
	assert := assert.New(t)
	t.Logf("Test GetDifferentSuffixWithSameLength started! \n")
	s1, l1 := "cia", getLengthInBit("cia")
	s2, l2 := "cic", getLengthInBit("cic")
	assert.NotEqual(s1, s2, "Strings should be not equal")
	assert.Equal(l1, l2, "But their length should be equal")

	b1, e1 := getBitData(s1)
	b2, e2 := getBitData(s2)
	assert.Nil(e1, "Error on converting first string")
	assert.Nil(e2, "Error on converting second string")
	assert.Equal(b1.Len, b2.Len, "Strings have the same length")
	assert.NotEqual(b1.Bits, b2.Bits, "Bitarrays are different")

	expectedSuffix := []bool{true, true}  // The different suffix should be "11"
	receivedSuffix, err := b1.getDifferentSuffix(b2)
	assert.Nil(err, "Unexpected error")
	assert.Equal(uint64(2), receivedSuffix.Len, "Prefix should be of 2 bits")
	for i := uint64(0);i<receivedSuffix.Len;i++ {
		receivedSuffixBit, err := receivedSuffix.Bits.GetBit(uint64(i))
		assert.Nil(err, "An error should not be catched")
		assert.Equal(expectedSuffix[i], receivedSuffixBit, "Suffix is as expected")
	}
}
