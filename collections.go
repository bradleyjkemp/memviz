package memviz

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func (m *mapper) mapStruct(structVal reflect.Value, depth uint32) (nodeID, string) {
	uType := structVal.Type()
	name := uType.Name()
	pkgPath := uType.PkgPath()
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
		if !m.config.includePrivateFields && !field.CanSet() {
			// ignore private field
			continue
		}
		// Get the field tag value to check for ignore
		f2 := uType.Field(index)
		tag := f2.Tag.Get(kTagName)
		tagElems := strings.Split(tag, ",")
		if _, ok := strArrContains(tagElems, "-", true); ok {
			// ignore element was requested
			continue
		}
		// access exported and unexported fields
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()

		// map next level
		fieldID, summary := m.mapValue(field, id, true, depth+1)

		// if field was inlined (id == 0) then print summary, else just the name and a link to the actual node
		if fieldID == 0 {
			fields += fmt.Sprintf("|{<f%d> %s | %s} ", index, f2.Name, summary)
		} else if fieldID > 0 {
			fields += fmt.Sprintf("|<f%d> %s", index, f2.Name)
			links = append(links, fmt.Sprintf("  %d:f%d -> %d:name;\n", id, index, fieldID))
		} else {
			// negative value means max depth was reached so don't link over to next node
		}
	}

	if !m.config.abbreviatedTypeNames && pkgPath != "" {
		name = pkgPath + "." + name
	}
	node := fmt.Sprintf("  %d [label=\"<name> %s %s \"];\n", id, name, fields)

	fmt.Fprint(m.writer, node)
	for _, link := range links {
		fmt.Fprint(m.writer, link)
	}

	return id, m.nodeSummaries[key]
}

func (m *mapper) mapSlice(sliceVal reflect.Value, parentID nodeID, inlineable bool, depth uint32) (nodeID, string) {
	sliceID := m.getNodeID(sliceVal)
	key := getNodeKey(sliceVal)
	st := sliceVal.Type()
	pkgPath := st.PkgPath()
	sliceType := escapeString(st.String())
	m.nodeSummaries[key] = sliceType

	if sliceVal.Len() == 0 {
		emptySlice := sliceType + "\\{\\}"
		if !m.config.abbreviatedTypeNames && pkgPath != "" {
			emptySlice = pkgPath + "." + emptySlice
		}
		m.nodeSummaries[key] = emptySlice

		if inlineable {
			return 0, m.nodeSummaries[key]
		}

		return m.newBasicNode(sliceVal, m.nodeSummaries[key]), sliceType
	}

	// sourceID is the nodeID that links will start from
	// if inlined then these come from the parent
	// if not inlined then these come from this node
	sourceID := sliceID
	if inlineable && (m.config.maxItemsToInline == 0 || sliceVal.Len() <= int(m.config.maxItemsToInline)) {
		sourceID = parentID
	}

	length := sliceVal.Len()
	var elements string
	var links []string
	for index := 0; index < length; index++ {
		indexID, summary := m.mapValue(sliceVal.Index(index), sliceID, true, depth+1)
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

	if inlineable && (m.config.maxItemsToInline == 0 || length <= int(m.config.maxItemsToInline)) {
		// inline slice
		// remove stored summary so this gets regenerated every time
		// we need to do this so that we get a chance to print out the new links
		delete(m.nodeSummaries, key)

		// have to remove invalid leading |
		return 0, "{" + elements[1:] + "}"
	}

	// else create a new node
	name := sliceType
	if !m.config.abbreviatedTypeNames && pkgPath != "" {
		name = pkgPath + "." + name
	}
	node := fmt.Sprintf("  %d [label=\"<name> %s %s \"];\n", sliceID, name, elements)
	fmt.Fprint(m.writer, node)

	return sliceID, m.nodeSummaries[key]
}

func (m *mapper) mapMap(mapVal reflect.Value, parentID nodeID, inlineable bool, depth uint32) (nodeID, string) {
	// create a string type while escaping graphviz special characters
	mt := mapVal.Type()
	mapType := escapeString(mt.String())
	pkgPath := mt.PkgPath()

	nodeKey := getNodeKey(mapVal)

	if mapVal.Len() == 0 {
		emptyMap := mapType + "\\{\\}"
		if !m.config.abbreviatedTypeNames && pkgPath != "" {
			emptyMap = pkgPath + "." + emptyMap
		}
		m.nodeSummaries[nodeKey] = emptyMap

		if inlineable {
			return 0, m.nodeSummaries[nodeKey]
		}

		return m.newBasicNode(mapVal, m.nodeSummaries[nodeKey]), mapType
	}

	mapID := m.getNodeID(mapVal)
	var id nodeID
	if inlineable && (m.config.maxItemsToInline == 0 || mapVal.Len() <= int(m.config.maxItemsToInline)) {
		m.nodeSummaries[nodeKey] = mapType
		id = parentID
	} else {
		id = mapID
	}

	var links []string
	var fields string
	for index, mapKey := range mapVal.MapKeys() {
		keyID, keySummary := m.mapValue(mapKey, id, true, depth+1)
		valueID, valueSummary := m.mapValue(mapVal.MapIndex(mapKey), id, true, depth+1)
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

	if inlineable && (m.config.maxItemsToInline == 0 || mapVal.Len() <= int(m.config.maxItemsToInline)) {
		// inline map
		// remove stored summary so this gets regenerated every time
		// we need to do this so that we get a chance to print out the new links
		delete(m.nodeSummaries, nodeKey)

		// have to remove invalid leading |
		return 0, "{" + fields[1:] + "}"
	}

	// else create a new node
	name := mapType
	if !m.config.abbreviatedTypeNames && pkgPath != "" {
		name = pkgPath + "." + name
	}
	node := fmt.Sprintf("  %d [label=\"<name> %s %s \"];\n", id, name, fields)
	fmt.Fprint(m.writer, node)

	return id, m.nodeSummaries[nodeKey]
}
