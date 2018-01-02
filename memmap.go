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

type nodeKey string
type nodeId int

var nilKey nodeKey = "nil0"

type mapper struct {
	writer        io.Writer
	nodeIds       map[nodeKey]nodeId
	nodeSummaries map[nodeKey]string
}

func (m *mapper) getNodeId(iVal reflect.Value) nodeId {
	// have to key on kind and address because a struct and its first element have the same UnsafeAddr()
	key := getNodeKey(iVal)
	var id nodeId
	var ok bool
	if id, ok = m.nodeIds[key]; !ok {
		id = nodeId(len(m.nodeIds))
		m.nodeIds[key] = id
		return id
	}

	return id
}

func getNodeKey(val reflect.Value) nodeKey {
	return nodeKey(fmt.Sprint(val.Kind(), val.UnsafeAddr()))
}

func (m *mapper) newBasicNode(iVal reflect.Value, text string) nodeId {
	id := m.getNodeId(iVal)
	fmt.Fprintf(m.writer, "  %d [label=\"<name> %s\"];\n", id, text)
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
	(&mapper{w, map[nodeKey]nodeId{nilKey: 0}, map[nodeKey]string{nilKey: "nil"}}).mapValue(iVal, false)
	fmt.Fprintln(w, "}")
}

func (m *mapper) mapValue(iVal reflect.Value, inlineable bool) (nodeId, string) {
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
		childId, summary := m.mapPtr(iVal, inlineable)
		m.nodeSummaries[key] = summary
		return childId, summary

	// Collections
	case reflect.Struct:
		childId, summary := m.mapStruct(iVal)
		m.nodeSummaries[key] = summary
		return childId, summary

	case reflect.Slice:
		fallthrough
	case reflect.Array:
		childId, summary := m.mapSlice(iVal)
		m.nodeSummaries[key] = summary
		return childId, summary

	// Simple types
	case reflect.String:
		quoted := fmt.Sprintf("\\\"%s\\\"", iVal.String())
		if inlineable {
			return 0, quoted
		}
		m.nodeSummaries[key] = "string"
		return m.newBasicNode(iVal, quoted), "string"

	// Numbers
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uint, reflect.Int:
		printed := strconv.Itoa(int(iVal.Int()))
		if inlineable {
			return 0, printed
		}
		m.nodeSummaries[key] = "int"
		return m.newBasicNode(iVal, printed), "int"

	// If we've missed anything then just fmt.Sprint it
	default:
		return m.newBasicNode(iVal, fmt.Sprint(iVal.Interface())), iVal.Kind().String()
	}
}

func (m *mapper) mapPtr(iVal reflect.Value, inlineable bool) (nodeId, string) {
	pointee := iVal.Elem()
	pointeeNode, summary := m.mapValue(pointee, false)

	if !inlineable {
		id := m.newBasicNode(iVal, "*"+summary)
		fmt.Fprintf(m.writer, "  %d:name -> %d:name;\n", id, pointeeNode)
		return id, "*" + summary
	}

	return pointeeNode, "*" + summary
}

func (m *mapper) mapStruct(structVal reflect.Value) (nodeId, string) {
	id := m.getNodeId(structVal)

	var fields []string
	var links []string

	uType := structVal.Type()
	for index := 0; index < uType.NumField(); index++ {
		field := structVal.Field(index)
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		fieldId, summary := m.mapValue(field, true)

		// if field was inlined (id == 0) then print summary, else just the name and a link to the actual
		if fieldId == 0 {
			fields = append(fields, fmt.Sprintf("%s: %s", uType.Field(index).Name, summary))
		} else {
			fields = append(fields, uType.Field(index).Name)
			links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", id, index, fieldId))
		}
	}

	node := fmt.Sprintf("  %d [label=\"<name> %s", id, structVal.Type().Name())
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

func (m *mapper) mapSlice(sliceVal reflect.Value) (nodeId, string) {
	id := m.getNodeId(sliceVal)

	length := sliceVal.Len()
	node := fmt.Sprintf("  %d [label=\"<name> %s", id, sliceVal.Type().String())
	var links []string
	for index := 0; index < length; index++ {
		indexId, summary := m.mapValue(sliceVal.Index(index), true)
		if indexId == 0 {
			// field was inlined
			node += fmt.Sprintf("|<f%d> %d: %s", index, index, summary)
		} else {
			// need a link to the new node and don't care about the summary
			node += fmt.Sprintf("|<f%d> %d", index, index)
			links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", id, index, indexId))
		}
	}
	node += "\"];\n"

	fmt.Fprint(m.writer, node)
	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	return id, sliceVal.Type().String()
}
