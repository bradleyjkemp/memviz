# memmap [![Build Status](https://travis-ci.org/bradleyjkemp/memmap.svg?branch=master)](https://travis-ci.org/bradleyjkemp/memmap) [![Coverage Status](https://coveralls.io/repos/github/bradleyjkemp/memmap/badge.svg)](https://coveralls.io/github/bradleyjkemp/memmap?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/bradleyjkemp/memmap)](https://goreportcard.com/report/github.com/bradleyjkemp/memmap) [![GoDoc](https://godoc.org/github.com/bradleyjkemp/memmap?status.svg)](https://godoc.org/github.com/bradleyjkemp/memmap) 

Take arbitrary data structures and turn them into a easy to understand graph:

![fibonacci](images/fib.svg)

Just pass a pointer to your data structure like so: ```memmap.Map(out, &data)``` and then pipe the output into graphviz.
