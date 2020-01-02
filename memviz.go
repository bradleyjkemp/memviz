package memviz // import "github.com/bradleyjkemp/memviz"

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

//var spewer = &spew.ConfigState{
//	Indent:                  "  ",
//	SortKeys:                true, // maps should be spewed in a deterministic order
//	DisablePointerAddresses: true, // don't spew the addresses of pointers
//	DisableCapacities:       true, // don't spew capacities of collections
//	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
//	MaxDepth:                1,
//}

// Map prints the given datastructure using the default config
func Map(w io.Writer, is ...interface{}) {
	defaultConfig().Map(w, is...)
}

// Map prints out a Graphviz digraph of the given datastructure to the given io.Writer
func (c *Config) Map(w io.Writer, is ...interface{}) {
	var iVals []reflect.Value
	for _, i := range is {
		iVal := reflect.ValueOf(i)
		if !iVal.CanAddr() {
			if iVal.Kind() != reflect.Ptr && iVal.Kind() != reflect.Interface {
				fmt.Fprint(w, "error: cannot map unaddressable value")
				return
			}

			iVal = iVal.Elem()

		}
		iVals = append(iVals, iVal)
	}

	m := &mapper{
		w,
		map[nodeKey]nodeID{nilKey: 0},
		map[nodeKey]string{nilKey: "nil"},
		c,
	}

	fmt.Fprintln(w, "digraph structs {")
	fmt.Fprintln(w, "  node [shape=Mrecord];")
	for _, iVal := range iVals {
		m.mapValue(iVal, 0, false, 0) // start at zero depth
	}
	fmt.Fprintln(w, "}")
}

var escaper = strings.NewReplacer(
	"{", "\\{",
	"}", "\\}",
	"\"", "\\\"",
)

func escapeString(s string) string {
	return escaper.Replace(s)
}
