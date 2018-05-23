package stringcoding

import (
	bd "github.com/dariodip/prefix-search/prefix-search/bitdata"
	"math"
)

// getEliasGammaLength computes the length of the Elias Gamma coding
// on the string set.
// TODO [optimization] GetLenthInBit called twice for each string
func getEliasGammaLength(strings []string) uint64 {
	count := uint64(0)
	for _, s := range strings {
		count += uint64(2*math.Floor(math.Log2(float64(bd.GetLengthInBit(s)))) + 1)
	}
	return count
}

// encodeEliasGamma appends Elias Gamma coding representation of the uint64 n
// to the Lengths bitdata.
// For more info check https://en.wikipedia.org/wiki/Elias_gamma_coding
func (c *Coding) encodeEliasGamma(n uint64) error {
	if c.Lengths == nil {
		return bd.ErrNotInitBitData
	}
	if n == uint64(0) { // a 0 length? sure?!?
		return bd.ErrZeroI
	}
	var (
		bigN = uint64(math.Floor(math.Log2(float64(n)))) // bigN is the first bit set to 1 in our n
	)
	for i := uint64(0); i < bigN; i++ { // write 0 bigN times
		if err := c.Lengths.AppendBit(false); err != nil {
			return err
		}
		c.NextLengthsIndex++
	}
	// once we wrote our |_log_2 (n) _| 0s, we have to convert our n to binary
	marker := uint64(1)                 // let's use a marker starting from 000...01
	for i := uint64(0); i < bigN; i++ { // let's shift our marker to the left, marking increasingly higher order bits
		marker = marker << 1 // until we reach bigN (the last bit set to 1 in our n)
	}
	for marker > 0 { // while marker is marking another valid bit
		if err := c.Lengths.AppendBit(marker&n != uint64(0)); err != nil { // let's add the i-th bit to our Lengths
			panic(err) // we got an error and must panic everything (we don't know how many bits have been written)
		}
		marker = marker >> 1 // moving marker on the lower order bit
		c.NextLengthsIndex++
	}
	return nil
}

// Given an index idx, returns the value of that index decoding
// the Elias Gamma coding
func (c *Coding) decodeIthEliasGamma(u uint64) (uint64, error) {
	if c.Lengths == nil {
		return uint64(0), bd.ErrNotInitBitData
	}
	if u >= c.Lengths.Len || u < uint64(0) {
		return uint64(0), bd.ErrIndexOutOfBound
	}

	currentIndex := uint64(0) // current index in the array
	currentNode := uint64(0)  // current node

	for u != currentNode { // while we are not in the desired node
		zeroCount, err := c.eliasGammaZeroCount(currentIndex) // count the total 0s in front of the coding
		if err != nil {
			return uint64(0), err
		}
		currentIndex += 2*zeroCount + 1 // advance the currentIndex of 2*zeroCount+1 positions
		currentNode++                   // let's move on the next node1
	}
	// now we found the desired node: currentNode and the index in which it starts
	zeroCount, err := c.eliasGammaZeroCount(currentIndex)
	if err != nil {
		return uint64(0), err
	}
	currentIndex += zeroCount                              // let's skip the first zeroCount bit, namely the 0s bit
	return c.extractNumFromBinary(currentIndex, zeroCount) // return our number as uint64
}

func (c *Coding) extractNumFromBinary(currentIndex uint64, zeroCount uint64) (uint64, error) {
	var (
		marker = uint64(1)    // 000...01
		n      = uint64(0)    // 000...00
		index  = currentIndex // index of Lengths in which we start
	)
	for i := uint64(0); i < zeroCount; i++ { // marker has 1 only in the position zeroCount
		marker = marker << 1
	}

	for marker > 0 { // marker has still a bit set to 1
		bit, err := c.Lengths.GetBit(index) // find the bit in Lengths bitdata
		index++
		if err != nil {
			return uint64(0), err // ops!
		}
		if bit { // our bit was set to 1
			n = n | marker // so let's mark our position in n
		}
		marker = marker >> 1 // marker now marks the lower bit
	}
	return n, nil
}

func (c *Coding) eliasGammaZeroCount(idx uint64) (uint64, error) {
	return c.eliasGammaZeroCountLoop(idx, 0)
}

// eliasGammaZerCountLoop uses stack in a recursive fashion in order to
// cound the number of 0s in an Elias gamma coding starting from the index idx
func (c *Coding) eliasGammaZeroCountLoop(idx uint64, zeroCount uint64) (uint64, error) {
	bit, err := c.Lengths.GetBit(idx)
	if err != nil {
		return uint64(0), err
	}
	if !bit {
		return c.eliasGammaZeroCountLoop(idx+1, zeroCount+1)
	}
	return zeroCount, nil
}
