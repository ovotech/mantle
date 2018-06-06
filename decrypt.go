package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	parser.AddCommand("decrypt",
		"Decrypts encrypted text, returning the plaintext data",
		"Decrypts the encrypted DEK via KMS, decrypts the data with the DEK, "+
			"outputs to file",
		&decryptCommand)
}

//DecryptCommand type
type DecryptCommand struct {
	Filepath string `short:"f" long:"filepath" description:"Path of file to get encrypted string from" default:"./cipher.txt"`
	Nonce    string `short:"n" long:"nonce" description:"Nonce for decryption" required:"true"`
}

var decryptCommand DecryptCommand

//Execute executes the DecryptCommand:
// 1. Obtains encrypted DEK from encrypted file
// 2. Decrypts DEK using KMS
// 3. Decrypts encrypted string from file using decrypted DEK
// 4. Outputs decrypted result to file
func (x *DecryptCommand) Execute(args []string) error {
	fmt.Println("Decrypting...")
	file, err := os.Open(x.Filepath)
	check(err)
	defer file.Close()
	s := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for s.Scan() {
		buffer.WriteString(s.Text())
	}
	cipherBytes, errb := base64.StdEncoding.DecodeString(buffer.String())
	check(errb)
	dekLength := 113
	encrypt := false
	outputFilepath := "./plain.txt"
	fileMode := os.FileMode.Perm(0644)
	encryptedDek := cipherBytes[len(cipherBytes)-dekLength : len(cipherBytes)]
	decryptedDek := googleKMSCrypto(encryptedDek, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, encrypt)
	plainText := cipherText(cipherBytes[0:len(cipherBytes)-dekLength],
		cipherblock(decryptedDek), []byte(x.Nonce), encrypt)
	ioutil.WriteFile(outputFilepath, plainText, fileMode)
	fmt.Printf("Decryption successful, plaintext available at %s\n",
		outputFilepath)
	check(zerofill(x.Filepath))
	return nil
}
