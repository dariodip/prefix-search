package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// psrcCmd represents the psrc command
var psrcCmd = &cobra.Command{
	Use:   "psrc",
	Short: "Run psrc algorithm",
	Long: `psrc (Prefix-Suffix Rear Coding) is an algorithm designed by Mattia Tomeo and Dario Di Pasquale, 
ispired by the paper "Compressed Cache-Oblivious String B-Tree". 

Our implementation takes in input two files: 
	- a file containing all the worlds to add to the dictionary (-i);
	- a file containing all the prefixes to search on the built dictionary (-p).

TODO...
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(inputFile)
		fmt.Println(inputPrefixFile)
		fmt.Println("psrc called") // TODO
	},
}

func init() {
	rootCmd.AddCommand(psrcCmd)

	psrcCmd.Flags().StringVarP(&inputFile, "INPUT_FILE", "i", "", "Input file containing"+
		" all the word to build up the dictionary")
	psrcCmd.MarkFlagRequired("INPUT_FILE")
	psrcCmd.MarkFlagFilename("INPUT_FILE")

	psrcCmd.Flags().StringVarP(&inputPrefixFile, "INPUT_PREFIX_FILE", "p", "", "Input"+
		" file containing all the prefix to search on the dictionary")
	psrcCmd.MarkFlagRequired("INPUT_PREFIX_FILE")
	psrcCmd.MarkFlagFilename("INPUT_PREFIX_FILE")
}
