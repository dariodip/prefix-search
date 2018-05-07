package stringcoding

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetBitData(t *testing.T) {
	assert := assert.New(t)
	s1 := "a" // 01100001
	s2 := "aa" //01100001 01100001

	b1, e1 := getBitData(s1)
	assert.Nil(e1, "Error in conversion 'a'")
	assert.Equal(uint64(8), b1.Len,"Length of 'a' in bits should be 7")
	s1check := []bool{true, false, false, false, false, true, true} // binary for a (reverse)
	for i, bitCheck := range s1check {
		check1bit, _ := b1.Bits.GetBit(uint64(i))
		assert.Equal(bitCheck, check1bit, "Not equals in position: " + string(i))
	}

	b2, e2 := getBitData(s2)
	assert.Nil(e2, "Error in conversion 'aa'")
	assert.Equal(uint64(16), b2.Len,"Length of 'aa' in bits should be 14")
	// binary for aa (reverse)
	s2check := []bool{true, false, false, false, false, true, true, false, true, false, false, false, false,
	true, true}
	for i, bitCheck := range s2check {
		check2bit, _ := b2.Bits.GetBit(uint64(i))
		assert.Equal(bitCheck, check2bit, "Not equals in position: " + string(i))
	}
}


func TestGetLengthInBit(t *testing.T)  {
	assert := assert.New(t)
	s1 := "ciao" // 4 * 8 bit
	assert.Equal(getLengthInBit(s1), uint64(4*8))

	s2 := "∂iao" // 3 * 8 bit + 24
	assert.Equal(getLengthInBit(s2), uint64(3*8+24))

	s3 := "世iao" // 3 * 8 bit + 24
	assert.Equal(getLengthInBit(s3), uint64(3*8+24))

	c1 := "a"
	assert.Equal(getLengthInBit(c1), uint64(8))
}