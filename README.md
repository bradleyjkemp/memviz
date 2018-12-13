# memviz [![Build Status](https://travis-ci.org/bradleyjkemp/memviz.svg?branch=master)](https://travis-ci.org/bradleyjkemp/memviz) [![Coverage Status](https://coveralls.io/repos/github/bradleyjkemp/memviz/badge.svg)](https://coveralls.io/github/bradleyjkemp/memviz?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/bradleyjkemp/memviz)](https://goreportcard.com/report/github.com/bradleyjkemp/memviz) [![GoDoc](https://godoc.org/github.com/bradleyjkemp/memviz?status.svg)](https://godoc.org/github.com/bradleyjkemp/memviz)

How would you rather debug a data structure?
<table>
  <tr>
    <td>"Pretty" printed</td>
    <td>Visual graph</td>
  </tr>
  <tr>
    <td>
        <pre>
(*test.fib)(0xc04204a5a0)({
 index: (int) 5,
 prev: (*test.fib)(0xc04204a580)({
  index: (int) 4,
  prev: (*test.fib)(0xc04204a560)({
   index: (int) 3,
   prev: (*test.fib)(0xc04204a540)({
    index: (int) 2,
    prev: (*test.fib)(0xc04204a520)({
     index: (int) 1,
     prev: (*test.fib)(0xc04204a500)({
      index: (int) 0,
      prev: (*test.fib)(<nil>),
      prevprev: (*test.fib)(<nil>)
     }),
     prevprev: (*test.fib)(<nil>)
    }),
    prevprev: (*test.fib)(0xc04204a500)({
     index: (int) 0,
     prev: (*test.fib)(<nil>),
     prevprev: (*test.fib)(<nil>)
    })
   }),
   .
   .
   .</pre>
    </td>
    <td width="60%"><image src=".github/fib.svg"></td>
  </tr>
</table>

## Usage
`memviz` takes a pointer to an arbitrary data structure and generates an easy to understand graph.

Simply pass in your data structure like so: ```memviz.Map(out, &data)``` and then pipe the output into graphviz.

For more complete examples see the tests in [memviz_test.go](https://github.com/bradleyjkemp/memviz/blob/master/memviz_test.go).
