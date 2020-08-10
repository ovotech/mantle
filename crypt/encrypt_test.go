package crypt

import (
	"testing"
)

var newLineTests = []struct {
	cipherText     string
	expectedResult string
}{
	{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\na"},
}

func TestInsertNewLines(t *testing.T) {
	for _, newLineTest := range newLineTests {
		newLineString := string(insertNewLines([]byte(newLineTest.cipherText)))
		expectedResult := newLineTest.expectedResult
		if newLineString != expectedResult {
			t.Errorf("Got %s, want %s", newLineString, expectedResult)
		}
	}
}
