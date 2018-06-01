package word_reader

import (
	"bufio"
	"os"
)

type WordReader struct {
	path    string
	Strings []string
}

func New(path string) *WordReader {
	return &WordReader{path, []string{}}
}

func (wr *WordReader) ReadLines() (int, error) {

	file, err := os.Open(wr.path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wr.Strings = append(wr.Strings, scanner.Text())
	}
	return len(wr.Strings), scanner.Err()

}
