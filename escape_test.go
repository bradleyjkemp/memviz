package memviz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	input  string
	output string
}{
	{"Hello world", "Hello world"},

	// double quotes are escaped
	{"\"Hello world\"", "\\\"Hello world\\\""},

	// brackets not escaped
	{"map[string]bool", "map[string]bool"},

	// braces escaped
	{"map[string]struct{}", "map[string]struct\\{\\}"},
}

func TestEscapeString(t *testing.T) {
	for _, tc := range cases {
		require.Equal(t, tc.output, escapeString(tc.input))
	}
}
