package stringcoding

import "errors"

var (
	ErrInvalidIndex = errors.New("index should point to a 0")
	ErrTooShortString = errors.New("The string is too short to contain a prefix of that length")
)
