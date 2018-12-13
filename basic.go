package memviz

import (
	"fmt"
	"reflect"
	"strconv"
)

func (m *mapper) mapPtrIface(iVal reflect.Value, inlineable bool) (nodeID, string) {
	pointee := iVal.Elem()
	key := getNodeKey(iVal)

	// inlineable=false so an invalid parentID is fine
	pointeeNode, pointeeSummary := m.mapValue(pointee, 0, false)
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
