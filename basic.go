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

func (m *mapper) mapNumber(numVal reflect.Value, inlineable bool) (nodeID, string) {
	printed := strconv.Itoa(int(numVal.Int()))
	if inlineable {
		return 0, printed
	}
	m.nodeSummaries[getNodeKey(numVal)] = "int"
	return m.newBasicNode(numVal, printed), "int"
}
