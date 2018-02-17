package memviz

import (
	"fmt"
	"reflect"
	"unsafe"
)

func (m *mapper) mapStruct(structVal reflect.Value) (nodeID, string) {
	uType := structVal.Type()
	id := m.getNodeID(structVal)
	key := getNodeKey(structVal)
	m.nodeSummaries[key] = escapeString(uType.String())

	var fields string
	var links []string
	for index := 0; index < uType.NumField(); index++ {
		field := structVal.Field(index)
		if !field.CanAddr() {
			// TODO: when does this happen? Can we work around it?
			continue
		}
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
	sliceID := m.getNodeID(sliceVal)
	key := getNodeKey(sliceVal)
	sliceType := escapeString(sliceVal.Type().String())
	m.nodeSummaries[key] = sliceType

	if sliceVal.Len() == 0 {
		m.nodeSummaries[key] = sliceType + "\\{\\}"

		if inlineable {
			return 0, m.nodeSummaries[key]
		}

		return m.newBasicNode(sliceVal, m.nodeSummaries[key]), sliceType
	}

	// sourceID is the nodeID that links will start from
	// if inlined then these come from the parent
	// if not inlined then these come from this node
	sourceID := sliceID
	if inlineable && sliceVal.Len() <= m.inlineableItemLimit {
		sourceID = parentID
	}

	length := sliceVal.Len()
	var elements string
	var links []string
	for index := 0; index < length; index++ {
		indexID, summary := m.mapValue(sliceVal.Index(index), sliceID, true)
		if indexID != 0 {
			// need pointer to value
			elements += fmt.Sprintf("|<%dindex%d> %d", sliceID, index, index)
			links = append(links, fmt.Sprintf("  %d:<%dindex%d> -> %d:name;\n", sourceID, sliceID, index, indexID))
		} else {
			// field was inlined so print summary
			elements += fmt.Sprintf("|{<%dindex%d> %d|<%dvalue%d> %s}", sliceID, index, index, sliceID, index, summary)
		}
	}

	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	if inlineable && length <= m.inlineableItemLimit {
		// inline slice
		// remove stored summary so this gets regenerated every time
		// we need to do this so that we get a chance to print out the new links
		delete(m.nodeSummaries, key)

		// have to remove invalid leading |
		return 0, "{" + elements[1:] + "}"
	}

	// else create a new node
	node := fmt.Sprintf("  %d [label=\"<name> %s %s \"];\n", sliceID, sliceType, elements)
	fmt.Fprint(m.writer, node)

	return sliceID, m.nodeSummaries[key]
}

func (m *mapper) mapMap(mapVal reflect.Value, parentID nodeID, inlineable bool) (nodeID, string) {
	// create a string type while escaping graphviz special characters
	mapType := escapeString(mapVal.Type().String())

	nodeKey := getNodeKey(mapVal)

	if mapVal.Len() == 0 {
		m.nodeSummaries[nodeKey] = mapType + "\\{\\}"

		if inlineable {
			return 0, m.nodeSummaries[nodeKey]
		}

		return m.newBasicNode(mapVal, m.nodeSummaries[nodeKey]), mapType
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
