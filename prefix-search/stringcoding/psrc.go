package stringcoding

import (
	"fmt"
	bd "github.com/dariodip/prefix-search/prefix-search/bitdata"
	"github.com/golang-collections/go-datastructures/bitarray"
)

// PSRCBitDataSize contains the size of all the data structures for PSRC
type PSRCBitDataSize struct {
	StringsSize        uint64
	StartsSize         uint64
	LengthsSize        uint64
	IsUncompressedSize uint64
	PrefixOrSuffixSize uint64
}

// PSRC contains all the data structures to run PSRC algorithm
type PSRC struct {
	coding                     *Coding
	Epsilon                    float64
	c                          float64
	latestCompressedBitWritten uint64
	strings                    []string
	isUncompressed             *bd.BitData
	isStoredSuffix             *bd.BitData
}

// NewPSRC return an implementation of PSRC: a storage method
// based on RC (Rear Coding) that stores a string s in an
// uncompressed way if the latest c|s| bits do not contain
// an uncompressed string.
func NewPSRC(strings []string, epsilon float64) PSRC {
	stringsCount := uint64(len(strings))
	if epsilon <= float64(0) { // check if epsilon is valid
		panic("epsilon should be greater than 0")
	}
	c := 2.0 + 2.0/epsilon
	return PSRC{New(strings),
		epsilon,
		c, 0,
		strings,
		bd.New(bitarray.NewBitArray(stringsCount), stringsCount),
		bd.New(bitarray.NewBitArray(stringsCount), stringsCount)}
}

// Populate populates all the trie
func (psrc *PSRC) Populate() error {
	for i, s := range psrc.strings {
		if err := psrc.add(s, uint64(i)); err != nil {
			return err
		}
	}
	return nil
}

func (psrc *PSRC) add(s string, index uint64) error {
	coding := psrc.coding // extracting our coding data structure

	s = string("\x00") + s + string("\x00")
	bdS, errGbd := bd.GetBitData(s) // 1: convert string s to a bitdata bdS
	if errGbd != nil {
		return errGbd
	}
	var (
		differentSuffix *bd.BitData
		differentPrefix *bd.BitData
		stringToAdd     *bd.BitData
		storeSuffix     bool
		li              uint64
	)
	if coding.LastString != nil { // this is not the first string
		var errGds error
		differentSuffix, errGds = coding.LastString.GetDifferentSuffix(bdS) // 2: get different suffix
		if errGds != nil {
			return errGds
		}
		var errGdp error
		differentPrefix, errGdp = coding.LastString.GetDifferentPrefix(bdS) // 2a: get different prefix
		if errGdp != nil {
			return errGdp
		}
		if differentSuffix.Len > differentPrefix.Len {
			stringToAdd = differentPrefix
			storeSuffix = false
		} else {
			stringToAdd = differentSuffix
			storeSuffix = true
		}
	} else {
		stringToAdd = bdS // 2b: this is the first string so we cannot have different suffix
	}

	saveUncompressed := saveUncompressedPSRC(stringToAdd, bdS, psrc) // should our string be saved uncompressed?
	if saveUncompressed {                                            // we have to save our string uncompressed
		stringToAdd = bdS                                         // so the string to save is the full string
		psrc.latestCompressedBitWritten = uint64(0)               // compressed bit written is now 0
		if err := psrc.isUncompressed.SetBit(index); err != nil { // We can compress s
			return err
		}
	}
	errAppendBit := coding.Strings.AppendBits(stringToAdd) // 3: append string to Strings bitdata
	if errAppendBit != nil {
		panic(errAppendBit) // we don't know if the method has written in the structure
		// so we have to stop all the process and redo... sorry :(
	}

	// 4: li is the number of bit to remove on the prefix (risp. suffix) in the preceding string
	li = bdS.Len - stringToAdd.Len // our string - different suffix
	if coding.LastString != nil {
		errAppUL := coding.encodeEliasGamma(calcLen(li, coding.LastString.Len))
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
		psrc.latestCompressedBitWritten += stringToAdd.Len
	}
	if storeSuffix {
		psrc.isStoredSuffix.SetBit(index)
	} else {
		psrc.isStoredSuffix.ClearBit(index) // just for safety
	}
	return nil
}

