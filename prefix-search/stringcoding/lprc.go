package stringcoding

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	bd "prefix-search/prefix-search/bitdata"
)

type LPRC struct {
	coding                     *Coding
	Epsilon                    float64
	c                          float64
	latestCompressedBitWritten uint64
	strings                    []string
	isCompressed               *bd.BitData
}

// LPRC (Locality Preserving Rear Coding) is a storage method
// based on RC (Rear Coding) that stores a string s in an
// uncompressed way if the latest c|s| bits do not contain
// an uncompressed string.
func NewLPRC(strings []string, epsilon float64) LPRC {
	stringsCount := uint64(len(strings))
	c := 2.0 + 2.0/epsilon
	return LPRC{New(strings),
		epsilon,
		c, 0,
		strings,
		bd.New(bitarray.NewBitArray(stringsCount), stringsCount)}
}

func calcLen(prefixLen, stringLen uint64) uint64 {
	return stringLen - prefixLen
}

// add adds the string s to the structure
func (lprc *LPRC) add(s string, index uint64) error {
	coding := lprc.coding // extracting our coding data structure

	bdS, errGbd := bd.GetBitData(s) // 1: convert string s to a bitdata bdS
	if errGbd != nil {
		return errGbd
	}
	var stringToAdd *bd.BitData
	if coding.LastString != nil { // this is not the first string
		var errGds error
		stringToAdd, errGds = coding.LastString.GetDifferentSuffix(bdS) // 2: get different suffix
		if errGds != nil {
			return errGds
		}
	} else {
		stringToAdd = bdS // 2b: this is the first string so we cannot have different suffix
	}

	saveUncompressed := saveUncompressed(stringToAdd, bdS, lprc) // should our string be saved uncompressed?
	if saveUncompressed {                                        // we have to save our string uncompressed
		stringToAdd = bdS                           // so the string to save is the full string
		lprc.latestCompressedBitWritten = uint64(0) // compressed bit written is now 0
		if err := lprc.isCompressed.SetBit(index); err != nil {
			return err
		}
	}
	errAppendBit := coding.Strings.AppendBits(stringToAdd) // 3: append string to Strings bitdata
	if errAppendBit != nil {
		panic(errAppendBit) // we don't know if the method has written in the structure
		// so we have to stop all the process and redo... sorry :(
	}

	// 4: append different suffix' length to Lengths
	prefixLen := bdS.Len - stringToAdd.Len // get suffix length
	errAppUL := coding.addUnaryLength(calcLen(prefixLen, bdS.Len))
	if errAppUL != nil { // as above...
		panic(errAppUL)
	}
	errSetSWO := coding.setStartsWithOffset(stringToAdd) // 5: set the bit of the next string in the Starts array
	if errSetSWO != nil {
		panic(errSetSWO)
	}
	coding.LastString = bdS // 6: update last string
	if !saveUncompressed {  // 7: if the string was saved compressed we have to update latestCompressedBitWritten counter
		lprc.latestCompressedBitWritten += stringToAdd.Len
	}
	return nil
}

func saveUncompressed(stringToAdd *bd.BitData, bdS *bd.BitData, lprc *LPRC) bool {
	return stringToAdd.Len == bdS.Len || float64(lprc.latestCompressedBitWritten) > lprc.c*float64(bdS.Len)
}

func (lprc *LPRC) run() error {
	for i, s := range lprc.strings {
		if err := lprc.add(s, uint64(i)); err != nil {
			return err
		}
	}
	return nil
}
