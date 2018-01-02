package memmap

import (
	"fmt"
	"reflect"
	"strconv"
)

func (m *mapper) mapString(stringVal reflect.Value, inlineable bool) (nodeID, string) {
	quoted := fmt.Sprintf("\\\"%s\\\"", stringVal.String())
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
