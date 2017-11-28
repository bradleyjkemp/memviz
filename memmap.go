package memmap

import (
	"fmt"
	"reflect"
	"unsafe"

	"io"
)

//var spewer = &spew.ConfigState{
//	Indent:                  "  ",
//	SortKeys:                true, // maps should be spewed in a deterministic order
//	DisablePointerAddresses: true, // don't spew the addresses of pointers
//	DisableCapacities:       true, // don't spew capacities of collections
//	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
//	MaxDepth:                1,
//}

func Map(w io.Writer, i interface{}) {
	iVal := reflect.ValueOf(i)
	fmt.Fprintln(w, "digraph structs {")
	fmt.Fprintln(w, "  node [shape=Mrecord];")
	imap(w, iVal, map[uintptr]int{})
	fmt.Fprintln(w, "}")
}

func imap(w io.Writer, iVal reflect.Value, nodeIds map[uintptr]int) (uintptr, map[uintptr]int) {
	if iVal.Kind() == reflect.Ptr || iVal.Kind() == reflect.Interface {
		if !iVal.IsNil() {
			return imap(w, iVal.Elem(), nodeIds)
		}
		return 0, nodeIds
	}

	if iVal.Kind() != reflect.Struct {
		return 0, nodeIds
	}

	iValAddr := iVal.UnsafeAddr()
	if _, ok := nodeIds[iValAddr]; ok {
		return iValAddr, nodeIds
	}

	nodeIds[iValAddr] = len(nodeIds)

	var fields []string
	var links []string

	uType := iVal.Type()
	for index := 0; index < uType.NumField(); index++ {
		field := iVal.Field(index)
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		if field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
			fields = append(fields, uType.Field(index).Name)
			ptr, rNodeIds := imap(w, field, nodeIds)
			if ptr != 0 {
				nodeIds = rNodeIds
				links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", nodeIds[iValAddr], index, nodeIds[ptr]))
			}
		} else {
			fields = append(fields, fmt.Sprintf("%s: %s", uType.Field(index).Name, field.Interface()))
		}
	}

	node := fmt.Sprintf("  %d [label=\"<name> %s", nodeIds[iValAddr], iVal.Type().Name())
	for index, name := range fields {
		node += fmt.Sprintf("|<f%d> %s", index, name)
	}
	node += "\"];\n"

	fmt.Fprint(w, node)
	for _, link := range links {
		fmt.Fprint(w, link)
	}

	return iValAddr, nodeIds
}
