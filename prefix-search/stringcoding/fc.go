package stringcoding

type FC struct {
	c *Coding
}

func NewFC(strings []string) FC {
	return FC{New(strings, func (prefixLen, stringLen uint64) uint64 {
		return prefixLen
	})}
}
