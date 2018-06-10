package stringcoding

import "errors"

var (
	// ErrTooShortString is returned when you are trying to access given an index that isn't defined
	ErrTooShortString = errors.New("the string is too short to contain a prefix of that length")
)
