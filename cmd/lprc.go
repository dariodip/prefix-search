package cmd

import (
	"fmt"

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
		fmt.Println("lprc called") // TODO
	},
}

func init() {

	rootCmd.AddCommand(lprcCmd)
	lprcCmd.Flags().StringVarP(&inputFile, "INPUT_FILE", "i", "", "Input file containing"+
		" all the word to build up the dictionary")
	lprcCmd.MarkFlagRequired("INPUT_FILE")
	lprcCmd.MarkFlagFilename("INPUT_FILE")

	lprcCmd.Flags().StringVarP(&inputPrefixFile, "INPUT_PREFIX_FILE", "p", "", "Input"+
		" file containing all the prefix to search on the dictionary")
	lprcCmd.MarkFlagRequired("INPUT_PREFIX_FILE")
	lprcCmd.MarkFlagFilename("INPUT_PREFIX_FILE")

}
