package memmap

import (
	"bytes"
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

	cupaloy.Snapshot(strings.Split(b.String(), "\n"))
}
