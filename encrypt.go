package main

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

//Execute executes the EncryptCommand:
// 1. Create new DEK
// 2. Encrypt data with the DEK
// 3. Encrypt DEK using KMS
// 4. Append encrypted DEK to encrypted data
// 5. Print out result to command-line, and to file
func (x *EncryptCommand) Execute(args []string) error {
	fmt.Println("Encrypting...")
	dekSize := 32
	dek := randByteSlice(dekSize)
	nonce := randByteSlice(nonceLength)
	dat, err := ioutil.ReadFile(x.Filepath)
	check(err)
	encrypt := true
	outputFilepath := "./cipher.txt"
	fileMode := os.FileMode.Perm(0644)
	encryptedDek := googleKMSCrypto(dek, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, encrypt)
	cipherBytes := []byte(base64.StdEncoding.EncodeToString(append(
		append(cipherText(dat, cipherblock(dek), nonce, encrypt),
			nonce...),
		encryptedDek...)))
	var cipherTexts []byte
	if x.SingleLine {
		cipherTexts = cipherBytes
	} else {
		cipherTexts = insertNewLines(cipherBytes)
	}
	fmt.Println("-----BEGIN (ENCRYPTED DATA + DEK) STRING-----")
	fmt.Printf("%s\n", cipherTexts)
	fmt.Println("-----END (ENCRYPTED DATA + DEK) STRING-----")
	ioutil.WriteFile(outputFilepath, cipherTexts, fileMode)
	fmt.Printf("Encryption successful, ciphertext available at %s\n",
		outputFilepath)
	check(zerofill(x.Filepath))
	return nil
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
