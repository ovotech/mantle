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
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func init() {
	Parser.AddCommand("decrypt",
		"Decrypts encrypted text, returning the plaintext data",
		"Decrypts the encrypted DEK via KMS, decrypts the data with the DEK, "+
			"outputs to file",
		&decryptCommand)
}

//DecryptCommand type
type DecryptCommand struct {
	Filepath         string `short:"f" long:"filepath" description:"Path of file to get encrypted string from" default:"./cipher.txt"`
	RetainCipherText bool   `short:"r" long:"retainCipherText" description:"Retain ciphertext after decryption"`
	TargetFilepath   string `short:"t" long:"targetFilepath" description:"Path of file to write decrypted string to" default:"./plain.txt"`
	Validate         bool   `short:"v" long:"validate" description:"Validate decryption works"`
	WriteToStdout    bool   `short:"o" long:"stdout" description:"Writes decrypted plaintext to console"`
}

var decryptCommand DecryptCommand

//Execute executes the DecryptCommand
func (x *DecryptCommand) Execute(args []string) error {
	if !x.WriteToStdout {
		fmt.Println("Decrypting...")
	}
	plaintext, err := PlainText(x.Filepath)
	outputFilepath := x.TargetFilepath
	fileMode := os.FileMode.Perm(0644)
	if x.Validate {
		fmt.Println("Validation completed successfully")
		os.Exit(0)
	}
	if x.WriteToStdout {
		fmt.Printf("%s\n", plaintext)
	} else {
		err = ioutil.WriteFile(outputFilepath, plaintext, fileMode)
		check(err)
		fmt.Printf("Decryption successful, plaintext available at %s\n",
			outputFilepath)
	}
	if !x.RetainCipherText {
		check(secureDelete(x.Filepath, x.WriteToStdout))
	}
	return err
}

func checkCipherTextLength(ciphertext []byte) {
	length := len(ciphertext)
	minLength := encDekLength + nonceLength
	if length < minLength {
		panic("CipherText was shorter (" + strconv.Itoa(length) +
			") than the smallest possible generated CipherText (" +
			strconv.Itoa(minLength) + ")")
	}
}

// PlainText returns a slice of bytes (the plaintext), decrypted from File
func PlainText(filepath string) (plaintext []byte, err error) {
	file, err := os.Open(filepath)
	check(err)
	defer file.Close()
	s := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for s.Scan() {
		buffer.WriteString(s.Text())
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(buffer.String())
	check(err)
	plaintext, err = PlainTextFromBytes(cipherBytes)
	return
}

// PlainTextFromBytes returns a slice of bytes (the plaintext), decrypted from
// a byte slice
func PlainTextFromBytes(cipherBytes []byte) (plaintext []byte, err error) {
	return PlainTextFromPrimitives(cipherBytes, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, defaultOptions.KeyName)
}

// PlainTextFromPrimitives returns a slice of bytes (the plaintext), decrypted from
// a byte slice
func PlainTextFromPrimitives(cipherBytes []byte, projectID, locationID, keyRingID,
	cryptoKeyID, keyName string) (plaintext []byte, err error) {
	checkCipherTextLength(cipherBytes)
	cipherLength := len(cipherBytes)
	encrypt := false
  if plaintext, err = plainTextWithDekLength(cipherBytes, projectID, locationID, keyRingID,
		cryptoKeyID, keyName, encDekLength, cipherLength, encrypt); err != nil{
		plaintext, err = plainTextWithDekLength(cipherBytes, projectID, locationID, keyRingID,
			cryptoKeyID, keyName, encDekLength - 1, cipherLength, encrypt)
	}
	return
}

func plainTextWithDekLength(cipherBytes []byte, projectID, locationID, keyRingID,
	cryptoKeyID, keyName string, encDekLength, cipherLength int, encrypt bool)(plaintext []byte, err error){
	encryptedDek := cipherBytes[cipherLength-encDekLength : cipherLength]
	nonce := cipherBytes[cipherLength-(encDekLength+nonceLength) : cipherLength-encDekLength]
	if decryptedDek, err := googleKMSCrypto(encryptedDek, projectID,
		locationID, keyRingID, cryptoKeyID, keyName, encrypt); err == nil{
			plaintext = cipherText(cipherBytes[0:len(cipherBytes)-(encDekLength+nonceLength)],
				cipherblock(decryptedDek), nonce, encrypt)
		}
	return
}
