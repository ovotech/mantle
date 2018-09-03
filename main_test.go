package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestZerofill(t *testing.T) {
	path := os.TempDir() + "testFile"
	defer os.Remove(path)

	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile(path, d1, 0644)
	check(err)

	er := zerofill(path)
	check(er)

	dat, err := ioutil.ReadFile(path)
	check(err)

	zerod := make([]byte, len(d1))
	if !reflect.DeepEqual(zerod, dat) {
		t.Error("data should be zero'ed array")
	}
}

func TestDeleteFile(t *testing.T) {
	path := os.TempDir() + "testFile"

	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile(path, d1, 0644)
	check(err)

	er := deleteFile(path)
	check(er)

	if _, err := os.Stat(path); err == nil {
		t.Error("file still exists")
	}
}

func TestSecureDelete(t *testing.T) {
	path := os.TempDir() + "testFile"

	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile(path, d1, 0644)
	check(err)

	er := secureDelete(path)
	check(er)

	if _, err := os.Stat(path); err == nil {
		t.Error("file still exists")
	}
}
