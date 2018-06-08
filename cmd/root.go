package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	VERSION         string
	inputFile       string
	inputPrefixFile string
	outputFile      string
	algorithm       string
	epsilon         float64
	verbose         bool
	LPRCconst       = "lprc"
	PSRCconst       = "psrc"
)

type ResultRow struct {
	Prefix            string
	PrefixedWordCount int
	SearchTime        float64
}

type Result struct {
	InitTime             float64
	Epsilon              float64
	StructureSize        map[string]uint64
	UncompressedDataSize uint64
	PrefixResult         []ResultRow
	TotalSearchTime      float64
}

func (res *Result) addResultRow(prefix string, wordCount int, searchTime time.Duration) {
	res.PrefixResult = append(res.PrefixResult,
		ResultRow{prefix, wordCount, toMilliseconds(searchTime)})
}

var rootCmd = &cobra.Command{
	Use:   "prefix-search",
	Short: "Prefix Search: A tool that implements Prefix Search algorithm based on a cache oblivious string b-tree",
	Long: `Prefix Search: A tool that implements Prefix Search algorithm based on a cache oblivious string b-tree.  
It is an implementation of the paper Compressed Cache-Oblivious String B-tree of Paolo Ferragina and Rossano Venturini.
We developed the proposed algorithm (LPRC) and a new one (PSRC), giving you the ability to deal with online dictionaries
of strings in an unspecified order.
Check out our GitHub repository for more info: https://github.com/dariodip/prefix-search .`,
	Run: run,
}

func Execute(version string) {

	VERSION = version

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}

// Returns a function used to print benchmark update based on the verbose flag
func updateResultTemplate(verbose bool, stringsCount int) func(string) {
	if verbose {
		return func(result string) {
			fmt.Printf(result)
		}
	} else {
		idx := 1
		return func(result string) {
			fmt.Printf("\rQuering prefix %d on %d", idx, stringsCount)
			if idx == stringsCount {
				fmt.Println()
			} else {
				idx++
			}
		}
	}
}

// Extracts the file name (without extension) from a path
func getFileName(path string) string {
	file := filepath.Base(path)
	return strings.Split(file, ".")[0]
}

// Saves to a file all the result memorized in res
func saveToFile(res *Result, filename string) {
	fp, err := os.Create(filename)
	defer fp.Close()

	if err != nil {
		fmt.Printf("Cannot open file %s. %s\n", filename, err)
	} else {
		encodedResults, err := json.Marshal(res)
		if err != nil {
			fmt.Printf("Cannot save file %s. %s\n", filename, err)
		} else {
			fp.Write(encodedResults)
		}
	}
}

// Returns the size, as bits, of a list of strings
func getBitSize(strings []string) uint64 {
	bitSize := uint64(0)
	for _, s := range strings {
		bitSize += uint64(len([]byte(s)) * 8)
	}

	return bitSize
}

// It converts a time.Duration value in milliseconds
func toMilliseconds(t time.Duration) float64 {
	return float64(t) / float64(time.Millisecond)
}
