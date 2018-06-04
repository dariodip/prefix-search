package cmd

import (
	"fmt"

	"bufio"
	"github.com/dariodip/prefix-search/prefix-search/stringcoding"
	"github.com/dariodip/prefix-search/word-reader"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"time"
)

var consoleMarker = "> "

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Start interactive console",
	Long: `Using "console" you can start interactive console that gives you the opportunity
to, given a preloaded dataset, to find more prefixes interactively.`,
	Run: runConsole,
}

func init() {
	rootCmd.AddCommand(consoleCmd)

	consoleCmd.Flags().StringVarP(&inputFile, "input_file", "i", "", "Input file containing"+
		" all the word to build up the dictionary.")
	consoleCmd.MarkFlagRequired("input_file")
	consoleCmd.MarkFlagFilename("input_file")

	consoleCmd.Flags().StringVarP(&algorithm, "algorithm", "a", "lprc", "Algorithm"+
		"to use")
	consoleCmd.MarkFlagFilename("algorithm")

	consoleCmd.Flags().Float64VarP(&epsilon, "epsilon", "e", 0, "Epsilon is the parameter"+
		"given to the algorithm in order to decide how many bits compress in the trie.")
	consoleCmd.MarkFlagRequired("epsilon")

}

func runConsole(cmd *cobra.Command, args []string) {

	fmt.Println("Welcome in prefix-search interactive console.")
	fmt.Println("Enter a prefix to search: ")

	var (
		wr   = word_reader.New(inputFile)
		impl stringcoding.PrefixSearch
	)

	startTime := time.Now()
	lines, err := wr.ReadLines() // read all lines from the file
	if err != nil {
		fmt.Errorf("error in load lines from file: %s \n", err)
	}

	if algorithm == LPRCconst {
		lprcImpl := stringcoding.NewLPRC(wr.Strings, epsilon)
		impl = &lprcImpl
	} else if algorithm == PSRCconst {
		// TODO
		fmt.Errorf("algorithm not yet implemented \n")
		os.Exit(1)
	} else {
		fmt.Errorf(`insert an algorithm between "lprc" and "psrc" \n`)
		os.Exit(1)
	}

	if err := impl.Populate(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Loaded %d words in %v \n", lines, time.Since(startTime))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go handleInterrupt(c)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(consoleMarker)
	for scanner.Scan() {
		prefix := scanner.Text()
		if prefix == "" {
			fmt.Print(consoleMarker)
			continue
		}
		fmt.Println("Searching for strings starting with ", prefix)

		startTime = time.Now()
		strings, err := impl.FullPrefixSearch(prefix)
		if err != nil {
			fmt.Errorf("error: %s", err)
		}
		fmt.Printf("Found %d strings in %v \n", len(strings), time.Since(startTime))
		if len(strings) == 0 {
			fmt.Println("No string found")
		}
		for i, s := range strings {
			endPrint := false
			if i%10 == 0 && i > 0 {
				usage := "[...] hit Enter to continue or q + Enter to end the visualization"
				fmt.Println(usage)
				for {
					scanner.Scan()
					text := scanner.Text()
					if len(text) > 1 || (len(text) == 1 && text != "q") {
						fmt.Println(usage)
					} else if text == "q" {  // The user want to end the visualization
						endPrint = true
						break
					} else {  // The user clicked only enter
						break
					}
				}
			}
			if endPrint {
				break
			}
			fmt.Printf("%d) %s \n", i+1, s)
		}
		fmt.Print(consoleMarker)
	}

}

func handleInterrupt(c chan os.Signal) {
	<-c
	fmt.Printf("Received interrupt signal \n")
	fmt.Println("Bye")
	os.Exit(0)
}
