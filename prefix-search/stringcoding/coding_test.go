package stringcoding

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

/*
func TestGetDifferentSuffixWithSameLength(t *testing.T) {
	assert := assert.New(t)
	s1, l1 := "cia", getLengthInBit("cia")
	s2, l2 := "cic", getLengthInBit("cic")
	// The different suffix should be "11"
	assert.NotEqual(s1, s2, "Strings should be not equal")
	assert.Equal(l1, l2, "But their length should be equal")

	b1, e1 := getBitArray(s1)
	b2, e2 := getBitArray(s2)
	assert.Nil(e1, "Error on converting first string")
	assert.Nil(e2, "Error on converting second string")


}
*/
func TestGetBitArray(t *testing.T) {
	assert := assert.New(t)
	s1 := "a" // 01100001 = 97
//	s2 := "c" // 01100011 = 99

	b1, l1, e1 := getBitArray(s1)
	assert.Nil(e1, "Error in conversion 'a'")
	assert.Equal(l1, uint64(7), "Length of 'a' in bits should be 8")
	fmt.Println(b1.ToNums())
	assert.Equal(len(b1.ToNums()), 1, "Numbs should be of length 1")
	n1 := b1.ToNums()[0]
	assert.NotNil(n1, "Number should be not nil")
	assert.Equal(n1, uint64(97), "Number should be 97")
}


func TestGetLengthInBit(t *testing.T)  {
	assert := assert.New(t)
	s1 := "ciao" // 4 * 8 bit
	assert.Equal(getLengthInBit(s1), uint64(4*8))

	s2 := "∂iao" // 3 * 8 bit + 24
	assert.Equal(getLengthInBit(s2), uint64(3*8+24))

	s3 := "世iao" // 3 * 8 bit + 24
	assert.Equal(getLengthInBit(s3), uint64(3*8+24))

}