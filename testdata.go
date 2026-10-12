package pocketsmith

// This file contains common types and functions used across tests.

import (
	"fmt"
	"os"
)

// testdata represents a helper type that loads and stores file content for
// testing purposes.
type testdata struct {
	path    string
	content []byte
}

// getContent reads the content from the test file and stores it in the
// 'content' field. It panics if there is an error reading the file, and is
// not intended to be used in production code.
func (t *testdata) getContent() {
	b, err := os.ReadFile(t.path)
	if err != nil {
		panic(fmt.Sprintf("failed to read test data from %s: %v", t.path, err))
	}
	t.content = b
}

// newTestdata creates a new instance of testdata for the given 'fileName'. The
// path to the file is expected to be located under the './testdata' directory
// with a '.json' file extension.
func newTestdata(fileName string) *testdata {
	td := &testdata{
		path: "./testdata/" + fileName + ".json",
	}
	td.getContent()
	return td
}
