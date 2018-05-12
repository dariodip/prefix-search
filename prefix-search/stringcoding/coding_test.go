package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"github.com/stretchr/testify/assert"
	bd "prefix-search/prefix-search/bitdata"
	"testing"
)

// Unit test in order to check out if the method getBitData
// works on a string of a single character
func TestGetBitDataSingleChar(t *testing.T) {
	const s1 = "a" // 01100001
	var (
		a       = assert.New(t)
		b1, e1  = bd.GetBitData(s1)
		s1check = []bool{true, false, false, false, false, true, true} // binary for a (reverse)
	)

	a.Nil(e1, "Error in conversion 'a'")
	a.Equal(uint64(8), b1.Len, "Length of 'a' in bits should be 7")
	for i, bitCheck := range s1check {
		check1bit, _ := b1.GetBit(uint64(i))
		a.Equal(bitCheck, check1bit, "Not equals in position: "+string(i))
	}
}

func TestAddAndGetUnaryLength(t *testing.T) {
	a := assert.New(t)
	const (
		n1              = uint64(5)         // length to set
		expectedLength  = uint64(1 + 5 + 1) // first 0 + 5 * '1' + last 0
		expected1sCount = uint64(5)
	)
	var (
		c = New([]string{"ciaos"}, func(u, u2 uint64) uint64 { // stub coding struct
			return 0
		})
		onesCounter uint64      // counter of 1s
		lastIndex   = uint64(1) // last index of the array
	)

	t.Log("Add unary value")
	err := c.addUnaryLength(n1)
	a.Nil(err, "Error should be nil")
	t.Log("Get lengths array")
	lengthBitData := c.Lengths
	a.NotNil(lengthBitData, "Lengths array should not be nil")
	a.Equal(expectedLength, lengthBitData.Len, "Len value should be 6 (5 + 1)")
	a.Equal(expectedLength, c.NextLengthsIndex, "NextLengthsIndex should be 6")
	for lastIndex >= 0 {
		bit, err := lengthBitData.GetBit(lastIndex)
		a.Nil(err, "Error should be nil")
		if !bit {
			break
		} else {
			onesCounter++
		}
		lastIndex++
	}
	a.Equal(expected1sCount, onesCounter, "Array contains 5 ones")

	checkUnaryToInt, err := c.unaryToInt(uint64(0))
	a.Nil(err, "Error should be nil")
	a.Equal(expected1sCount, checkUnaryToInt, "Check unary to int should be 5")
}

// Unit test in order to check out if the method getBitData
// works on a string of two characters
func TestGetBitDataTwoChar(t *testing.T) {
	const s2 = "ab" //01100001 01100010
	var (
		a       = assert.New(t)
		b2, e2  = bd.GetBitData(s2)
		s2check = []bool{false, true, false, false, false, true, true, false, true, false, false, false, false,
			true, true} // binary for ab (reverse)
	)

	a.Nil(e2, "Error in conversion 'aa'")
	a.Equal(uint64(16), b2.Len, "Length of 'aa' in bits should be 14")
	for i, bitCheck := range s2check {
		check2bit, _ := b2.GetBit(uint64(i))
		a.Equal(bitCheck, check2bit, "Not equals in position: "+string(i))
	}
}

// Unit test in order to check out if the method GetLengthInBit
// works on a string of two characters
func TestGetLengthInBit(t *testing.T) {
	a := assert.New(t)
	s1 := "ciao" // 4 * 8 bit
	a.Equal(bd.GetLengthInBit(s1), uint64(4*8))

	s2 := "∂iao" // 3 * 8 bit + 24
	a.Equal(bd.GetLengthInBit(s2), uint64(3*8+24))

	s3 := "世iao" // 3 * 8 bit + 24
	a.Equal(bd.GetLengthInBit(s3), uint64(3*8+24))

	c1 := "a"
	a.Equal(bd.GetLengthInBit(c1), uint64(8))
}

