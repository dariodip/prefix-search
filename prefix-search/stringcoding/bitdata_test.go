package stringcoding

import (
	"bytes"
	"github.com/dariodip/prefix-search/prefix-search/bitdata"
	"github.com/golang-collections/go-datastructures/bitarray"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitData_GetBit(t *testing.T) {
	var (
		a              = assert.New(t)
		checkBit       = []bool{false, true, false}
		expectedLength uint64
		bd             = bitdata.New(bitarray.NewBitArray(8), 0)
	)
	a.Equal(expectedLength, bd.Len, "Initially the BitData is empty")
	for _, bit := range checkBit {
		err := bd.AppendBit(bit)
		a.Nil(err, "Error should be nil")
		expectedLength++
		a.Equal(expectedLength, bd.Len, "Length should be %d", expectedLength)
	}
	a.Equal(uint64(3), expectedLength, "ExpectedLegth should be 3")

	for i, checkBit := range checkBit {
		bit, err := bd.GetBit(uint64(i))
		a.Nil(err, "Error should be nil")
		a.Equal(checkBit, bit, "bit should be %s", checkBit)
	}
	//	t.Log(runtime.Caller(0))
}

func TestBitData_SetBit(t *testing.T) {
	var (
		a        = assert.New(t)
		bd       = bitdata.New(bitarray.NewBitArray(8), 0)
		checkBit = []bool{false, true, true}
	)

	for i, bit := range checkBit {
		err := bd.AppendBit(bit)
		a.Nil(err, "Cannot set bit %d", i)
	}

	bitToSet := uint64(0)
	err := bd.SetBit(bitToSet)
	a.Nil(err, "Cannot set bit %d", bitToSet)

	bit, err := bd.GetBit(bitToSet)
	a.Nil(err, "Something goes wrong")
	a.Equal(bit, true, "Bit %d has a not valid value. Expected %t, found %t",
		bitToSet, true, bit)

	a.Equal(uint64(len(checkBit)), bd.Len, "BitData len should be 3")
}

func TestBitData_AppendBits(t *testing.T) {
	var (
		a              = assert.New(t)
		checkBit       = []bool{false, true, false}
		expectedLength uint64
		bd             = bitdata.New(bitarray.NewBitArray(8), 0)
		bd2            = bitdata.New(bitarray.NewBitArray(8), 0)
	)
	a.Equal(expectedLength, bd.Len, "Initially the BitData is empty")
	for _, bit := range checkBit {
		err := bd.AppendBit(bit)
		a.Nil(err, "Error should be nil")
		expectedLength++
		a.Equal(expectedLength, bd.Len, "Length should be %d", expectedLength)
	}

	err := bd2.AppendBits(bd)
	a.Nil(err, "Error should be nil")
	a.Equal(bd.Len, bd2.Len, "Initially the BitData is empty")
	for i, checkBit := range checkBit {
		bit, err := bd2.GetBit(uint64(i))
		a.Nil(err, "Error should be nil")
		a.Equal(checkBit, bit, "bit should be %s", checkBit)
	}
}

// Unit test in order to check out if the method GetDifferentSuffix
// works on two string having the same length
func TestGetDifferentSuffixWithSameLength(t *testing.T) {
	t.Logf("Test GetDifferentSuffixWithSameLength started! \n")
	var (
		a      = assert.New(t)
		s1, l1 = "cia", bitdata.GetLengthInBit("cia")
		s2, l2 = "cic", bitdata.GetLengthInBit("cic")
	)

	a.NotEqual(s1, s2, "Strings should be not equal")
	a.Equal(l1, l2, "But their length should be equal")

	b1, e1 := bitdata.GetBitData(s1)
	b2, e2 := bitdata.GetBitData(s2)
	a.Equal(l1, b1.Len, "String (1) length should be as expected")
	a.Equal(l2, b2.Len, "String (2) length should be as expected")
	a.Nil(e1, "Error on converting first string")
	a.Nil(e2, "Error on converting second string")
	a.Equal(b1.Len, b2.Len, "Strings have the same length")
	a.NotEqual(b1, b2, "Bitarrays are different")

	expectedSuffix := []bool{true, true} // The different suffix should be "11"
	receivedSuffix, err := b1.GetDifferentSuffix(b2)
	a.Nil(err, "Unexpected error")
	a.Equal(uint64(2), receivedSuffix.Len, "Prefix should be of 2 bits")
	for i := uint64(0); i < receivedSuffix.Len; i++ {
		receivedSuffixBit, err := receivedSuffix.GetBit(uint64(i))
		a.Nil(err, "An error should not be caught")
		a.Equal(expectedSuffix[i], receivedSuffixBit, "Suffix is as expected")
	}
}

// Unit test in order to check out if the method GetDifferentSuffix
// works on two string having different length
func TestGetDifferentSuffixWithDifferentLength(t *testing.T) {
	t.Logf("Test GetDifferentSuffixWithSameLength started! \n")
	var (
		a      = assert.New(t)
		s1, l1 = "cia", bitdata.GetLengthInBit("cia")
		s2, l2 = "cica", bitdata.GetLengthInBit("cica")
	)
	a.NotEqual(s1, s2, "Strings should be not equal")
	a.NotEqual(l1, l2, "Their length should not be equal")

	b1, e1 := bitdata.GetBitData(s1)
	b2, e2 := bitdata.GetBitData(s2)
	a.Nil(e1, "Error on converting first string")
	a.Nil(e2, "Error on converting second string")
	a.NotEqual(b1.Len, b2.Len, "Strings have not the same length")
	a.NotEqual(b1, b2, "Bitarrays are different")

	// The different suffix should be the reverse of "11 01100001"
	expectedSuffix := []bool{true, false, false, false, false, true, true, false, true, true}
	receivedSuffix, err := b1.GetDifferentSuffix(b2)
	a.Nil(err, "Unexpected error")
	a.Equal(uint64(10), receivedSuffix.Len, "Prefix should be of 10 bits")
	for i := uint64(0); i < receivedSuffix.Len; i++ {
		receivedSuffixBit, err := receivedSuffix.GetBit(uint64(i))
		a.Nil(err, "An error should not be occurred")
		a.Equal(expectedSuffix[i], receivedSuffixBit, "Suffix should be as expected")
	}
}

// Unit test in order to check out if the method GetDifferentPrefix
// works on two string having the same length
func TestGetDifferentPrefixWithSameLength(t *testing.T) {
	t.Logf("Test GetDifferentPrefixWithSameLength started! \n")
	var (
		a      = assert.New(t)
		s1, l1 = "cia", bitdata.GetLengthInBit("cia")
		s2, l2 = "aia", bitdata.GetLengthInBit("aia")
	)

	a.NotEqual(s1, s2, "Strings should be not equal")
	a.Equal(l1, l2, "But their length should be equal")

	b1, e1 := bitdata.GetBitData(s1)
	b2, e2 := bitdata.GetBitData(s2)
	a.Equal(l1, b1.Len, "String (1) length should be as expected")
	a.Equal(l2, b2.Len, "String (2) length should be as expected")
	a.Nil(e1, "Error on converting first string")
	a.Nil(e2, "Error on converting second string")
	a.Equal(b1.Len, b2.Len, "Strings have the same length")
	a.NotEqual(b1, b2, "Bitarrays are different")

	// The different prefix should be "0110000"
	expectedPrefix := []bool{false, false, false, false, true, true, false}
	receivedPrefix, err := b1.GetDifferentPrefix(b2)
	a.Nil(err, "Unexpected error")
	a.Equal(uint64(len(expectedPrefix)), receivedPrefix.Len, "Prefix should be of %d bits",
		uint64(len(expectedPrefix)))
	for i := uint64(0); i < receivedPrefix.Len; i++ {
		receivedPrefixBit, err := receivedPrefix.GetBit(uint64(i))
		a.Nil(err, "An error should not be caught")
		a.Equal(expectedPrefix[i], receivedPrefixBit, "Prefix is as expected")
	}
}

// Unit test in order to check out if the method GetDifferentPrefix
// works on two string having different length
func TestGetDifferentPrefixWithDifferentLength(t *testing.T) {
	t.Logf("Test GetDifferentPrefixWithSameLength started! \n")
	var (
		a      = assert.New(t)
		s1, l1 = "fare", bitdata.GetLengthInBit("fare")
		s2, l2 = "stare", bitdata.GetLengthInBit("stare")
	)
	a.NotEqual(s1, s2, "Strings should be not equal")
	a.NotEqual(l1, l2, "Their length should not be equal")

	b1, e1 := bitdata.GetBitData(s1)
	b2, e2 := bitdata.GetBitData(s2)
	a.Nil(e1, "Error on converting first string")
	a.Nil(e2, "Error on converting second string")
	a.NotEqual(b1.Len, b2.Len, "Strings have not the same length")
	a.NotEqual(b1, b2, "Bitarrays are different")

	// The different prefix should be the reverse of "01110011 0111010"
	expectedPrefix := []bool{
		false, true, false, true, true, true, false,
		true, true, false, false, true, true, true, false,
	}
	receivedPrefix, err := b1.GetDifferentPrefix(b2)
	a.Nil(err, "Unexpected error")
	a.Equal(uint64(len(expectedPrefix)), receivedPrefix.Len, "Prefix should be of %d bits",
		uint64(len(expectedPrefix)))
	for i := uint64(0); i < receivedPrefix.Len; i++ {
		receivedPrefixBit, err := receivedPrefix.GetBit(uint64(i))
		a.Nil(err, "An error should not be occurred")
		a.Equal(expectedPrefix[i], receivedPrefixBit, "Prefix should be as expected")
	}
}

// Unit test in order to check out if the method bitToByte
// works on a string of 4 characters
func TestBitToByte(t *testing.T) {
	const s = "ciao"
	var (
		a    = assert.New(t)
		b, e = bitdata.GetBitData(s)
	)
	a.Nil(e)
	a.Equal(uint64(4*8), b.Len)
	bt, e := b.BitToByte()
	a.Nil(e)
	resultString := bytes.NewBuffer(bt).String()
	a.Equal(s, resultString)
}

// Unit test in order to check out if the method bitToString
// works on a string of 4 characters
func TestBitToString(t *testing.T) {
	const s = "ciao"
	var (
		a    = assert.New(t)
		b, e = bitdata.GetBitData(s)
	)
	a.Nil(e)
	a.Equal(uint64(4*8), b.Len)
	sCheck, e := b.BitToString()
	a.Nil(e)
	a.Equal(s, sCheck)
}

func TestSelect1(t *testing.T) {
	const (
		s       = "ciao" // 01100011011010010110000101101111
		sBitLen = 4 * 8
		check1  = uint64(0)
		check2  = uint64(3)
	)
	var (
		bd, errBd = bitdata.GetBitData(s)
		a         = assert.New(t)
	)
	if errBd != nil {
		a.Fail(errBd.Error())
	}
	selectOne1, errSelect1 := bd.Select1(uint64(1))
	if errSelect1 != nil {
		a.Fail(errSelect1.Error())
	}
	a.Equal(check1, selectOne1, "select1(1) should be 2")

	selectOne2, errSelect2 := bd.Select1(uint64(4))
	if errSelect2 != nil {
		a.Fail(errSelect2.Error())
	}
	a.Equal(check2, selectOne2, "select1(4) should be 7")

	select1Zero, errSelect1Zero := bd.Select1(uint64(0))
	a.Equal(uint64(0), select1Zero, "there should be an error so we have 0")
	a.NotNil(errSelect1Zero, "error message should not be nil")

	select133, errSelect133 := bd.Select1(uint64(33))
	a.Equal(uint64(0), select133, "there should be an error so we have 0")
	a.NotNil(errSelect133, "error message should not be nil")

}

func TestRank1(t *testing.T) {
	const (
		s       = "ciao" // 01100011011010010110000101101111
		sBitLen = 4 * 8
		check1  = uint64(2)
		check2  = uint64(4)
	)
	var (
		bd, errBd = bitdata.GetBitData(s)
		a         = assert.New(t)
	)
	if errBd != nil {
		a.Fail(errBd.Error())
	}

	rankOne1, errRankOne1 := bd.Rank1(uint64(2))
	a.Nil(errRankOne1, "there should be no error")
	a.Equal(check1, rankOne1, "rank1(2) should be 2")

	rankOne2, errRankOne2 := bd.Rank1(uint64(4))
	a.Nil(errRankOne2, "there should be no error")
	a.Equal(check2, rankOne2, "rank1(4) should be 4")

	rank1Zero, errRank1Zero := bd.Rank1(uint64(0))
	a.Equal(uint64(0), rank1Zero, "there should be an error so we have 0")
	a.NotNil(errRank1Zero, "error message should not be nil")

	rank133, errRank133 := bd.Select1(uint64(33))
	a.Equal(uint64(0), rank133, "there should be an error so we have 0")
	a.NotNil(errRank133, "error message should not be nil")

}
