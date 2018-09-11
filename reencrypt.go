package main

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

//Execute executes the ReencryptCommand:
func (x *ReencryptCommand) Execute(args []string) error {
	fmt.Println("Reencrypting...")
	plaintext, err := PlainText(x.Filepath)
	check(err)
	err = CipherText(plaintext, x.Filepath, x.SingleLine)
	return err
}
