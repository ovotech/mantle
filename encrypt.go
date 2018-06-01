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
	cipherTexts := []byte(base64.StdEncoding.EncodeToString(append(
		cipherText(dat, cipherblock(dek), []byte(x.Nonce), true),
		encryptedDek...)))
	fmt.Println("-----BEGIN (DATA + ENCRYPTED DEK) STRING-----")
	fmt.Printf("%s\n", cipherTexts)
	fmt.Println("-----END (DATA + ENCRYPTED DEK) STRING-----")
	ioutil.WriteFile("cipher.txt", cipherTexts, 0644)
	check(zerofill(x.Filepath))
	return nil
}

//zerofill zerofills the desired file, and removes it
func zerofill(filepath string) (err error) {
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	check(err)
	defer file.Close()
	fileInfo, err := file.Stat()
	check(err)
	if fileInfo.IsDir() {
		fmt.Printf("%s\n", "Didn't zerofill/delete unencrypted file \""+
			filepath+"\" as it's not a file")
	} else {
		zeroBytes := make([]byte, fileInfo.Size())
		copy(zeroBytes[:], "0")
		n, err := file.Write([]byte(zeroBytes))
		check(err)
		fmt.Printf("Wiped %v bytes from %s.\n", n, filepath)
		err = os.Remove(filepath)
		check(err)
	}
	return
}
