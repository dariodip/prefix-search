package bitdata

import (
	"errors"
	"fmt"
)

var (
	// ErrIndexOutOfBound represents an index out of bound error
	ErrIndexOutOfBound = errors.New("index out of bound")
	// ErrNotInitBitData tells that the BitData on which you are operating is not initialized
	ErrNotInitBitData = errors.New("cannot append to a non initialized BitData")
	// ErrInvalidString indicates that the string on which you are working is not valid
	ErrInvalidString = errors.New("bitdata should be a valid string")
	// ErrLessThanIOnes is returned by Select1 or Rank1 when you are trying to access to a "1" that does not exist
	ErrLessThanIOnes = errors.New("there are less than i 1s in the array")
	// ErrInvalidI is returned when i is greater than the length of the array
	ErrInvalidI = errors.New("i should not be greater than the length of the array")
	// ErrZeroI is returned when you are passing a value of i equal to 0
	ErrZeroI = errors.New("i should be greater than 0")
)

// ErrInvalidPosition is returned when you are trying to access to an invalid position
type ErrInvalidPosition struct {
	index uint64
}

func (e *ErrInvalidPosition) Error() string {
	return fmt.Sprintf("cannot access bitarray in position: %d", e.index)
}
