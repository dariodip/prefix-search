package stringcoding

type RC struct {
	c *Coding
}

func NewRC(strings []string) RC {
	return RC{New(strings, func (prefixLen, stringLen uint64) uint64 {
		return stringLen - prefixLen
	})}
}
