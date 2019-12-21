package memviz

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// writeDotStringToPng takes a content of a dot file in a string and makes a graph using Graphviz
// to ouput an image
//
// sourced from https://github.com/Arafatk/DataViz/blob/master/utils/utils.go
func writeDotStringToPng(fileName string, dotFileString string) (ok bool) {
	byteString := []byte(dotFileString) // Converting the string to byte slice to write to a file
	tmpFile, _ := ioutil.TempFile("", "TemporaryDotFile")
	tmpFile.Write(byteString)            // Writing the string to a temporary file
	dotPath, err := exec.LookPath("dot") // Looking for dot command
	if err != nil {
		fmt.Printf("Error: Running the Visualizer command. Please install Graphviz. %s", err)
		return false
	}
	dotCommandResult, err := exec.Command(dotPath, "-Tpng", tmpFile.Name()).Output() // Running the command
	if err != nil {
		fmt.Printf("Error: Running the Visualizer command. Please install Graphviz. %s", err)
		return false
	}
	ioutil.WriteFile(fileName, dotCommandResult, os.FileMode(int(0777)))
	return true
}

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
