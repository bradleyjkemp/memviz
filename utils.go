package memviz

import (
	"strings"
)

// Finds a string in a string array
//
// sarr = the string array to search in
// s = the string to search
// cs = true/false for case sensitive search
//
// Returns the index of the item and true/false whether it's found in the array
func strArrContains(sarr []string, s string, cs bool) (int, bool) {
	lks := strings.ToLower(s)
	for idx, s2 := range sarr {
		if (cs && s == s2) || (!cs && strings.ToLower(s2) == lks) {
			return idx, true
		}
	}
	return -1, false
}