func saveUncompressedPSRC(stringToAdd *bd.BitData, bdS *bd.BitData, psrc *PSRC) bool {
	return stringToAdd.Len == bdS.Len || float64(psrc.latestCompressedBitWritten) > psrc.c*float64(bdS.Len)
}

func (psrc *PSRC) getStringLength(i uint64) (uint64, error) {
	isUncomp, err := psrc.isUncompressed.GetBit(i)
	if err != nil {
		return uint64(0), err
	}
	if i == uint64(0) || isUncomp { // uncompressed string
		return psrc.getLengthInStrings(i)
	}
	lengthStringPI, err := psrc.getStringLength(i - 1) // length of the parent
	if err != nil {
		return uint64(0), err
	}
	li, err := psrc.coding.decodeIthEliasGamma(i) // li is the number of bits to remove in string(p(i)) in order to
	if err != nil {                               // obtain the prefix for string(i)
		return uint64(0), err
	}
	ni := lengthStringPI - li
	lengthI, err := psrc.getLengthInStrings(i) // lengthI is the number of bits stored in Strings
	if err != nil {
		return uint64(0), err
	} // our string length is the number of bits stored in String +
	return lengthI + ni, nil // number of bit saved by the coding
}

// Retrieval (u, l) returns the prefix of the string string(u) with length l.
// So the returned prefix ends up in the edge (p(u), u).
func (psrc *PSRC) Retrieval(u uint64, l uint64) (string, error) {
	l += 8
	var (
		stringBuffer *bd.BitData
	)
	isUncompressedStringU, errIsCompressed := psrc.isUncompressed.GetBit(u) // check if our string is compressed (we hope no)
	if errIsCompressed != nil {                                             // isUncompressed has gone wrong
		return "", errIsCompressed
	}
	if isUncompressedStringU { // our string is stored uncompressed
		ll, err := psrc.getLengthInStrings(u)
		if err != nil {
			return "", err
		} else { // no error
			if l > ll { // l is greater than our string
				l = ll // we can only return a string as big as our string
			}
		} // end else
		stringBuffer := bd.New(bitarray.NewBitArray(l), l)
		err = psrc.populateBuffer(stringBuffer, l, u, uint64(0), l) // we get the first l bits of that string
		if err != nil {
			panic(err)
		}
		return stringBuffer.BitToTrimmedString()
	} else { // our string is stored compressed
		// we'll do Select1(V, Rank1(V, u))
		v, err := psrc.isUncompressed.Rank1(u) // extract the number of 1s before u in isUncompressed
		if err != nil {                        // i.e. the number of uncompressed strings before u
			return "", err
		}
		vPosition, err := psrc.isUncompressed.Select1(v) // extract the position of the v-th string
		if err != nil {                                  // i.e. the first uncompressed string before u
			return "", err
		}
		vStarts, err := psrc.coding.Starts.Select1(vPosition + 1) // give me the position where the string v starts in Strings
		if err != nil {                                           // where v is the first uncompressed string before u
			return "", err
		}
		vNextStarts, err := psrc.coding.Starts.Select1(vPosition + 1 + 1) // give me the position of the string next to v
		if err != nil {                                                   // in order to extract the size of string(v)
			return "", err
		}
		lengthStringV := vNextStarts - vStarts // that's the length of string(v)
		stringBuffer = bd.New(bitarray.NewBitArray(lengthStringV), lengthStringV)
		err = psrc.populateBuffer(stringBuffer, lengthStringV, vPosition, 0, lengthStringV) // insert the first l bits of string(v) in the buffer
		if err != nil {
			panic(err)
		}
		for i := vPosition + 1; i <= u; i++ { // for each string i between v and u (we follow the path on the trie in dfs order)
			li, err := psrc.coding.decodeIthEliasGamma(i) // li is the number of bits to remove in string(p(i)) in order to
			if err != nil {                               // obtain the common string for string(i)
				return "", err
			}
			ni := lengthStringV - li // this is the length of the common bits between string(p(i))
			// and string(i)
			lengthI, err := psrc.getLengthInStrings(i) // length of the suffix of string(i) in Strings
			lengthStringV = lengthI + ni               // total length of string(i)
			if err != nil {
				return "", err
			}
			isStoredSuffix, err := psrc.isStoredSuffix.GetBit(u)
			if err != nil {
				return "", err
			}
			newBuffer := bd.New(bitarray.NewBitArray(lengthStringV), lengthStringV)
			if isStoredSuffix {
				// we didn't store LastString.Len - li but li
				sbLen := stringBuffer.Len
				for i := uint64(0); i < ni; i++ {
					bit, err := stringBuffer.GetBit(sbLen - 1 - i)
					if err != nil {
						return "", err
					}
					if bit {
						newBuffer.SetBit(newBuffer.Len - 1 - i)
					}
				}

				uPosition, err := psrc.coding.Starts.Select1(i + 1) // We need to now where the next string starts
				if err != nil {
					return "", err
				}
				lengthI, err := psrc.getLengthInStrings(i)
				if err != nil {
					return "", err
				}

				for i := uint64(0); i < lengthI; i++ {
					bit, err := psrc.coding.Strings.GetBit(uPosition + i)
					if err != nil {
						return "", err
					}
					if bit {
						newBuffer.SetBit(i)
					}
				}
				stringBuffer = newBuffer
			} else {

				var uPosition uint64
				if (u + 1) == uint64(len(psrc.strings)) {
					uPosition = psrc.coding.Strings.Len // u is the last string memorized!
				} else {
					var err error
					uPosition, err = psrc.coding.Starts.Select1(u + 1 + 1) // We need to now where the next string starts
					if err != nil {
						return "", err
					}
				}
				uPosition -= 1
				for i := uint64(0); i < lengthI; i++ {
					bit, err := psrc.coding.Strings.GetBit(uPosition - i)
					if err != nil {
						return "", err
					}
					if bit {
						err := newBuffer.SetBit(newBuffer.Len - 1 - i)
						if err != nil {
							return "", err
						}
					}
				}
				for i := uint64(0); i < ni; i++ {
					bit, err := stringBuffer.GetBit(i)
					if err != nil {
						return "", err
					}
					if bit {
						err := newBuffer.SetBit(i)
						if err != nil {
							return "", err
						}
					}
				}
				stringBuffer = newBuffer
			} // end else !isStoredSuffix
		} // end for
	} //end else !isUncompressedStringU

	if stringBuffer.Len < l { // Our string is too short
		return "", ErrTooShortString
	}
	firstLBits, err := stringBuffer.GetFirstLBits(l)
	if err != nil {
		return "", err
	}
	return firstLBits.BitToTrimmedString()
}

