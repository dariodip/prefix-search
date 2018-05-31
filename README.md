# Prefix Search

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/prefix-search/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link)
[![Build Status](https://travis-ci.com/dariodip/prefix-search.svg?token=NZ9VK4sB4UsVShV1p8wD&branch=master)](https://travis-ci.com/dariodip/prefix-search)

prefix-search is an implementation of the [paper](https://link.springer.com/chapter/10.1007/978-3-642-40450-4_40) *Compressed Cache-Oblivious String B-tree* of Paolo Ferragina and Rossano Venturini. We developed the proposed algorithm (LPRC) and a new one (PSRC), giving you the ability to deal with online dictionaries of strings in an unspecified order.

## Getting Started

With these instructions you will get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites
You need make and [Golang](https://golang.org/) installed in order to run all the command listed next. 

You need [Python 3](https://www.python.org/) if you want to run our scripts.
### Installing

In order to install all the required dependencies, in the project root directory run the following:
```
make install
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
