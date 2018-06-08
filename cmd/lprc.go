package cmd

import (
	"fmt"
	"github.com/dariodip/prefix-search/prefix-search/stringcoding"
	"github.com/dariodip/prefix-search/word-reader"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type LPRCResult struct {
	InitTime             float64
	Epsilon              float64
	StructureSize        stringcoding.LPRCBitDataSize
	UncompressedDataSize uint64
	PrefixResult         []ResultRow
	TotalSearchTime      float64
}

// lprcCmd represents the lprc command
var lprcCmd = &cobra.Command{
	Use:   "lprc",
	Short: "Run lprc algorithm",
	Long: `lprc (Locality Preserving Rear Coding) is an algorithm designed by Paolo Ferragina and Rossano Venturini
in their paper "Compressed Cache-Oblivious String B-Tree". 

Our implementation takes in input: 
	- a file containing all the worlds to add to the dictionary (-i);
	- a file containing all the prefixes to search on the built dictionary (-p).
	- the epsilon to use in order to build our structure

All the results will be saved into a json file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		lprcBenchmark()
	},
}

func init() {

	rootCmd.AddCommand(lprcCmd)
	lprcCmd.Flags().StringVarP(&inputFile, "input_file", "i", "", "Input file containing"+
		" all the word to build up the dictionary.")
	lprcCmd.MarkFlagRequired("input_file")
	lprcCmd.MarkFlagFilename("input_file")

	lprcCmd.Flags().StringVarP(&inputPrefixFile, "input_p_file", "p", "", "Input"+
		" file containing all the prefix to search on the dictionary.")
	lprcCmd.MarkFlagRequired("input_p_file")
	lprcCmd.MarkFlagFilename("input_p_file")

	lprcCmd.Flags().Float64VarP(&epsilon, "epsilon", "e", 0, "Epsilon is the parameter"+
		"given to the algorithm in order to decide how many bits compress in the trie.")
	lprcCmd.MarkFlagRequired("epsilon")

	lprcCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Detailed Output ")

	lprcCmd.Flags().StringVarP(&outputFile, "output_file", "o", "", "Output file"+
		" containing the final output of lprc, with information about the memory usage and the time elapsed.\n"+
		"Default <word filename>-<prefix file name>-<epsilon>.json")
	lprcCmd.MarkFlagFilename("output_file")
}

func lprcBenchmark() {
	// load words
	wr := word_reader.New(inputFile)
	wr.ReadLines()

	// load prefix
	wrp := word_reader.New(inputPrefixFile)
	wrp.ReadLines()

	lprcImpl, initTime, err := initLPRC(wr.Strings, epsilon)
	if err != nil {
		fmt.Printf("Unable to complete the benchmark: %s\n", err)
		os.Exit(-1)
	}

	bdSize := lprcImpl.GetBitDataSize()
	totalBitSize := bdSize["StringsSize"] + bdSize["StartsSize"] + bdSize["LengthsSize"] + bdSize["IsUncompressedSize"]
	fmt.Printf("Initialization time:   %v\n", initTime)
	fmt.Printf("Size of the structure: %d bits\n", totalBitSize)
	fmt.Println()

	if outputFile == "" { // no output file specified
		outputFile = fmt.Sprintf("%s-%s-%.2f.json", getFileName(inputFile), getFileName(inputPrefixFile),
			lprcImpl.Epsilon)
	}
	finalResults := &Result{
		InitTime:             toMilliseconds(initTime),
		Epsilon:              lprcImpl.Epsilon,
		StructureSize:        bdSize,
		UncompressedDataSize: getBitSize(wr.Strings),
	}
	defer saveToFile(finalResults, outputFile)

	var searchTime time.Time
	var elapsedTime time.Duration
	totalSearchTime := time.Duration(0)
	updateResult := updateResultTemplate(verbose, len(wrp.Strings))
	for _, prefix := range wrp.Strings {
		searchTime = time.Now()
		result, err := lprcImpl.FullPrefixSearch(prefix)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		elapsedTime = time.Since(searchTime)

		updateResult(fmt.Sprintf("Full-Prefix-Search for prefix %s -> strings found: %d, time elapsed: %v\n",
			prefix, len(result), elapsedTime))

		totalSearchTime += elapsedTime
		finalResults.addResultRow(prefix, len(result), elapsedTime)
	}

	fmt.Println()
	fmt.Printf("Full-Prefix-Search total elapsed time: %v\n", totalSearchTime)

	finalResults.TotalSearchTime = toMilliseconds(totalSearchTime)
}

func initLPRC(strings []string, epsilon float64) (*stringcoding.LPRC, time.Duration, error) {
	startTime := time.Now()
	lprcImpl := stringcoding.NewLPRC(strings, epsilon)
	if err := lprcImpl.Populate(); err != nil {
		return nil, time.Duration(0), err
	}

	return &lprcImpl, time.Since(startTime), nil
}
