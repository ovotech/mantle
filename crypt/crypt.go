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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type kmsProvider interface {
	crypto(payload []byte, projectid, locationid, keyringid,
		cryptokeyid, keyname string, encrypt bool) (resultText []byte, err error)
	encryptedDekLength() int
}

//Defaults type defining input flags
type Defaults struct {
	CryptoKeyID string `short:"c" long:"cryptokeyId" description:"Google KMS crytoKeyId" required:"false"`
	KeyRingID   string `short:"k" long:"keyringId" description:"Google KMS keyRingId" required:"false"`
	KeyName     string `short:"n" long:"keyName" description:"Google KMS keyName or AWS KMS keyId" required:"false"`
	LocationID  string `short:"l" long:"locationId" description:"Google KMS locationId" required:"false"`
	ProjectID   string `short:"p" long:"projectId" description:"Google projectId" required:"false"`
	KMSProvider string `short:"m" long:"kmsProvider" description:"KMS provider" required:"false"`
}

var (
	defaultOptions = Defaults{}
	//Parser is a new Parser with default options
	Parser       = flags.NewParser(&defaultOptions, flags.Default)
	kmsProviders = map[string]kmsProvider{
		"AWS": awsKms{},
		"GCP": gcpKms{},
	}
)

const (
	nonceLength = 12
	dekLength   = 32
)

func getKmsProvider(provider string) (kmsProvider kmsProvider, err error) {
	if provider == "" {
		return gcpKms{}, nil
	}
	kmsProvider, ok := kmsProviders[strings.ToUpper(provider)]
	if !ok {
		err = fmt.Errorf("KMS Provider %v not supported", provider)
	}
	return
}

//check panics if error is not nil
func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}

//byteSliceToString converts a byte slice to a string, and returns it
func byteSliceToString(dat []byte) (resultString string) {
	resultString = fmt.Sprint(string(dat[:]))
	return
}

//randByteSlice creates and returns a random byte slice, of desired size
func randByteSlice(size int) (bytes []byte) {
	bytes = make([]byte, size)
	_, err := io.ReadFull(rand.Reader, bytes)
	check(err)
	return
}

//secureDelete zerofills the desired file, and removes it
func secureDelete(filepath string, stdOut bool) (err error) {
	zerofill(filepath, stdOut)
	deleteFile(filepath)
	return
}

//zerofill zerofills the desired file
func zerofill(filepath string, stdOut bool) (err error) {
	fi, err := os.Stat(filepath)
	check(err)
	switch mode := fi.Mode(); {
	case mode.IsDir():
		if !stdOut {
			fmt.Printf("%s\n", "Didn't zerofill/delete unencrypted file \""+
				filepath+"\" as it's not a file")
		}
	case mode.IsRegular():
		file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
		check(err)
		defer file.Close()
		fileInfo, err := file.Stat()
		check(err)
		zeroBytes := make([]byte, fileInfo.Size())
		n, err := file.Write(zeroBytes)
		check(err)
		if !stdOut {
			fmt.Printf("Wiped %v bytes from %s.\n", n, filepath)
		}
	}
	return
}

//delete file removes the file
func deleteFile(filepath string) (err error) {
	err = os.Remove(filepath)
	check(err)
	return
}

//cipherblock creates and returns a new aes cipher.Block
func cipherblock(dek []byte) (cipherblock cipher.Block) {
	cipherblock, err := aes.NewCipher(dek)
	check(err)
	return
}

//cipherText seals or opens the text
func cipherText(text []byte, cipherblock cipher.Block, nonce []byte,
	seal bool) (ciphertext []byte) {
	aesgcm, err := cipher.NewGCM(cipherblock)
	check(err)
	var errm error
	if seal {
		ciphertext = aesgcm.Seal(nil, nonce, text, nil)
	} else {
		ciphertext, errm = aesgcm.Open(nil, nonce, text, nil)
	}
	check(errm)
	return
}
