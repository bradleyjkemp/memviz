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

	var fields []string
	var links []string

	for index := 0; index < uType.NumField(); index++ {
		field := structVal.Field(index)
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		fieldID, summary := m.mapValue(field, true)

		// if field was inlined (id == 0) then print summary, else just the name and a link to the actual
		if fieldID == 0 {
			fields = append(fields, fmt.Sprintf("%s: %s", uType.Field(index).Name, summary))
		} else {
			fields = append(fields, uType.Field(index).Name)
			links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", id, index, fieldID))
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

	return id, m.nodeSummaries[key]
}

func (m *mapper) mapSlice(sliceVal reflect.Value) (nodeID, string) {
	id := m.getNodeID(sliceVal)
	key := getNodeKey(sliceVal)
	m.nodeSummaries[key] = sliceVal.Type().String()

	length := sliceVal.Len()
	node := fmt.Sprintf("  %d [label=\"<name> %s", id, sliceVal.Type().String())
	var links []string
	for index := 0; index < length; index++ {
		indexID, summary := m.mapValue(sliceVal.Index(index), true)
		if indexID == 0 {
			// field was inlined
			node += fmt.Sprintf("|<f%d> %d: %s", index, index, summary)
		} else {
			// need a link to the new node and don't care about the summary
			node += fmt.Sprintf("|<f%d> %d", index, index)
			links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", id, index, indexID))
		}
	}
	node += "\"];\n"

	fmt.Fprint(m.writer, node)
	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	return id, m.nodeSummaries[key]
}

func (m *mapper) mapMap(mapVal reflect.Value) (nodeID, string) {
	// create a string type while escaping graphviz special characters
	mapType := fmt.Sprint(mapVal.Type().String())
	mapType = strings.Replace(mapType, "[", "\\[", -1)
	mapType = strings.Replace(mapType, "]", "\\]", -1)

	if mapVal.Len() == 0 {
		return 0, mapType
	}

	id := m.getNodeID(mapVal)

	var links []string
	var fields string
	for index, key := range mapVal.MapKeys() {
		keyID, keySummary := m.mapValue(key, true)
		valueID, valueSummary := m.mapValue(mapVal.MapIndex(key), true)
		fields += fmt.Sprintf("|{<key%d> %s| <value%d> %s}", index, keySummary, index, valueSummary)
		if keyID != 0 {
			links = append(links, fmt.Sprintf("  %d:key%d -> %d:name;\n", id, index, keyID))
		}
		if valueID != 0 {
			links = append(links, fmt.Sprintf("  %d:value%d -> %d:name;\n", id, index, valueID))
		}
	}

	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	if mapVal.Len() <= 3 {
		// inline map
		// have to remove invalid leading |
		return 0, fields[1:]
	}

	// else create a new node
	node := fmt.Sprintf("  %d [label=\"<name> %s %s \"];\n", id, mapType, fields)
	fmt.Fprint(m.writer, node)

	key := getNodeKey(mapVal)
	m.nodeSummaries[key] = mapVal.Type().String()

	return id, m.nodeSummaries[key]
}
