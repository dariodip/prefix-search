package word_reader

import (
	"bufio"
	"os"
)

type WordReader struct {
	path    string
	strings []string
}

func New(path string) *WordReader {
	return &WordReader{path, []string{}}
}

func (wr *WordReader) readLines() (int, error) {

	file, err := os.Open(wr.path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wr.strings = append(wr.strings, scanner.Text())
	}
	return len(wr.strings), scanner.Err()

}
