package memviz_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/bradleyjkemp/memviz"
)

type basicNumerics struct {
	uint8      uint8
	uint32     uint32
	uint64     uint64
	int8       int8
	int16      int16
	int32      int32
	int64      int64
	float32    float32
	float64    float64
	complex64  complex64
	complex128 complex128
	byte       byte
	rune       rune
	uint       uint
	int        int
	uintptr    uintptr

	Ptruint32     *uint32
	Ptruint64     *uint64
	Ptrint8       *int8
	Ptrint16      *int16
	Ptrint32      *int32
	Ptrint64      *int64
	Ptrfloat32    *float32
	Ptrfloat64    *float64
	Ptrcomplex64  *complex64
	Ptrcomplex128 *complex128
	Ptrbyte       *byte
	Ptrrune       *rune
	Ptruint       *uint
	Ptrint        *int
	Ptruintptr    *uintptr
}

type basics struct {
	numerics *basicNumerics
	string   string
	slice    []string
	ptr      *string
	iface    interface{}
}

func TestBasicTypes(t *testing.T) {
	str := "Hello"
	b := &basics{
		new(basicNumerics),
		"Hi",
		[]string{"Hello", "World"},
		&str,
		"interfaceValue",
	}

	v := reflect.ValueOf(b.numerics).Elem()
	for i := 0; i < v.NumField(); i++ {
		if f := v.Field(i); f.Kind() == reflect.Ptr {
			fv := reflect.New(f.Type().Elem())
			f.Set(fv)
		}
	}

	buf := &bytes.Buffer{}
	memviz.Map(buf, b)
	fmt.Println(buf.String())
	cupaloy.SnapshotT(t, buf.Bytes())
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
	memviz.Map(b, root)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b.Bytes())
}

func TestVariadicArguments(t *testing.T) {
	leaf := &tree{
		0,
		nil,
		nil,
	}
	inner1 := &tree{
		1,
		nil,
		leaf,
	}
	inner2 := &tree{
		2,
		leaf,
		nil,
	}
	root1 := &tree{
		3,
		inner1,
		inner2,
	}
	root2 := &tree{
		4,
		inner2,
		nil,
	}

	b := &bytes.Buffer{}
	memviz.Map(b, root1, root2)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b.Bytes())
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
	memviz.Map(b, &slice)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b.Bytes())
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
	memviz.Map(b, f5)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b.Bytes())
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
	memviz.Map(b, parent)
	fmt.Println(b.String())

	// TODO: enable snapshot assertion once map keys are sorted (and so this has stable output)
	err := cupaloy.Snapshot(b.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func TestPointerChain(t *testing.T) {
	str := "Hello world"
	str2 := &str
	str3 := &str2
	str4 := &str3

	b := &bytes.Buffer{}
	memviz.Map(b, &str4)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b.Bytes())
}

func TestPointerAliasing(t *testing.T) {
	leaf := "leaf"
	parent0 := &leaf
	parent1 := &parent0
	parent2 := &leaf
	root := struct {
		left  **string
		right *string
	}{
		parent1,
		parent2,
	}

	b := &bytes.Buffer{}
	memviz.Map(b, &root)
	fmt.Println(b.String())
	cupaloy.SnapshotT(t, b.Bytes())
}
