package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
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

func googleKMSEncrypt(dek []byte, projectid, locationid, keyringid,
	cryptokeyid string) (encrypteddek string) {
	plaintext := "plaintextexample"
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
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
		projectid, locationid, projectid, cryptokeyid)

	req := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}
	resp, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(parentName, req).Do()
	if err != nil {
		log.Fatal(err)
	}
	encrypteddek = resp.Ciphertext
	return
	//fmt.Printf("returned from KMS: %x\n", resp.Ciphertext)
}

func dek() (dek []byte) {
	dek = make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, dek)
	check(err)
	return
}

func nonce() (nonce []byte) {
	//Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce = make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, nonce)
	check(err)
	return
}

func main() {

	defaultOptions := Defaults{}
	parser := flags.NewParser(&defaultOptions, flags.Default)
	parser.Parse()
	fmt.Printf("Verbose: %v\n", defaultOptions.ProjectId)
	fmt.Printf("Terse: %v\n", defaultOptions.LocationId)
	fmt.Printf("Terse: %v\n", defaultOptions.KeyRingId)
	fmt.Printf("Terse: %v\n", defaultOptions.CryptoKeyId)

	dek := dek()
	nonce := nonce()
	cipherblock := cipherblock(dek)
	fmt.Println(encryptedText(cipherblock, nonce))
	fmt.Println(googleKMSEncrypt(dek, projectid, locationid, keyringid,
		cryptokeyid))
	// //TODO: shred plaintext file
}

func cipherblock(dek []byte) (cipherblock cipher.Block) {
	cipherblock, err := aes.NewCipher(dek)
	check(err)
	return
}

func encryptedText(cipherblock cipher.Block, nonce []byte) (ciphertext []byte) {
	plaintext := []byte("exampleplaintext")

	aesgcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		panic(err.Error())
	}
	ciphertext = aesgcm.Seal(nil, nonce, plaintext, nil)
	return
}

func ExampleNewGCMDecrypter() {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	key := []byte("AES256Key-32Characters1234567890")
	ciphertext, _ := hex.DecodeString("***REMOVED***")

	nonce, _ := hex.DecodeString("***REMOVED***")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n", string(plaintext))
}
