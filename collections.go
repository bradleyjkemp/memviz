package memmap

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func (m *mapper) mapStruct(structVal reflect.Value) (nodeID, string) {
	uType := structVal.Type()
	id := m.getNodeID(structVal)
	key := getNodeKey(structVal)
	m.nodeSummaries[key] = uType.String()

	var fields string
	var links []string
	for index := 0; index < uType.NumField(); index++ {
		field := structVal.Field(index)
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		fieldID, summary := m.mapValue(field, id, true)

		// if field was inlined (id == 0) then print summary, else just the name and a link to the actual
		if fieldID == 0 {
			fields += fmt.Sprintf("|{<f%d> %s | %s} ", index, uType.Field(index).Name, summary)
		} else {
			fields += fmt.Sprintf("|<f%d> %s", index, uType.Field(index).Name)
			links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", id, index, fieldID))
		}
	}

	node := fmt.Sprintf("  %d [label=\"<name> %s %s \"];\n", id, structVal.Type().Name(), fields)

	fmt.Fprint(m.writer, node)
	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	return id, m.nodeSummaries[key]
}

func (m *mapper) mapSlice(sliceVal reflect.Value, parentID nodeID, inlineable bool) (nodeID, string) {
	id := m.getNodeID(sliceVal)
	key := getNodeKey(sliceVal)
	m.nodeSummaries[key] = sliceVal.Type().String()

	length := sliceVal.Len()
	node := fmt.Sprintf("  %d [label=\"<name> %s", id, sliceVal.Type().String())
	var links []string
	for index := 0; index < length; index++ {
		indexID, summary := m.mapValue(sliceVal.Index(index), id, true)
		node += fmt.Sprintf("|{<index%d> %d|<value%d> %s}", index, index, index, summary)
		if indexID != 0 {
			// need pointer to value
			links = append(links, fmt.Sprintf("  %d:value%d -> %d:name;\n", id, index, indexID))
		}
	}
	node += "\"];\n"

	fmt.Fprint(m.writer, node)
	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	return id, m.nodeSummaries[key]
}

func (m *mapper) mapMap(mapVal reflect.Value, parentID nodeID, inlineable bool) (nodeID, string) {
	// create a string type while escaping graphviz special characters
	mapType := fmt.Sprint(mapVal.Type().String())
	mapType = strings.Replace(mapType, "[", "\\[", -1)
	mapType = strings.Replace(mapType, "]", "\\]", -1)

	nodeKey := getNodeKey(mapVal)

	if mapVal.Len() == 0 {
		m.nodeSummaries[nodeKey] = mapType + "{}"

		if inlineable {
			return 0, mapType
		}

		return m.newBasicNode(mapVal, mapType+"{}"), mapType
	}

	mapID := m.getNodeID(mapVal)
	var id nodeID
	if inlineable && mapVal.Len() <= m.inlineableItemLimit {
		m.nodeSummaries[nodeKey] = mapType
		id = parentID
	} else {
		id = mapID
	}

	var links []string
	var fields string
	for index, mapKey := range mapVal.MapKeys() {
		keyID, keySummary := m.mapValue(mapKey, id, true)
		valueID, valueSummary := m.mapValue(mapVal.MapIndex(mapKey), id, true)
		fields += fmt.Sprintf("|{<%dkey%d> %s| <%dvalue%d> %s}", mapID, index, keySummary, mapID, index, valueSummary)
		if keyID != 0 {
			links = append(links, fmt.Sprintf("  %d:<%dkey%d> -> %d:name;\n", id, mapID, index, keyID))
		}
		if valueID != 0 {
			links = append(links, fmt.Sprintf("  %d:<%dvalue%d> -> %d:name;\n", id, mapID, index, valueID))
		}
	}

	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	if inlineable && mapVal.Len() <= m.inlineableItemLimit {
		// inline map
		// remove stored summary so this gets regenerated every time
		// we need to do this so that we get a chance to print out the new links
		delete(m.nodeSummaries, nodeKey)

		// have to remove invalid leading |
		return 0, "{" + fields[1:] + "}"
	}

	// else create a new node
	node := fmt.Sprintf("  %d [label=\"<name> %s %s \"];\n", id, mapType, fields)
	fmt.Fprint(m.writer, node)

	return id, m.nodeSummaries[nodeKey]
}
