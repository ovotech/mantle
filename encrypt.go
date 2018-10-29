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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	//defaultOptions := Defaults{}
	//parser := flags.NewParser(&defaultOptions, flags.Default)
	parser.AddCommand("encrypt",
		"Encrypts your data, returning everything required for future decryption",
		"Creates a new DEK, encrypts data with DEK, encrypts the DEK using KMS, "+
			"spits out encrypted data + encrypted DEK.",
		&encryptCommand)
}

//EncryptCommand type
type EncryptCommand struct {
	Filepath   string `short:"f" long:"filepath" description:"Path of file to encrypt" default:"./plain.txt"`
	SingleLine bool   `short:"s" long:"singleLine" description:"Disable use of newline chars in ciphertext"`
}

var encryptCommand EncryptCommand

//Execute executes the EncryptCommand
func (x *EncryptCommand) Execute(args []string) (err error) {
	fmt.Println("Encrypting...")
	dat, err := ioutil.ReadFile(x.Filepath)
	check(err)
	err = CipherText(dat, x.Filepath, x.SingleLine)
	check(secureDelete(x.Filepath, false))
	return err
}

// insertNewLines inserts a newline char at specific intervals
func insertNewLines(cipherTexts []byte) (newLineText []byte) {
	interval := 40
	for i, char := range cipherTexts {
		if i > 0 && (i%interval == 0) {
			newLineText = append(newLineText, []byte("\n")...)
		}
		newLineText = append(newLineText, char)
	}
	return
}

// CipherText creates a ciphertext encrypted from a slice of bytes
// (the plaintext), and writes to File and Console.
func CipherText(plaintext []byte, filepath string, singleLine bool) (err error) {
	dekSize := 32
	dek := randByteSlice(dekSize)
	nonce := randByteSlice(nonceLength)
	check(err)
	encrypt := true
	encryptedDek := googleKMSCrypto(dek, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, defaultOptions.KeyName, encrypt)
	cipherBytes := []byte(base64.StdEncoding.EncodeToString(append(
		append(cipherText(plaintext, cipherblock(dek), nonce, encrypt),
			nonce...),
		encryptedDek...)))
	check(err)
	var cipherTexts []byte
	if singleLine {
		cipherTexts = cipherBytes
	} else {
		cipherTexts = insertNewLines(cipherBytes)
	}
	outputFilepath := "./cipher.txt"
	fileMode := os.FileMode.Perm(0644)
	fmt.Println("-----BEGIN (ENCRYPTED DATA + DEK) STRING-----")
	fmt.Printf("%s\n", cipherTexts)
	fmt.Println("-----END (ENCRYPTED DATA + DEK) STRING-----")
	ioutil.WriteFile(outputFilepath, cipherTexts, fileMode)
	fmt.Printf("Encryption successful, ciphertext available at %s\n",
		outputFilepath)
	return
}