func (psrc *PSRC) getLengthInStrings(i uint64) (uint64, error) {
	startPositionI, err := psrc.coding.Starts.Select1(i + 1)
	if err != nil {
		return uint64(0), err
	}
	var startPositionSuccI uint64
	if (i + 1) == uint64(len(psrc.strings)) {
		startPositionSuccI = psrc.coding.Starts.Len // u is the last string memorized!
	} else {
		startPositionSuccI, err = psrc.coding.Starts.Select1(i + 1 + 1) // We need to now where the next string starts
		if err != nil {
			return uint64(0), err
		}
	}

	return startPositionSuccI - startPositionI, nil
}

func (psrc *PSRC) populateBuffer(stringBuffer *bd.BitData, l uint64, u uint64, ni uint64, vlen uint64) error {
	var (
		uPosition uint64
		maxIt     uint64
	)
	if (u + 1) == uint64(len(psrc.strings)) {
		uPosition = psrc.coding.Strings.Len // u is the last string memorized!
	} else {
		var err error
		uPosition, err = psrc.coding.Starts.Select1(u + 1 + 1) // We need to now where the next string starts
		if err != nil {
			return err
		}
	}
	if vlen < l {
		maxIt = vlen
	} else {
		maxIt = l - ni
	}
	uPosition = uPosition - 1            // The most significant bits are at the end
	for i := uint64(0); i < maxIt; i++ { // let's iterate for i = 0 up to l - 1 (l times)
		// We start from the most significant bits
		lastBit, lastBitErr := psrc.coding.Strings.GetBit(uPosition - i) // take the i-th bit of string(u)
		if lastBitErr != nil {                                           // getBit has gone wrong
			return lastBitErr
		}

		// The last ((l - 1) - ni) bits in stringBuffer are the number of most significant bit in common between
		// the two consecutive strings
		indexToUpdate := ((l - 1) - ni) - i
		if lastBit {
			stringBuffer.SetBit(indexToUpdate)
		} else {
			stringBuffer.ClearBit(indexToUpdate)
		}
	}
	return nil
}

