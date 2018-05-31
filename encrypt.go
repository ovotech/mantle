package main

import (
	"bufio"
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
			"spits out (data + encrypted DEK).",
		&encryptCommand)
}

//EncryptCommand type
type EncryptCommand struct {
	Filepath string `short:"f" long:"filepath" description:"Path of file to encrypt" default:"./plain.txt"`
	Nonce    string `short:"n" long:"nonce" description:"Nonce for encryption" required:"true"`
}

var encryptCommand EncryptCommand

//Execute executes the EncryptCommand:
// 1. Create new DEK
// 2. Encrypt data with the DEK
// 3. Encrypt DEK using KMS
// 4. Append encrypted DEK to encrypted data
// 5. Print out result to command-line, and to file
func (x *EncryptCommand) Execute(args []string) error {
	fmt.Println("Encrypting...")
	dek := randByteSlice(32)
	dat, err := ioutil.ReadFile(x.Filepath)
	check(err)
	encryptedDek := googleKMSCrypto(dek, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, true)
	cipherTexts := base64.StdEncoding.EncodeToString(append(
		cipherText(dat, cipherblock(dek), []byte(x.Nonce), true),
		encryptedDek...))
	fmt.Println("-----BEGIN (DATA + ENCRYPTED DEK) STRING-----")
	fmt.Println(cipherTexts)
	fmt.Println("-----END (DATA + ENCRYPTED DEK) STRING-----")
	f, err := os.Create("cipher.txt")
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = fmt.Fprint(w, cipherTexts)
	check(err)
	w.Flush()
	return nil
}
