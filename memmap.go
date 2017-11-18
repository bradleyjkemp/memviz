package memmap

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"io"

	"github.com/davecgh/go-spew/spew"
)

var spewer = &spew.ConfigState{
	Indent:                  "  ",
	SortKeys:                true, // maps should be spewed in a deterministic order
	DisablePointerAddresses: true, // don't spew the addresses of pointers
	DisableCapacities:       true, // don't spew capacities of collections
	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
	MaxDepth:                1,
}

func Map(w io.Writer, i interface{}) {
	iVal := reflect.ValueOf(i)
	fmt.Fprintln(w, "digraph {")
	imap(w, iVal)
	fmt.Fprintln(w, "}")
}

func imap(w io.Writer, iVal reflect.Value) uintptr {
	if iVal.Kind() == reflect.Ptr || iVal.Kind() == reflect.Interface {
		if !iVal.IsNil() {
			return imap(w, iVal.Elem())
		}
		return 0
	}

	if iVal.Kind() != reflect.Struct {
		return 0
	}

	iValAddr := iVal.UnsafeAddr()
	nodeLabel := strings.Replace(spewer.Sdump(iVal.Interface()), "\n", "\\n", -1)
	fmt.Fprintf(w, "  %d [label=\"%s\"];\n", iValAddr, nodeLabel)

	uType := iVal.Type()
	for i := 0; i < uType.NumField(); i++ {
		field := iVal.Field(i)
		if field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
			field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
			ptr := imap(w, field)
			if ptr != 0 {
				fmt.Fprintf(w, "  %d -> %d;\n", iValAddr, ptr)
			}
		}
	}

	return iValAddr
}
