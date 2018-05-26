package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type Defaults struct {
	ProjectId   string `short:"p" long:"projectId" description:"Show verbose debug information" required:"true"`
	LocationId  string `short:"l" long:"locationId" description:"Shows terse output" required:"true"`
	KeyRingId   string `short:"k" long:"keyringId" description:"Shows terse output" required:"true"`
	CryptoKeyId string `short:"c" long:"cryptokeyId" description:"Shows terse output" required:"true"`
}

func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func googleKMSCrypto(payload []byte, projectid, locationid, keyringid,
	cryptokeyid string, encrypt bool) (resultText []byte) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	// Create the KMS client.
	kmsService, err := cloudkms.New(client)
	if err != nil {
		log.Fatal(err)
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectid, locationid, keyringid, cryptokeyid)

	//ar errm error
	if encrypt {
		req := &cloudkms.EncryptRequest{
			Plaintext: base64.StdEncoding.EncodeToString(payload),
		}
		resp, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(parentName, req).Do()
		check(err)
		var errm error
		resultText, errm = base64.StdEncoding.DecodeString(resp.Ciphertext)
		check(errm)
	} else {
		fmt.Printf("decrypt payload: %s", payload)
		req := &cloudkms.DecryptRequest{
			Ciphertext: base64.StdEncoding.EncodeToString(payload),
		}
		resp, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(parentName, req).Do()
		check(err)
		var errm error
		resultText, errm = base64.StdEncoding.DecodeString(resp.Plaintext)
		check(errm)
	}
	if err != nil {
		log.Fatal(err)
	}
	return
}

func dek() (dek []byte) {
	dek = make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, dek)
	check(err)
	return
}

func nonce() (nonce []byte) {
	//Never use more than 2^32 random nonces with a given key because of the risk
	// of a repeat.
	nonce = make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, nonce)
	check(err)
	return
}

func validateInput(defaultOptions Defaults) {
	//TODO: validate the input flags, panic if something's not right
}

func main() {
	defaultOptions := Defaults{}
	parser := flags.NewParser(&defaultOptions, flags.Default)
	_, err := parser.Parse()
	check(err)
	validateInput(defaultOptions)
	dek := dek()
	fmt.Printf("dek: %s", dek)
	nonce := nonce()

	//encrypt data using aes-256-gcm
	cipherTexts := cipherText([]byte("exampleplaintext"), cipherblock(dek), nonce, true)

	//encrypt DEK
	encryptedDek := googleKMSCrypto(dek, defaultOptions.ProjectId,
		defaultOptions.LocationId, defaultOptions.KeyRingId,
		defaultOptions.CryptoKeyId, true)
	fmt.Printf("encrypted dek: %s", encryptedDek)

	decryptedDek := googleKMSCrypto(encryptedDek, defaultOptions.ProjectId,
		defaultOptions.LocationId, defaultOptions.KeyRingId,
		defaultOptions.CryptoKeyId, false)
	fmt.Printf("decrypted dek: %s", decryptedDek)

	plainText := cipherText(cipherTexts, cipherblock(decryptedDek), nonce, false)
	fmt.Printf("plaintext: %s", plainText)

	// //TODO: shred plaintext file
}

func cipherblock(dek []byte) (cipherblock cipher.Block) {
	cipherblock, err := aes.NewCipher(dek)
	check(err)
	return
}

func cipherText(text []byte, cipherblock cipher.Block, nonce []byte, seal bool) (ciphertext []byte) {
	//TODO: get this from file or input flag
	//plaintext := []byte("exampleplaintext")

	aesgcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		panic(err.Error())
	}
	var errm error
	if seal {
		ciphertext = aesgcm.Seal(nil, nonce, text, nil)
	} else {
		ciphertext, errm = aesgcm.Open(nil, nonce, text, nil)
	}
	check(errm)

	return
}
