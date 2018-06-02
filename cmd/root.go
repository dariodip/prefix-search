package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	VERSION         string
	inputFile       string
	inputPrefixFile string
	algorithm       string
	epsilon         float64
	LPRCconst       = "lprc"
	PSRCconst       = "psrc"
)

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
