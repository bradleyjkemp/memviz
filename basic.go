package memviz

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type nodeKey string
type nodeID int

var (
	// for values that aren't addressable keep an incrementing counter instead
	keyCounter int
)

const (
	kTagName         = "memviz"
	nilKey   nodeKey = "nil0"
)

type mapper struct {
	writer        io.Writer
	nodeIDs       map[nodeKey]nodeID
	nodeSummaries map[nodeKey]string
	config        *Config
}

func (m *mapper) mapPtrIface(iVal reflect.Value, inlineable bool, depth uint32) (nodeID, string) {
	pointee := iVal.Elem()
	key := getNodeKey(iVal)

	// inlineable=false so an invalid parentID is fine
	pointeeNode, pointeeSummary := m.mapValue(pointee, 0, false, depth+1)
	summary := escapeString(iVal.Type().String())
	m.nodeSummaries[key] = summary

	if !pointee.IsValid() {
		m.nodeSummaries[key] += "(" + pointeeSummary + ")"
		return pointeeNode, m.nodeSummaries[key]
	}

	if !inlineable {
		id := m.newBasicNode(iVal, summary)
		fmt.Fprintf(m.writer, "  %d:name -> %d:name;\n", id, pointeeNode)
		return id, summary
	}

	return pointeeNode, summary
}

func (m *mapper) mapString(stringVal reflect.Value, inlineable bool) (nodeID, string) {
	// We want the output to look like a Go quoted string literal. The first
	// Quote achieves that. The second is to quote it for graphviz itself.
	quoted := strconv.Quote(strconv.Quote(stringVal.String()))
	// Lastly, quoting adds quotation-marks around the string, but it is
	// inserted into a graphviz string literal, so we have to remove those.
	quoted = quoted[1 : len(quoted)-1]
	if inlineable {
		return 0, quoted
	}
	m.nodeSummaries[getNodeKey(stringVal)] = "string"
	return m.newBasicNode(stringVal, quoted), "string"
}

func (m *mapper) mapBool(stringVal reflect.Value, inlineable bool) (nodeID, string) {
	value := fmt.Sprintf("%t", stringVal.Bool())
	if inlineable {
		return 0, value
	}
	m.nodeSummaries[getNodeKey(stringVal)] = "bool"
	return m.newBasicNode(stringVal, value), "bool"
}

func (m *mapper) mapInt(numVal reflect.Value, inlineable bool) (nodeID, string) {
	printed := strconv.Itoa(int(numVal.Int()))
	if inlineable {
		return 0, printed
	}
	m.nodeSummaries[getNodeKey(numVal)] = "int"
	return m.newBasicNode(numVal, printed), "int"
}

func (m *mapper) mapUint(numVal reflect.Value, inlineable bool) (nodeID, string) {
	printed := strconv.Itoa(int(numVal.Uint()))
	if inlineable {
		return 0, printed
	}
	m.nodeSummaries[getNodeKey(numVal)] = "uint"
	return m.newBasicNode(numVal, printed), "uint"
}

func (m *mapper) getNodeID(iVal reflect.Value) nodeID {
	// have to key on kind and address because a struct and its first element have the same UnsafeAddr()
	key := getNodeKey(iVal)
	var id nodeID
	var ok bool
	if id, ok = m.nodeIDs[key]; !ok {
		id = nodeID(len(m.nodeIDs))
		m.nodeIDs[key] = id
		return id
	}

	return id
}

func (m *mapper) newBasicNode(iVal reflect.Value, text string) nodeID {
	id := m.getNodeID(iVal)
	fmt.Fprintf(m.writer, "  %d [label=\"<name> %s\"];\n", id, text)
	return id
}

func (m *mapper) mapValue(iVal reflect.Value, parentID nodeID, inlineable bool, depth uint32) (nodeID, string) {
	if !iVal.IsValid() {
		// zero value => probably result of nil pointer
		return m.nodeIDs[nilKey], m.nodeSummaries[nilKey]
	}

	key := getNodeKey(iVal)
	if summary, ok := m.nodeSummaries[key]; ok {
		// already seen this address so no need to map again
		return m.nodeIDs[key], summary
	}
	switch iVal.Kind() {
	// Indirections
	case reflect.Ptr, reflect.Interface:
		return m.mapPtrIface(iVal, inlineable, depth)
	// Collections
	case reflect.Struct:
		if m.config.maxDepth > 0 && depth > m.config.maxDepth {
			return -1, "" // max depth reached
		}
		return m.mapStruct(iVal, depth)
	case reflect.Slice, reflect.Array:
		if m.config.maxDepth > 0 && depth > m.config.maxDepth {
			return -1, "" // max depth reached
		}
		return m.mapSlice(iVal, parentID, inlineable, depth)
	case reflect.Map:
		if m.config.maxDepth > 0 && depth > m.config.maxDepth {
			return -1, "" // max depth reached
		}
		return m.mapMap(iVal, parentID, inlineable, depth)

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

func getNodeKey(val reflect.Value) nodeKey {
	if val.CanAddr() {
		return nodeKey(fmt.Sprint(val.Kind()) + fmt.Sprint(val.UnsafeAddr()))
	}

	// reverse order of type and "address" to prevent (incredibly unlikely) collisions
	keyCounter++
	return nodeKey(fmt.Sprint(keyCounter) + fmt.Sprint(val.Kind()))
}