// Unit test in order to check out if the method add works
// by inserting 1 or 2 strings
func TestCoding_Add(t *testing.T) {
	const (
		s1 = "cia"
		s2 = "cic"
	)

	var (
		a       = assert.New(t)
		lenCalc = func(prefixLen, stringLen uint64) uint64 {
			return stringLen - prefixLen
		}
		c           = New([]string{s1, s2}, lenCalc)
		stringsBits = []bool{
			true, false, false, false, false, true, true, false,
			true, false, false, true, false, true, true, false,
			true, true, false, false, false, true, true, false,
			true, true,
		}

		startsBits = []bool{
			true, false, false, false, false, false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false, false, false, false, true, false,
		}
	)

	a.NotEqual(s1, s2, "The two test strings must not be equal")

	c.add(s1)

	// Strings and LastString check, since we have only one strings for now we can use getBitData
	s1bits, err := bd.GetBitData(s1)
	a.Nil(err, "Something goes wrong while converting %s in a BitData: %s", s1, err)
	a.Equal(s1bits, c.Strings, "Wrong conversion on Strings") // TODO
	a.Equal(s1bits.Len, c.Strings.Len, "Wrong len on Strings")

	a.Equal(s1bits, c.LastString, "Wrong conversion on LastString, should be equal to %s",
		s1) // TODO
	a.Equal(s1bits.Len, c.LastString.Len, "Wrong len on LastString, should be %d", s1bits.Len)

	// Starts check
	bit, err := c.Starts.GetBit(0)
	a.Nil(err, "Cannot access bit %d. %s", 0, err)
	a.Equal(bit, true, "Wrong bit at position %d. Found %t, expected %t", 0, bit, true)
	for i := uint64(1); i < c.Starts.Len; i++ {
		bit, err := c.Starts.GetBit(i)
		a.Nil(err, "Cannot access bit %d. %s", i, err)
		a.Equal(bit, false, "Wrong bit at position %d. Found %t, expected %t", i, bit, false)
	}

	// Lenghts
	s1len, err := c.unaryToInt(0)
	a.Nil(err, "Something goes wrong: %s", err)
	a.Equal(s1len, bd.GetLengthInBit(s1), "Some bit are missing in Lenghts. Found %d, expected %d", s1len,
		len(s1))

	c.add(s2)

	// Check if lastString is now s2
	s2bits, err := bd.GetBitData(s2)
	a.Nil(err, "Something goes wrong while converting %s in a BitData: %s", s1, err)
	a.Equal(s2bits, c.LastString, "Wrong conversion on LastString, should be equal to %s",
		s2)
	a.Equal(s2bits.Len, c.LastString.Len, "Wrong len on LastString, should be %d", s2bits.Len)

	a.Equal(c.Strings.Len, uint64(len(stringsBits)), "String len should be %d, not %d",
		uint64(len(stringsBits)), c.Strings.Len)

	// Strings check
	for i := uint64(0); i < c.Strings.Len; i++ {
		bit, err := c.Strings.GetBit(i)
		a.Nil(err, "Cannot access bit %d. %s", i, err)
		a.Equal(bit, stringsBits[i], "Wrong bit at position %d. Found %t, expected %t",
			i, bit, stringsBits[i])
	}

	// Starts check
	for i := uint64(0); i < c.Starts.Len; i++ {
		bit, err := c.Starts.GetBit(i)
		a.Nil(err, "Cannot access bit %d. %s", i, err)
		a.Equal(bit, startsBits[i], "Wrong bit at position %d. Found %t, expected %t",
			i, bit, startsBits[i])
	}

	s2suffixLen, err := c.unaryToInt(bd.GetLengthInBit(s1) + 1)
	a.Nil(err, "Something goes wrong: %s", err)
	a.Equal(s2suffixLen, uint64(2), "Some bit are missing in Lenghts. Found %d, expected %d",
		s2suffixLen, uint64(2))
}

func TestSetStartsWithOffset(t *testing.T) {
	var (
		a = assert.New(t)
		c = New([]string{"ciao"}, func(len1, len2 uint64) uint64 {
			return 0 // we don't need that function here, so we simply use a stub method
		})
		bitD         = bd.New(bitarray.NewBitArray(8), 0)
		checkBit     = []bool{false, true, false}
		expectedBits = []bool{true, false, false, true, false, false}
	)

	for i, bit := range checkBit {
		err := bitD.AppendBit(bit)
		a.Nil(err, "Cannot set bit %d", i)
	}

	err1 := c.setStartsWithOffset(bitD)
	a.Nil(err1, "Something goes wrong %s", err1)
	err2 := c.setStartsWithOffset(bitD)
	a.Nil(err2, "Something goes wrong %s", err2)

	a.Equal(c.Starts.Len, bitD.Len*2, "Lengths are not equal. Found %d, expected %d",
		c.Starts.Len, bitD.Len*2)

	for i := uint64(0); i < c.Starts.Len; i++ {
		bit, err := c.Starts.GetBit(i)
		a.Nil(err, "Cannot access to bit %d", i)
		a.Equal(bit, expectedBits[i], "Wrong bit at position %d. Found %t, expected %t",
			i, bit, expectedBits[i])
	}
}
