package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type Defaults struct {
	ProjectID   string "short:\"p\" long:\"projectId\" description:\" Show verbose debug information\" required:\"true\""
	LocationID  string `short:"l" long:"locationId" description:"Shows terse output" required:"true"`
	KeyRingID   string `short:"k" long:"keyringId" description:"Shows terse output" required:"true"`
	CryptoKeyID string `short:"c" long:"cryptokeyId" description:"Shows terse output" required:"true"`
}

func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func kmsClient() (kmsService *cloudkms.Service) {
	ctx := context.Background()
	client, errc := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	check(errc)
	kmsService, errk := cloudkms.New(client)
	check(errk)
	return
}

func googleKMSCrypto(payload []byte, projectid, locationid, keyringid,
	cryptokeyid string, encrypt bool) (resultText []byte) {
	kmsService := kmsClient()
	parentName := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectid,
		locationid, keyringid, cryptokeyid)
	if encrypt {
		resultText = googleKMSEncrypt(payload, parentName, kmsService)

	} else {
		resultText = googleKMSDecrypt(payload, parentName, kmsService)
	}
	return
}

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

func googleKMSDecrypt(payload []byte, parentName string,
	kmsService *cloudkms.Service) (resultText []byte) {
	fmt.Printf("decrypt payload: %x\n", payload)
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
	fmt.Printf("%x\n", dek)
	nonce := nonce()

	//encrypt data using aes-256-gcm
	cipherTexts := cipherText([]byte("exampleplaintext"), cipherblock(dek),
		nonce, true)

	//encrypt DEK
	encryptedDek := googleKMSCrypto(dek, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, true)
	fmt.Printf("encrypted dek: %x\n", encryptedDek)

	decryptedDek := googleKMSCrypto(encryptedDek, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, false)
	fmt.Printf("decrypted dek: %x\n", decryptedDek)

	plainText := cipherText(cipherTexts, cipherblock(decryptedDek), nonce, false)
	fmt.Printf("plaintext: %s\n", plainText)

	// //TODO: shred plaintext file
}

func cipherblock(dek []byte) (cipherblock cipher.Block) {
	cipherblock, err := aes.NewCipher(dek)
	check(err)
	return
}

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
