// stringcoding package provides an easy way to deal
// with strings in a bit-to-bit fashion
package stringcoding

import (
	"errors"
	"fmt"
	"github.com/golang-collections/go-datastructures/bitarray"
	bd "prefix-search/prefix-search/bitdata"
)

type Coding struct {
	// Strings consists of all the concatenated bit sequences
	// corresponding to the suffixes L[i] of S's strings.
	Strings *bd.BitData
	// Starts consists of a sequence of bits in which each bit
	// set to 1 marks the first bit of each of those suffixes
	// in the aforementioned array (Strings).
	Starts *bd.BitData
	// Lengths encodes in unary the length of the shared prefixes
	// between consecutive strings.
	Lengths *bd.BitData
	// LastString contains the last processed string as a sequence
	// of bit. It can be used for a more efficient processing of the
	// current string to deal with.
	LastString *bd.BitData
	// NextIndex marks the last index in Strings (risp. Starts) arrays.
	NextIndex uint64
	// NextLengthsIndex marks the last index in the Lengths array.
	NextLengthsIndex uint64
	// Function that computes the value to insert in the Lengths array.
	LengthCalcFunction func(uint64, uint64) uint64
}

// Create creates and returns a new Coding structure inserting the strings
// that are in the array of strings.
func New(strings []string, lenCalc func(uint64, uint64) uint64) *Coding {
	maxCapacity := bd.GetTotalBitCount(strings)
	maxLengthCapacity := maxCapacity + uint64(len(strings)-1)
	fc := Coding{
		Strings:            bd.New(bitarray.NewBitArray(maxCapacity), 0),
		Starts:             bd.New(bitarray.NewBitArray(maxCapacity), 0),
		Lengths:            bd.New(bitarray.NewBitArray(maxLengthCapacity), 1),
		NextLengthsIndex:   uint64(1),
		LengthCalcFunction: lenCalc,
	}
	// TODO insert
	return &fc
}

// add adds the string s to the structure
func (c *Coding) add(s string) error {
	bdS, errGbd := bd.GetBitData(s) // 1: convert string s to a bitdata bdS
	if errGbd != nil {
		return errGbd
	}

	var differentSuffix *bd.BitData
	if c.LastString != nil {
		var errGds error
		differentSuffix, errGds = c.LastString.GetDifferentSuffix(bdS)
		if errGds != nil {
			return errGds
		}
	} else {
		differentSuffix = bdS
	}

	// 2: get different suffix
	errAppendBit := c.Strings.AppendBits(differentSuffix) // 3: append string to Strings bitdata
	if errAppendBit != nil {
		panic(errAppendBit) // we don't know if the method has written in the structure
		// so we have to stop all the process and redo... sorry :(
	}

	// 4: append different suffix' length to Lengths
	prefixLen := bd.GetLengthInBit(s) - differentSuffix.Len
	errAppUL := c.addUnaryLength(c.LengthCalcFunction(prefixLen, bd.GetLengthInBit(s)))
	if errAppUL != nil { // as above...
		panic(errAppUL)
	}
	errSetSWO := c.setStartsWithOffset(differentSuffix) // 5: set the bit of the next string in the Starts array
	if errSetSWO != nil {
		panic(errSetSWO)
	}
	c.LastString = bdS // 6: update
	return nil
}

// setStartsWithOffset sets the bit in the Starts bitdata in order
// to state where the suffix in Strings starts.
func (c *Coding) setStartsWithOffset(differentSuffix *bd.BitData) error {
	if differentSuffix.Len == 0 {
		return nil // nothing to do here
	}

	c.Starts.AppendBit(true)

	c.Starts.Len += differentSuffix.Len - 1

	return nil
}

// addUnaryLength appends unary representation of the uint64 n
// to the Lengths bitdata.
func (c *Coding) addUnaryLength(n uint64) error {
	if c.Lengths == nil {
		return errors.New("error in trying to add on a non initialized BitData")
	}
	for i := uint64(0); i < n; i++ {
		if err := c.Lengths.AppendBit(true); err != nil {
			return err
		}
		c.NextLengthsIndex++
	}
	if err := c.Lengths.AppendBit(false); err != nil {
		return err
	}
	c.NextLengthsIndex++
	return nil
}

// Given an index, returns the idx-th value of the unary array
func (c *Coding) unaryToInt(idx uint64) (uint64, error) {
	if c.Lengths == nil {
		return uint64(0), errors.New("error in trying to add on a non initialized BitData")
	}
	if bit, err := c.Lengths.GetBit(idx); err == nil {
		if bit && idx != uint64(0) {
			return uint64(0), errors.New("index should point to a 0")
		}
	} else {
		return uint64(0), err
	}

	var val uint64
	current := idx
	for {
		current++
		if bit, err := c.Lengths.GetBit(current); err == nil {
			if bit {
				val++
			} else {
				break
			}
		} else {
			return uint64(0), err
		}
	}
	return val, nil
}

func (c *Coding) String() string {
	return fmt.Sprintf("type:%T, Strings: %v, Starts:%v, Lengths:%v, LastString:%v, NextIndex:%v, NextLengthsIndex:%v, LengthCalcFunction:%T",
		c, c.Strings, c.Starts, c.Lengths, c.LastString, c.NextIndex, c.NextLengthsIndex, c.LengthCalcFunction)
}
