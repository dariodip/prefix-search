package main

import "github.com/dariodip/prefix-search/cmd"

var (
	// VERSION specifies the version of the build
	VERSION = "0.0.1"
)

func main() {

	cmd.Execute(VERSION)
}
