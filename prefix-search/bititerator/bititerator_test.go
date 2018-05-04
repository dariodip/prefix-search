package bititerator

import (
	"testing"
	"bytes"
)

func TestBitIterator(t *testing.T) {
	const s = "ciao"  // base variable on which iterate
	t.Logf("Starting test on BitIterator using a simple string: %q", s)
	check := []byte{0, 0, 0, 0}
	t.Logf("Initially check arrays is: %s", string(check[:]))
	if bytes.Equal(check, []byte(s)) {  // they should not be equal
		t.FailNow()
	}
	t.Log("Both s and check are not equal")
	t.Log("Let's create a new bit iterator")
	bitIt := NewStringToBitIterator(s) // bit iterator created

	var currentBit uint  // current bit in the byte
	currentByte := len(check) - 1  // current byte in the slice

	t.Log("Let's retrieve our bits")
	for bitIt.HasNext() {
		b, err := bitIt.Next()
		if err != nil {  // we got an error, the test should not pass :D
			t.Errorf("Error in the iterator: %s", err.Error())
			t.FailNow()
		}
		if b {  // bit set to 1, we should edit our byte
			check[currentByte] = check[currentByte] | 1 << currentBit
		}
		if currentBit == 7 { // last bit, let's switch to the preceding byte
			currentByte--
		}
		currentBit = (currentBit + 1) % 8
	}
	t.Log("Let's convert our bytes to string")
	checkS := bytes.NewBuffer(check).String()
	if checkS == s {
		t.Logf("%q equals to %q. Test passed!", checkS, s)
	} else {
		t.Errorf("%q not equals to %q. Test failed!", checkS, s)
		t.FailNow()
	}
}