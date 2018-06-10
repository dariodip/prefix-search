package stringcoding

// PrefixSearch interface contains all the methods in order to run both LPRC and PSRC
type PrefixSearch interface {
	Populate() error
	add(string, uint64) error
	Retrieval(uint64, uint64) (string, error)
	FullPrefixSearch(prefix string) ([]string, error)
	GetBitDataSize() map[string]uint64
	checkInterface()
}
