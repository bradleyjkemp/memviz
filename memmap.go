package memmap

import (
	"fmt"
	"io"
	"reflect"
)

//var spewer = &spew.ConfigState{
//	Indent:                  "  ",
//	SortKeys:                true, // maps should be spewed in a deterministic order
//	DisablePointerAddresses: true, // don't spew the addresses of pointers
//	DisableCapacities:       true, // don't spew capacities of collections
//	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
//	MaxDepth:                1,
//}

type nodeKey string
type nodeID int

var nilKey nodeKey = "nil0"

type mapper struct {
	writer        io.Writer
	nodeIds       map[nodeKey]nodeID
	nodeSummaries map[nodeKey]string
}

// Map prints out a Graphviz digraph of the given datastructure to the given io.Writer
func Map(w io.Writer, i interface{}) {
	iVal := reflect.ValueOf(i)
	if !iVal.CanAddr() {
		if iVal.Kind() != reflect.Ptr && iVal.Kind() != reflect.Interface {
			fmt.Fprint(w, "error: cannot map unaddressable value")
			return
		}

		iVal = iVal.Elem()
	}

	fmt.Fprintln(w, "digraph structs {")
	fmt.Fprintln(w, "  node [shape=Mrecord];")
	(&mapper{w, map[nodeKey]nodeID{nilKey: 0}, map[nodeKey]string{nilKey: "nil"}}).mapValue(iVal, false)
	fmt.Fprintln(w, "}")
}

// for values that aren't addressable keep an incrementing counter instead
var keyCounter int

func getNodeKey(val reflect.Value) nodeKey {
	if val.CanAddr() {
		return nodeKey(fmt.Sprint(val.Kind(), val.UnsafeAddr()))
	}

	// reverse order of type and "address" to prevent collisions
	return nodeKey(fmt.Sprint(keyCounter, val.Kind()))
}

func (m *mapper) getNodeID(iVal reflect.Value) nodeID {
	// have to key on kind and address because a struct and its first element have the same UnsafeAddr()
	key := getNodeKey(iVal)
	var id nodeID
	var ok bool
	if id, ok = m.nodeIds[key]; !ok {
		id = nodeID(len(m.nodeIds))
		m.nodeIds[key] = id
		return id
	}

	return id
}

func (m *mapper) newBasicNode(iVal reflect.Value, text string) nodeID {
	id := m.getNodeID(iVal)
	fmt.Fprintf(m.writer, "  %d [label=\"<name> %s\"];\n", id, text)
	return id
}

func (m *mapper) mapValue(iVal reflect.Value, inlineable bool) (nodeID, string) {
	if !iVal.IsValid() {
		// zero value => probably result of nil pointer
		return m.nodeIds[nilKey], m.nodeSummaries[nilKey]
	}

	key := getNodeKey(iVal)
	if id, ok := m.nodeIds[key]; ok {
		// already seen this address so no need to map again
		return id, m.nodeSummaries[key]
	}

	switch iVal.Kind() {
	case reflect.Ptr:
		childID, summary := m.mapPtr(iVal, inlineable)
		m.nodeSummaries[key] = summary
		return childID, summary

	// Collections
	case reflect.Struct:
		return m.mapStruct(iVal)
	case reflect.Slice, reflect.Array:
		return m.mapSlice(iVal)
	case reflect.Map:
		return m.mapMap(iVal)

	// Simple types
	case reflect.Bool:
		return m.mapBool(iVal, inlineable)
	case reflect.String:
		return m.mapString(iVal, inlineable)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return m.mapInt(iVal, inlineable)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return m.mapUint(iVal, inlineable)

	// If we've missed anything then just fmt.Sprint it
	default:
		return m.newBasicNode(iVal, fmt.Sprint(iVal.Interface())), iVal.Kind().String()
	}
}

func (m *mapper) mapPtr(iVal reflect.Value, inlineable bool) (nodeID, string) {
	pointee := iVal.Elem()
	pointeeNode, summary := m.mapValue(pointee, false)

	if !inlineable {
		id := m.newBasicNode(iVal, "*"+summary)
		fmt.Fprintf(m.writer, "  %d:name -> %d:name;\n", id, pointeeNode)
		return id, "*" + summary
	}

	return pointeeNode, "*" + summary
}
