package wordreader

import (
	"bufio"
	"os"
)

// WordReader contains a path to a file and the list of strings in that file
type WordReader struct {
	path    string
	Strings []string
}

// New returns a new WordReader
func New(path string) *WordReader {
	return &WordReader{path, []string{}}
}

// ReadLines reads all the lines in a file
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
