package main

import (
	"bytes"
	"github.com/bradleyjkemp/memviz"
	"io/ioutil"
)

type tree struct {
	id    int
	left  *tree
	right *tree
}

func main() {
	root := &tree{
		id: 0,
		left: &tree{
			id: 1,
		},
		right: &tree{
			id: 2,
		},
	}
	leaf := &tree{
		id: 3,
	}

	root.left.right = leaf
	root.right.left = leaf

	buf := &bytes.Buffer{}
	memviz.Map(buf, &root)
	err := ioutil.WriteFile("example-tree-data", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
