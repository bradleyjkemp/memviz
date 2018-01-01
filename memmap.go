package memmap

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"unsafe"
)

//var spewer = &spew.ConfigState{
//	Indent:                  "  ",
//	SortKeys:                true, // maps should be spewed in a deterministic order
//	DisablePointerAddresses: true, // don't spew the addresses of pointers
//	DisableCapacities:       true, // don't spew capacities of collections
//	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
//	MaxDepth:                1,
//}

type mapper struct {
	writer        io.Writer
	nodeIds       map[uintptr]nodeId
	nodeSummaries map[uintptr]string
}

func (m *mapper) getNodeId(addr uintptr) nodeId {
	var id nodeId
	var ok bool
	if id, ok = m.nodeIds[addr]; !ok {
		id = nodeId(len(m.nodeIds))
		m.nodeIds[addr] = id
		return id
	}

	return id
}

func (m *mapper) newBasicNode(addr uintptr, text string) nodeId {
	id := m.getNodeId(addr)
	m.nodeIds[addr] = id
	fmt.Fprintf(m.writer, "  %d [label=\"<name> %s\"];\n", m.nodeIds[addr], text)
	return id
}

// Map prints out a Graphviz digraph of the given datastructure to the given io.Writer
func Map(w io.Writer, i interface{}) {
	iVal := reflect.ValueOf(i)
	if !iVal.CanAddr() {
		if iVal.Kind() != reflect.Ptr && iVal.Kind() != reflect.Interface {
			w.Write([]byte("error: cannot map unaddressable value"))
			return
		}

		iVal = iVal.Elem()
	}

	fmt.Fprintln(w, "digraph structs {")
	fmt.Fprintln(w, "  node [shape=Mrecord];")
	(&mapper{w, map[uintptr]nodeId{0: 0}, map[uintptr]string{0: "nil"}}).mapValue(iVal, false)
	fmt.Fprintln(w, "}")
}

type nodeId int

func (m *mapper) mapValue(iVal reflect.Value, inlineable bool) (nodeId, string) {
	if !iVal.IsValid() {
		// zero value => probably result of nil pointer
		return m.nodeIds[0], m.nodeSummaries[0]
	}

	iValAddr := iVal.UnsafeAddr()
	if _, ok := m.nodeIds[iValAddr]; ok {
		// already seen this address so no need to map again
		return m.nodeIds[iValAddr], m.nodeSummaries[iValAddr]
	}

	switch iVal.Kind() {
	case reflect.Ptr:
		childId, summary := m.mapPtr(iVal, inlineable)
		m.nodeSummaries[iValAddr] = summary
		return childId, summary

	// Collections
	case reflect.Struct:
		childId, summary := m.mapStruct(iVal)
		m.nodeSummaries[iValAddr] = summary
		return childId, summary

	case reflect.Slice:
		fallthrough
	case reflect.Array:
		childId, summary := m.mapSlice(iVal)
		m.nodeSummaries[iValAddr] = summary
		return childId, summary

	// Simple types
	case reflect.String:
		m.nodeSummaries[iValAddr] = "string"
		return m.newBasicNode(iValAddr, "\\\""+iVal.String()+"\\\""), "string"
	case reflect.Int:
		m.nodeSummaries[iValAddr] = "int"
		return m.newBasicNode(iValAddr, strconv.Itoa(int(iVal.Int()))), "string"
	default:
		fmt.Println(iVal.Kind())
		return m.newBasicNode(iValAddr, fmt.Sprint(iVal.Interface())), iVal.Kind().String()
	}
}

func (m *mapper) mapPtr(iVal reflect.Value, inlineable bool) (nodeId, string) {
	iValAddr := iVal.UnsafeAddr()
	pointee := iVal.Elem()
	pointeeNode, summary := m.mapValue(pointee, false)

	if !inlineable {
		id := m.newBasicNode(iValAddr, "*"+summary)
		fmt.Fprintf(m.writer, "  %d:name -> %d:name;\n", m.nodeIds[iValAddr], pointeeNode)
		return id, "*" + summary
	}

	return pointeeNode, "*" + summary
}

func (m *mapper) mapStruct(iVal reflect.Value) (nodeId, string) {
	iValAddr := iVal.UnsafeAddr()

	id := m.getNodeId(iValAddr)

	var fields []string
	var links []string

	uType := iVal.Type()
	for index := 0; index < uType.NumField(); index++ {
		field := iVal.Field(index)
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		if field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
			fieldAddr, summary := m.mapValue(field, true)
			fields = append(fields, fmt.Sprintf("%s: %s", uType.Field(index).Name, summary))
			if fieldAddr != 0 {
				links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", id, index, fieldAddr))
			}
		} else {
			fields = append(fields, fmt.Sprintf("%s: %s", uType.Field(index).Name, fmt.Sprint(field.Interface())))
		}
	}

	node := fmt.Sprintf("  %d [label=\"<name> %s", id, iVal.Type().Name())
	for index, name := range fields {
		node += fmt.Sprintf("|<f%d> %s", index, name)
	}
	node += "\"];\n"

	fmt.Fprint(m.writer, node)
	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	return id, uType.String()
}

func (m *mapper) mapSlice(slice reflect.Value) (nodeId, string) {
	addr := slice.UnsafeAddr()
	id := m.getNodeId(addr)

	length := slice.Len()
	node := fmt.Sprintf("  %d [label=\"<name> %s", id, slice.Type().String())
	var links []string
	for index := 0; index < length; index++ {
		indexId, summary := m.mapValue(slice.Index(index), true)
		node += fmt.Sprintf("|<f%d> %d: %s", index, index, summary)
		links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", id, index, indexId))
	}
	node += "\"];\n"

	fmt.Fprint(m.writer, node)
	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	return id, slice.Type().String()
}
