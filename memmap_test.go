package memmap

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

type basics struct {
	int    int
	string string
	slice  []string
	ptr    *string
}

func TestBasicTypes(t *testing.T) {
	str := "Hello"
	b := &basics{
		1,
		"Hi",
		[]string{"Hello", "World"},
		&str,
	}

	buf := &bytes.Buffer{}
	Map(buf, b)
	fmt.Println(buf.String())
	cupaloy.SnapshotT(t, buf)
}

type tree struct {
	id    int
	left  *tree
	right *tree
}

func TestTree(t *testing.T) {
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

	b := &bytes.Buffer{}
	Map(b, root)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b)
}

func TestSliceTree(t *testing.T) {
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

	slice := []*tree{root, root.left, root.right, leaf}

	b := &bytes.Buffer{}
	Map(b, &slice)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b)
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
	cupaloy.SnapshotT(t, b)
}

type structMap struct {
	id    string
	links map[*structMap]bool
}

func TestMap(t *testing.T) {
	leaf := &structMap{
		"leaf",
		nil,
	}

	leaf2 := &structMap{
		"leaf2",
		nil,
	}

	parent := &structMap{
		"parent",
		map[*structMap]bool{
			leaf:  true,
			leaf2: true,
		},
	}

	leaf.links = map[*structMap]bool{parent: true}
	parent.links[parent] = true

	b := &bytes.Buffer{}
	Map(b, parent)
	fmt.Println(b.String())

	// TODO: enable snapshot once map keys are sorted (and so this has stable output)
	//cupaloy.SnapshotT(t, b)
}

func TestPointerChain(t *testing.T) {
	str := "Hello world"
	str2 := &str
	str3 := &str2
	str4 := &str3

	b := &bytes.Buffer{}
	Map(b, &str4)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b)
}
