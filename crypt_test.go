// Copyright 2018 OVO Technology
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypt

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

	er := zerofill(path, false)
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

	er := secureDelete(path, false)
	check(er)

	if _, err := os.Stat(path); err == nil {
		t.Error("file still exists")
	}
}

// func TestPlainText(t *testing.T) {
// 	path := os.TempDir() + "plain.txt"
// 	ciphertextLength := 125
// 	b := make([]byte, ciphertextLength)
//
// 	s1 := base64.StdEncoding.EncodeToString(b)
//
// 	err := ioutil.WriteFile(path, []byte(s1), 0644)
// 	fmt.Println(s1)
// 	check(err)
// 	PlainText(path)
// }
