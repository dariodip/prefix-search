package stringcoding

type PrefixSearch interface {
	Populate() error
	add(string, uint64) error
	Retrieval(uint64, uint64) (string, error)
	FullPrefixSearch(prefix string) ([]string, error)
	checkInterface()
	GetBitDataSize() interface{}
}
