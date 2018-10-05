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

package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

//Defaults type defining input flags
type Defaults struct {
	CryptoKeyID string `short:"c" long:"cryptokeyId" description:"Google kms crytoKeyId" required:"false"`
	KeyRingID   string `short:"k" long:"keyringId" description:"Google kms keyRingId" required:"false"`
	KeyName     string `short:"n" long:"keyName" description:"Google kms keyName" required:"false"`
	LocationID  string `short:"l" long:"locationId" description:"Google kms locationId" required:"false"`
	ProjectID   string `short:"p" long:"projectId" description:"Google projectId" required:"false"`
}

var defaultOptions = Defaults{}

var parser = flags.NewParser(&defaultOptions, flags.Default)

var nonceLength = 12

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

//kmsClient returns a kms service created from a default google client
func kmsClient() (kmsService *cloudkms.Service) {
	ctx := context.Background()
	client, errc := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	check(errc)
	kmsService, errk := cloudkms.New(client)
	check(errk)
	return
}

//googleKMSCrypto uses google kms to either encrypt or decrypt a byte slice
func googleKMSCrypto(payload []byte, projectid, locationid, keyringid,
	cryptokeyid, keyname string, encrypt bool) (resultText []byte) {
	kmsService := kmsClient()
	var parentName string
	if len(keyname) > 0 {
		parentName = keyname
	} else {
		parentName = fmt.Sprintf(
			"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectid,
			locationid, keyringid, cryptokeyid)
	}
	if encrypt {
		resultText = googleKMSEncrypt(payload, parentName, kmsService)
	} else {
		resultText = googleKMSDecrypt(payload, parentName, kmsService)
	}
	return
}

//googleKMSEncrypt uses google kms to encypt a bite slice
func googleKMSEncrypt(payload []byte, parentName string,
	kmsService *cloudkms.Service) (resultText []byte) {
	req := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString(payload),
	}
	resp, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.
		Encrypt(parentName, req).Do()
	check(err)
	var errm error
	resultText, errm = base64.StdEncoding.DecodeString(resp.Ciphertext)
	check(errm)
	return
}

//googleKMSDecrypt uses google kms to decypt a bite slice
func googleKMSDecrypt(payload []byte, parentName string,
	kmsService *cloudkms.Service) (resultText []byte) {
	req := &cloudkms.DecryptRequest{
		Ciphertext: base64.StdEncoding.EncodeToString(payload),
	}
	resp, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.
		Decrypt(parentName, req).Do()
	check(err)
	var errm error
	resultText, errm = base64.StdEncoding.DecodeString(resp.Plaintext)
	check(errm)
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

func main() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
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
