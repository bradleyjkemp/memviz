package memmap

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

type btree struct {
	id    int
	left  *btree
	right *btree
}

func TestTree(t *testing.T) {
	root := &btree{
		id: 0,
		left: &btree{
			id: 1,
		},
		right: &btree{
			id: 2,
		},
	}
	leaf := &btree{
		id: 3,
	}

	root.left.right = leaf
	root.right.left = leaf

	b := &bytes.Buffer{}
	Map(b, root)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, strings.Split(b.String(), "\n"))
}

type fib struct {
	index    int
	prev     *fib
	prevprev *fib
}

func TestFib(t *testing.T) {
	f0 := &fib{
		0,
		nil,
		nil,
	}
	f1 := &fib{
		1,
		f0,
		nil,
	}
	f2 := &fib{
		2,
		f1,
		f0,
	}
	f3 := &fib{
		3,
		f2,
		f1,
	}
	f4 := &fib{
		4,
		f3,
		f2,
	}
	f5 := &fib{
		5,
		f4,
		f3,
	}

	b := &bytes.Buffer{}
	Map(b, f5)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, strings.Split(b.String(), "\n"))
}
