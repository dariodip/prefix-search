package cmd

import (
	"fmt"

	"github.com/dariodip/prefix-search/prefix-search/stringcoding"
	"github.com/dariodip/prefix-search/word-reader"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// psrcCmd represents the psrc command
var psrcCmd = &cobra.Command{
	Use:   "psrc",
	Short: "Run psrc algorithm",
	Long: `psrc (Prefix-Suffix Rear Coding) is an algorithm designed by Mattia Tomeo and Dario Di Pasquale, 
inspired by the paper "Compressed Cache-Oblivious String B-Tree". 

Our implementation takes in input two files: 
	- a file containing all the worlds to add to the dictionary (-i);
	- a file containing all the prefixes to search on the built dictionary (-p).

All the results will be saved into a json file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		psrcBenchmark()
	},
}

func init() {
	rootCmd.AddCommand(psrcCmd)

	psrcCmd.Flags().StringVarP(&inputFile, "input_file", "i", "", "Input file containing"+
		" all the word to build up the dictionary")
	psrcCmd.MarkFlagRequired("input_file")
	psrcCmd.MarkFlagFilename("input_file")

	psrcCmd.Flags().StringVarP(&inputPrefixFile, "input_p_file", "p", "", "Input"+
		" file containing all the prefix to search on the dictionary")
	psrcCmd.MarkFlagRequired("input_p_file")
	psrcCmd.MarkFlagFilename("input_p_file")

	psrcCmd.Flags().Float64VarP(&epsilon, "epsilon", "e", 0, "Epsilon is the parameter"+
		"given to the algorithm in order to decide how many bits compress in the trie.")
	psrcCmd.MarkFlagRequired("epsilon")

	psrcCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Detailed Output ")

	psrcCmd.Flags().StringVarP(&outputFile, "output_file", "o", "", "Output file"+
		" containing the final output of lprc, with information about the memory usage and the time elapsed.\n"+
		"Default <word filename>-<prefix file name>-<epsilon>.json")
	psrcCmd.MarkFlagFilename("output_file")
}

func psrcBenchmark() {
	// load words
	wr := word_reader.New(inputFile)
	wr.ReadLines()

	// load prefix
	wrp := word_reader.New(inputPrefixFile)
	wrp.ReadLines()

	psrcImpl, initTime, err := initPSRC(wr.Strings, epsilon)
	if err != nil {
		fmt.Printf("Unable to complete the benchmark: %s\n", err)
		os.Exit(-1)
	}

	bdSize := psrcImpl.GetBitDataSize()
	totalBitSize := bdSize["StringsSize"] + bdSize["StartsSize"] + bdSize["LengthsSize"] + bdSize["IsUncompressedSize"] + bdSize["PrefixOrSuffixSize"]
	fmt.Printf("Initialization time:   %v\n", initTime)
	fmt.Printf("Size of the structure: %d bits\n", totalBitSize)
	fmt.Println()

	if outputFile == "" { // no output file specified
		outputFile = fmt.Sprintf("%s-%s-%.2f.json", getFileName(inputFile), getFileName(inputPrefixFile),
			psrcImpl.Epsilon)
	}
	finalResults := &Result{
		InitTime:             toMilliseconds(initTime),
		Epsilon:              psrcImpl.Epsilon,
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
		result, err := psrcImpl.FullPrefixSearch(prefix)
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

func initPSRC(strings []string, epsilon float64) (*stringcoding.PSRC, time.Duration, error) {
	startTime := time.Now()
	psrcImpl := stringcoding.NewPSRC(strings, epsilon)
	if err := psrcImpl.Populate(); err != nil {
		return nil, time.Duration(0), err
	}

	return &psrcImpl, time.Since(startTime), nil
}
