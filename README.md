# memmap [![Build Status](https://travis-ci.org/bradleyjkemp/memmap.svg?branch=master)](https://travis-ci.org/bradleyjkemp/memmap) [![Coverage Status](https://coveralls.io/repos/github/bradleyjkemp/memmap/badge.svg)](https://coveralls.io/github/bradleyjkemp/memmap?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/bradleyjkemp/memmap)](https://goreportcard.com/report/github.com/bradleyjkemp/memmap) [![GoDoc](https://godoc.org/github.com/bradleyjkemp/memmap?status.svg)](https://godoc.org/github.com/bradleyjkemp/memmap)

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
    <td width="60%"><image src="images/fib.svg"></td>
  </tr>
</table>

`memmap` takes a pointer to an arbitrary data structure and generates an easy to understand graph.

Simply pass in your data structure like so: ```memmap.Map(out, &data)``` and then pipe the output into graphviz.
