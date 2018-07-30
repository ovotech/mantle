package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestZerofill(t *testing.T) {

	path := os.TempDir() + "/testFile"

	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile(path, d1, 0644)
	check(err)

	er := zerofill(path)
	check(er)

	if _, err := os.Stat(path); err == nil {
		t.Error("file still exists")
	}

}
