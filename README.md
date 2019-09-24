# Prefix Search

[![Build Status](https://travis-ci.com/dariodip/prefix-search.svg?token=NZ9VK4sB4UsVShV1p8wD&branch=master)](https://travis-ci.com/dariodip/prefix-search)
[![GoDoc](https://godoc.org/github.com/spf13/cobra?status.svg)](https://godoc.org/github.com/dariodip/prefix-search)
[![Go Report Card](https://goreportcard.com/badge/github.com/dariodip/prefix-search)](https://goreportcard.com/report/github.com/dariodip/prefix-search)


prefix-search is an implementation of the [paper](https://link.springer.com/chapter/10.1007/978-3-642-40450-4_40) *Compressed Cache-Oblivious String B-tree* of Paolo Ferragina and Rossano Venturini. We developed the proposed algorithm (LPRC) and a new one (PSRC), giving you the ability to deal with online dictionaries of strings in an unspecified order.

## Getting Started

With these instructions, you will get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites
You need [make](https://www.gnu.org/software/make/) and [Golang](https://golang.org/) installed in order to run all the command listed next. 

You also need [Python 3](https://www.python.org/) to run our scripts.

### Installing

In order to install all the required dependencies, in the project root directory run the following:
```
make install
```

In order to create an executable, in the project root directory run the followings:
```
make build
```
this will install all the dependencies, run all the tests and create an executable inside your `$GOPATH/bin` directory.

## Release
If you want to create an executable for a specific platform, run the following:
```
make [VERSION=<version>] <platforms>
```
where the platforms are `windows`, `linux` and `darwin`. If not specified, the value of the parameter `VERSION` will be
`vlatest`.

If you want to create an executable for each platform, then simply run:
```
make [VERSION=<version>] release
```
## Usage
Prefix search can be used with 4 commands:
* **console**: 
```
prefix-search console --help                                                                12:35   08.06.18 
Using "console" you can start an interactive console that gives you the opportunity
to, given a preloaded dataset, to find prefixes interactively.

Usage:
  prefix-search console [flags]

Flags:
  -a, --algorithm string    Algorithmto use (default "lprc")
  -e, --epsilon float       Epsilon is the parametergiven to the algorithm in order to decide how many bits compress in the trie.
  -h, --help                help for console
  -i, --input_file string   Input file containing all the word to build up the dictionary.
```
* **lprc**:
```
prefix-search lprc --help
lprc (Locality Preserving Rear Coding) is an algorithm designed by Paolo Ferragina and Rossano Venturini
in their paper "Compressed Cache-Oblivious String B-Tree". 

Our implementation takes in input: 
	- a file containing all the worlds to add to the dictionary (-i);
	- a file containing all the prefixes to search on the built dictionary (-p).
	- the epsilon to use in order to build our structure

All the results will be saved into a json file.

Usage:
  prefix-search lprc [flags]

Flags:
  -e, --epsilon float         Epsilon is the parametergiven to the algorithm in order to decide how many bits compress in the trie.
  -h, --help                  help for lprc
  -i, --input_file string     Input file containing all the word to build up the dictionary.
  -p, --input_p_file string   Input file containing all the prefix to search on the dictionary.
  -o, --output_file string    Output file containing the final output of lprc, with information about the memory usage and the time elapsed.
                              Default <word filename>-<prefix file name>-<epsilon>.json
  -v, --verbose               Detailed Output
```
* **psrc**:
```
prefix-search psrc --help 
psrc (Prefix-Suffix Rear Coding) is an algorithm designed by Mattia Tomeo and Dario Di Pasquale, 
inspired by the paper "Compressed Cache-Oblivious String B-Tree". 

Our implementation takes in input two files: 
	- a file containing all the worlds to add to the dictionary (-i);
	- a file containing all the prefixes to search on the built dictionary (-p).

All the results will be saved into a json file.

Usage:
  prefix-search psrc [flags]

Flags:
  -e, --epsilon float         Epsilon is the parametergiven to the algorithm in order to decide how many bits compress in the trie.
  -h, --help                  help for psrc
  -i, --input_file string     Input file containing all the word to build up the dictionary
  -p, --input_p_file string   Input file containing all the prefix to search on the dictionary
  -o, --output_file string    Output file containing the final output of lprc, with information about the memory usage and the time elapsed.
                              Default <word filename>-<prefix file name>-<epsilon>.json
  -v, --verbose               Detailed Output
```
* **fullbenchmark**:
```
prefix-search fullbenchmark --help
This command gives you the ability to run a sophisticated benchmarking test.
You can select the file to open as dataset, the file to open as prefix, the lower value
of epsilon, the higher value of epsilon and the step with which increase the value of it.

The test gives you only a JSON file containing all the results of the test.

Usage:
  prefix-search fullbenchmark [flags]

Flags:
  -a, --algorithm string      Algorithmto use (default "lprc")
  -h, --help                  help for fullbenchmark
  -i, --input_file string     Input file containing all the word to build up the dictionary.
  -p, --input_p_file string   Input file containing all the prefix to search on the dictionary.
  -x, --max_epsilon float     Maximum value of Epsilon: the parameter given to the algorithm in order to decide how many bits compress in the trie.
  -n, --min_epsilon float     Minimum value of Epsilon: the parameter given to the algorithm in order to decide how many bits compress in the trie.
  -o, --output_file string    Output file containing the final output of lprc, with information about the memory usage and the time elapsed.
                              Default <algorithm>-<word filename>-<prefix file name>-<min_epsilon>-<max_epsilon>.json
  -s, --step float            Step value with which increment the value of epsilon
  -v, --verbose               Detailed Output
```
## Running the tests

All the test are built using the package [testing](https://golang.org/pkg/testing/).

To execute all our tests, run the following
```
make test
```

## Scripts

In our project we have implemented some utility scripts:
* **resources/generate_dataset.py** - create word datasets of different size from the single dataset **words.txt**
## Built With

* [Golang](https://golang.org/)

## Authors

* [**Dario Di Pasquale**](https://github.com/dariodip)
* [**Mattia Tomeo**](https://github.com/mattiatomeo)
