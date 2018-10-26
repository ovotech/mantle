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
	"fmt"
)

func init() {
	parser.AddCommand("reencrypt",
		"Decrypts encrypted text, returning the plaintext data",
		"Decrypts the encrypted DEK via KMS, decrypts the data with the DEK, "+
			"outputs to file",
		&reencryptCommand)
}

//ReencryptCommand type
type ReencryptCommand struct {
	Filepath   string `short:"f" long:"filepath" description:"Path of file to get encrypted string from" default:"./cipher.txt"`
	SingleLine bool   `short:"s" long:"singleLine" description:"Disable use of newline chars in ciphertext"`
}

var reencryptCommand ReencryptCommand

//Execute executes the ReencryptCommand
func (x *ReencryptCommand) Execute(args []string) error {
	fmt.Println("Reencrypting...")
	return Reencrypt(x.Filepath, x.SingleLine)
}

//Reencrypt decrypts into a plaintext byte array, and encrypts back to ciphertext file
func Reencrypt(filepath string, singleLine bool) error {
	plaintext, err := PlainText(filepath)
	check(err)
	err = CipherText(plaintext, filepath, singleLine)
	return err
}
