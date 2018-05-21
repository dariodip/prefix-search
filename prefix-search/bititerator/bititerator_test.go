package bititerator

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// A unit test in order to check out if the operations
// on BitIterator work in the correct way
func TestBitIterator(t *testing.T) {
	a := assert.New(t)
	const s = "ciao" // base variable on which iterate
	t.Logf("Starting test on BitIterator using a simple string: %q", s)
	check := []byte{0, 0, 0, 0}
	t.Logf("Initially check arrays is: %s", string(check[:]))
	// Check equality
	a.NotEqual(check, []byte(s), "They should be not equal")
	t.Log("Both s and check are not equal")
	t.Log("Let's create a new bit iterator")
	bitIt := NewStringToBitIterator(s) // bit iterator created
	var currentBit uint                // current bit in the byte
	currentByte := len(check) - 1      // current byte in the slice
	counter := 0
	t.Log("Let's retrieve our bits")
	for bitIt.HasNext() {
		counter++
		b, err := bitIt.Next()
		a.Nil(err, "Error")
		if b { // bit set to 1, we should edit our byte
			check[currentByte] = check[currentByte] | 1<<currentBit
		}
		if currentBit == 7 { // last bit, let's switch to the preceding byte
			currentByte--
		}
		currentBit = (currentBit + 1) % 8
	}
	t.Log("Let's convert our bytes to string")
	checkS := bytes.NewBuffer(check).String()
	a.Equal(s, checkS, s+" equals to "+checkS)
	a.Equal(4*8, counter)
}
