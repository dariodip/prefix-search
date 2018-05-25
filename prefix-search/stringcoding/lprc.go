package stringcoding

import (
	bd "github.com/dariodip/prefix-search/prefix-search/bitdata"
	"github.com/golang-collections/go-datastructures/bitarray"
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
	} else {
		if err := lprc.isCompressed.SetBit(index); err != nil { // We can compress s
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
	if coding.LastString != nil {
		errAppUL := coding.encodeEliasGamma(calcLen(prefixLen, coding.LastString.Len))
		if errAppUL != nil { // as above...
			panic(errAppUL)
		}
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

// Retrieval(u, l) returns the prefix of the string string(u) with length l.
// So the returned prefix ends up in the edge (p(u), u).
func (lprc *LPRC) Retrieval(u uint64, l uint64) (string, error) {
	var (
		stringBuffer         = bd.New(bitarray.NewBitArray(l), 0) // let's create a buffer in order to store our prefix
		uPosition, errSelect = lprc.coding.Starts.Select1(u)      // uPosition is the position of the u-th string
	)
	if errSelect != nil { // select has gone wrong
		return "", errSelect
	}
	isCompressedStringU, errIsCompressed := lprc.isCompressed.GetBit(u) // check if our string is compressed (we hope no)
	if errIsCompressed != nil {                                         // isCompressed has gone wrong
		return "", errIsCompressed
	}
	if !isCompressedStringU { // our string is stored uncompressed
		for i := uint64(0); i < l; i++ { // let's iterate for i = 0 up to l - 1 (l times)
			lastBit, lastBitErr := lprc.coding.Strings.GetBit(uPosition + i) // take the i-th bit of string(u)
			if lastBitErr != nil {                                           // getBit has gone wrong
				return "", lastBitErr
			}
			stringBuffer.AppendBit(lastBit) // populate our buffer
		}
		return stringBuffer.BitToString() // return our buffer encoded as a string
	} else {
		panic("to implement") // TODO implement compressed case
	}
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
