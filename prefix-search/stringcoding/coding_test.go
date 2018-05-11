package stringcoding

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang-collections/go-datastructures/bitarray"
)

// Unit test in order to check out if the method getBitData
// works on a string of a single character
func TestGetBitDataSingleChar(t *testing.T) {
	const s1 = "a" // 01100001
	var (
		assert = assert.New(t)
		b1, e1 = getBitData(s1)
		s1check = []bool{true, false, false, false, false, true, true} // binary for a (reverse)
	)

	assert.Nil(e1, "Error in conversion 'a'")
	assert.Equal(uint64(8), b1.Len,"Length of 'a' in bits should be 7")
	for i, bitCheck := range s1check {
		check1bit, _ := b1.bits.GetBit(uint64(i))
		assert.Equal(bitCheck, check1bit, "Not equals in position: " + string(i))
	}
}

func TestAddAndGetUnaryLength(t *testing.T) {
	assert := assert.New(t)
	const (
		n1 = uint64(5)										// length to set
		expectedLength = uint64(1 + 5 + 1)					// first 0 + 5 * '1' + last 0
		expected1sCount = uint64(5)
	)
	var (
		c = New([]string{"ciaos"}, func(u, u2 uint64) uint64 {  // stub coding struct
			return 0
		})
		onesCounter uint64									// counter of 1s
		lastIndex = uint64(1)								// last index of the array
	)

	t.Log("Add unary value")
	err := c.addUnaryLength(n1)
	assert.Nil(err, "Error should be nil")
	t.Log("Get lengths array")
	lengthBitData := c.Lengths
	assert.NotNil(lengthBitData, "Lengths array should not be nil")
	assert.Equal(expectedLength, lengthBitData.Len, "Len value should be 6 (5 + 1)")
	assert.Equal(expectedLength, c.NextLengthsIndex, "NextLengthsIndex should be 6")
	for lastIndex >= 0 {
		bit, err := lengthBitData.bits.GetBit(lastIndex)
		assert.Nil(err, "Error should be nil")
		if !bit {
			break
		} else {
			onesCounter++
		}
		lastIndex++
	}
	assert.Equal(expected1sCount, onesCounter, "Array contains 5 ones")

	checkUnaryToInt, err := c.unaryToInt(uint64(0))
	assert.Nil(err, "Error should be nil")
	assert.Equal(expected1sCount, checkUnaryToInt, "Check unary to int should be 5")
}

// Unit test in order to check out if the method getBitData
// works on a string of two characters
func TestGetBitDataTwoChar(t *testing.T) {
	const s2 = "ab" //01100001 01100010
	var (
		assert = assert.New(t)
		b2, e2 = getBitData(s2)
		s2check = []bool{false, true, false, false, false, true, true, false, true, false, false, false, false,
			true, true}	// binary for ab (reverse)
	)

	assert.Nil(e2, "Error in conversion 'aa'")
	assert.Equal(uint64(16), b2.Len,"Length of 'aa' in bits should be 14")
	for i, bitCheck := range s2check {
		check2bit, _ := b2.bits.GetBit(uint64(i))
		assert.Equal(bitCheck, check2bit, "Not equals in position: " + string(i))
	}
}

// Unit test in order to check out if the method getLengthInBit
// works on a string of two characters
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

// Unit test in order to check out if the methon add works
// by inserting 1 or 2 strings
func TestCoding_Add(t *testing.T) {

}

func TestSetStartsWithOffset(t *testing.T) {
	var (
		assert = assert.New(t)
		c = New([]string{"ciao"}, func(len1, len2 uint64) (uint64) {
			return 0  // we don't need that function here, so we simply use a stub method
		})
		bd = NewBitData(bitarray.NewBitArray(8), 0)
		checkBit = []bool{false, true, false}
		expectedBits = []bool{true, false, false, true, false, false}
	)

	for i, bit := range checkBit {
		err := bd.AppendBit(bit)
		assert.Nil(err, "Cannot set bit %d", i)
	}

	err1 := c.setStartsWithOffset(bd)
	assert.Nil(err1, "Something goes wrong %s", err1)
	err2 := c.setStartsWithOffset(bd)
	assert.Nil(err2, "Something goes wrong %s", err2)

	assert.Equal(c.Starts.Len, bd.Len * 2, "Lengths are not equal. Found %d, expected %d",
											c.Starts.Len, bd.Len * 2)

	for i := uint64(0); i < c.Starts.Len; i++ {
		bit, err := c.Starts.GetBit(i)
		assert.Nil(err, "Cannot access to bit %d", i)
		assert.Equal(bit, expectedBits[i], "Wrong bit at position %d. Found %t, expected %t",
					 i, bit, expectedBits[i])
	}
}