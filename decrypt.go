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
	Filepath      string `short:"f" long:"filepath" description:"Path of file to get encrypted string from" default:"./cipher.txt"`
	Validate      bool   `short:"v" long:"validate" description:"Validate decryption works; don't produce a plain.txt"`
	WriteToStdout bool   `short:"o" long:"stdout" description:"Writes decrypted plaintext to console"`
}

var decryptCommand DecryptCommand

//Execute executes the DecryptCommand:
// 1. Obtains encrypted DEK from encrypted file
// 2. Decrypts DEK using KMS
// 3. Decrypts encrypted string from file using decrypted DEK
// 4. Outputs decrypted result to file
func (x *DecryptCommand) Execute(args []string) error {
	if !x.WriteToStdout {
		fmt.Println("Decrypting...")
	}
	var err error
	file, err := os.Open(x.Filepath)
	check(err)
	defer file.Close()
	s := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for s.Scan() {
		buffer.WriteString(s.Text())
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(buffer.String())
	check(err)
	dekLength := 113
	cipherLength := len(cipherBytes)
	encrypt := false
	outputFilepath := "./plain.txt"
	fileMode := os.FileMode.Perm(0644)
	encryptedDek := cipherBytes[cipherLength-dekLength : cipherLength]
	nonce := cipherBytes[cipherLength-(dekLength+nonceLength) : cipherLength-dekLength]
	decryptedDek := googleKMSCrypto(encryptedDek, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, defaultOptions.KeyName, encrypt)
	plainText := cipherText(cipherBytes[0:len(cipherBytes)-(dekLength+nonceLength)],
		cipherblock(decryptedDek), nonce, encrypt)
	if x.Validate {
		fmt.Println("Validation completed successfully")
		os.Exit(0)
	}
	if x.WriteToStdout {
		fmt.Printf("%s\n", plainText)
	} else {
		err = ioutil.WriteFile(outputFilepath, plainText, fileMode)
		check(err)
		fmt.Printf("Decryption successful, plaintext available at %s\n",
			outputFilepath)
	}
	check(secureDelete(x.Filepath))
	return err
}