func (psrc *PSRC) checkInterface() {
	checkFunc := func(search PrefixSearch) bool {
		return true
	}
	var sPs PrefixSearch
	psrcImpl := NewPSRC([]string{}, 1.0)
	sPs = &psrcImpl
	checkFunc(sPs)
}

func (psrc *PSRC) String() string {
	return fmt.Sprintf(`type:%T coding:%v, Epsilon:%v, c:%v, strings:%v, isUncompressed:%v, isStoredSuffix:%v`,
		psrc, psrc.coding, psrc.Epsilon, psrc.c, psrc.strings, psrc.isUncompressed, psrc.isStoredSuffix)
}

// FullPrefixSearch , given a prefix *prefix* returns all the strings that start with that prefix.
func (psrc *PSRC) FullPrefixSearch(prefix string) ([]string, error) {
	var (
		lenPrefix    = uint64(len(prefix) * 8) // |prefix|
		totalStrings = uint64(len(psrc.strings))
		stringBuffer = []string{}
		prefixBuffer = []uint64{}
	)

	for i := uint64(0); i < totalStrings; i++ {
		retrievalI, err := psrc.Retrieval(i, lenPrefix)
		if err != nil && err != ErrTooShortString { // If the string is too short, then we simply skip it
			return nil, err // if error was found
		}
		if retrievalI == prefix { // we found the first node having
			prefixBuffer = append(prefixBuffer, i)
		}
	}
	if len(prefixBuffer) == 0 {
		return []string{}, nil
	}

	for _, index := range prefixBuffer {
		stringLength, err := psrc.getStringLength(index)
		if err != nil {
			return []string{}, err
		}
		stringLength -= 8
		prefixedString, err := psrc.Retrieval(index, stringLength)
		if err != nil {
			return []string{}, err
		}
		stringBuffer = append(stringBuffer, prefixedString)
	}

	return stringBuffer, nil
}

// GetBitDataSize returns the size in bits of the BitData used to compress the strings
func (psrc *PSRC) GetBitDataSize() map[string]uint64 {
	sizes := make(map[string]uint64)
	sizes["StringSize"] = psrc.coding.Strings.Len
	sizes["StartsSize"] = psrc.coding.Starts.Len
	sizes["LenghtsSize"] = psrc.coding.Lengths.Len
	sizes["IsUncompressedSize"] = psrc.isUncompressed.Len
	sizes["PrefixOrSuffixSize"] = psrc.isStoredSuffix.Len

	return sizes
}
