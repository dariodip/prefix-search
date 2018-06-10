package cmd

import (
	"fmt"

	"github.com/dariodip/prefix-search/prefix-search/stringcoding"
	"github.com/dariodip/prefix-search/word-reader"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

// fullbenchmarkCmd represents the fullbenchmark command
var fullbenchmarkCmd = &cobra.Command{
	Use:   "fullbenchmark",
	Short: "Run a complete benchmark",
	Long: `This command gives you the ability to run a sophisticated benchmarking test.
You can select the file to open as dataset, the file to open as prefix, the lower value
of epsilon, the higher value of epsilon and the step with which increase the value of it.

The test gives you only a JSON file containing all the results of the test.`,
	Run: func(cmd *cobra.Command, args []string) {
		fullBenchmark()
	},
}

func init() {
	rootCmd.AddCommand(fullbenchmarkCmd)

	fullbenchmarkCmd.Flags().StringVarP(&inputFile, "input_file", "i", "", "Input file containing"+
		" all the word to build up the dictionary.")
	fullbenchmarkCmd.MarkFlagRequired("input_file")
	fullbenchmarkCmd.MarkFlagFilename("input_file")

	fullbenchmarkCmd.Flags().StringVarP(&inputPrefixFile, "input_p_file", "p", "", "Input"+
		" file containing all the prefix to search on the dictionary.")
	fullbenchmarkCmd.MarkFlagRequired("input_p_file")
	fullbenchmarkCmd.MarkFlagFilename("input_p_file")

	fullbenchmarkCmd.Flags().StringArrayVarP(&epsilonList, "epsilon_list", "l", []string{}, "List"+
		" of epsilon value with which test the algorithm.")
	fullbenchmarkCmd.MarkFlagRequired("epsilon_list")

	fullbenchmarkCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Detailed Output ")

	fullbenchmarkCmd.Flags().StringVarP(&algorithm, "algorithm", "a", "lprc", "Algorithm"+
		"to use")
	fullbenchmarkCmd.MarkFlagRequired("algorithm")

	fullbenchmarkCmd.Flags().StringVarP(&outputFile, "output_file", "o", "", "Output file"+
		" containing the final output of lprc, with information about the memory usage and the time elapsed.\n"+
		"Default <algorithm>-<word filename>-<prefix file name>-<min_epsilon>-<max_epsilon>.json")
	fullbenchmarkCmd.MarkFlagFilename("output_file")
}

func fullBenchmark() {

	// load words
	wr := word_reader.New(inputFile)
	wr.ReadLines()

	// load prefix
	wrp := word_reader.New(inputPrefixFile)
	wrp.ReadLines()

	if outputFile == "" { // no output file specified
		outputFile = fmt.Sprintf("%s-%s-%s-(%d).json", algorithm, getFileName(inputFile),
			getFileName(inputPrefixFile), time.Now().Unix())
	}

	var (
		allResults       = []*Result{}
		epsilonListFloat []float64
	)
	for _, e := range epsilonList {
		eFloat, err := strconv.ParseFloat(e, 64)
		if err != nil {
			fmt.Println("invalid epsilon list")
			os.Exit(1)
		}
		epsilonListFloat = append(epsilonListFloat, eFloat)
	}

	for _, eps := range epsilonListFloat {
		var impl stringcoding.PrefixSearch
		var initTime time.Duration

		if algorithm == LPRCconst {
			lprcImpl, iTime, err := initLPRC(wr.Strings, eps)
			if err != nil {
				fmt.Printf("Unable to complete the benchmark: %s\n", err)
				os.Exit(-1)
			}
			impl = lprcImpl
			initTime = iTime
		} else if algorithm == PSRCconst {
			psrcImpl, iTime, err := initPSRC(wr.Strings, eps)
			if err != nil {
				fmt.Printf("Unable to complete the benchmark: %s\n", err)
				os.Exit(-1)
			}
			impl = psrcImpl
			initTime = iTime
		} else {
			fmt.Errorf(`insert an algorithm between "lprc" and "psrc" \n`)
			os.Exit(1)
		}

		bdSize := impl.GetBitDataSize()
		totalBitSize := totalSize(bdSize)
		fmt.Printf("Initialization time:   %v\n", initTime)
		fmt.Printf("Size of the structure: %d bits\n", totalBitSize)
		fmt.Println()

		finalResults := &Result{
			InitTime:             toMilliseconds(initTime),
			Epsilon:              eps,
			StructureSize:        bdSize,
			UncompressedDataSize: getBitSize(wr.Strings),
		}

		var searchTime time.Time
		var elapsedTime time.Duration
		totalSearchTime := time.Duration(0)
		updateResult := updateResultTemplate(verbose, len(wrp.Strings))
		for _, prefix := range wrp.Strings {
			searchTime = time.Now()
			result, err := impl.FullPrefixSearch(prefix)
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

		allResults = append(allResults, finalResults)
	}
	saveAllToFile(allResults, outputFile)
}

func totalSize(bdSize map[string]uint64) uint64 {
	var totalBitSize uint64
	for _, size := range bdSize {
		totalBitSize += size
	}
	return totalBitSize
}
