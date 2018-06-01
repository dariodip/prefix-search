package bitdata

import (
	"errors"
	"fmt"
)

var (
	ErrIndexOutOfBound = errors.New("index out of bound")
	ErrNotInitBitData  = errors.New("cannot append to a non initialized BitData")
	ErrInvalidString   = errors.New("bitdata should be a valid string")
	ErrLessThanIOnes   = errors.New("there are less than i 1s in the array")
	ErrInvalidI        = errors.New("i should not be greater than the length of the array")
	ErrZeroI           = errors.New("i should be greater than 0")
)

type ErrInvalidPosition struct {
	index uint64
}

func (e *ErrInvalidPosition) Error() string {
	return fmt.Sprintf("cannot access bitarray in position: %d", e.index)
}
