package cmd

import (
	"fmt"

	"github.com/dariodip/prefix-search/prefix-search/stringcoding"
	"github.com/dariodip/prefix-search/word-reader"
	"github.com/spf13/cobra"
)

// lprcCmd represents the lprc command
var lprcCmd = &cobra.Command{
	Use:   "lprc",
	Short: "Run lprc algorithm",
	Long: `lprc (Locality Preserving Rear Coding) is an algorithm designed by Paolo Ferragina and Rossano Venturini
in their paper "Compressed Cache-Oblivious String B-Tree". 

Our implementation takes in input two files: 
	- a file containing all the worlds to add to the dictionary (-i);
	- a file containing all the prefixes to search on the built dictionary (-p).

TODO...
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(inputFile)
		fmt.Println(inputPrefixFile)
		runLPRC()
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

}

func runLPRC() {

	// load words
	wr := word_reader.New(inputFile)
	wr.ReadLines()

	// load prefix
	wrp := word_reader.New(inputPrefixFile)
	wrp.ReadLines()

	lprcImpl := stringcoding.NewLPRC(wr.Strings, epsilon)
	if err := lprcImpl.PopulateLPRC(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(lprcImpl.String()) // TODO
}
