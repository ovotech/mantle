package crypt

import (
	"context"
	"encoding/base64"
	"fmt"

	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type gcpKms struct{}

func (g gcpKms) encryptedDekLength() int {
	return 114
}

// uses google kms to either encrypt or decrypt a byte slice
func (g gcpKms) crypto(payload []byte, projectid, locationid, keyringid,
	cryptokeyid, keyname string, encrypt bool) (resultText []byte, err error) {
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
		resultText, err = googleKMSEncrypt(payload, parentName, kmsService)
	} else {
		resultText, err = googleKMSDecrypt(payload, parentName, kmsService)
	}
	return
}

//googleKMSEncrypt uses google kms to encypt a bite slice
func googleKMSEncrypt(payload []byte, parentName string,
	kmsService *cloudkms.Service) (resultText []byte, err error) {
	req := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString(payload),
	}
	var resp *cloudkms.EncryptResponse
	resp, err = kmsService.Projects.Locations.KeyRings.CryptoKeys.
		Encrypt(parentName, req).Do()
	check(err)
	var errm error
	resultText, errm = base64.StdEncoding.DecodeString(resp.Ciphertext)
	check(errm)
	return
}

//googleKMSDecrypt uses google kms to decypt a bite slice
func googleKMSDecrypt(payload []byte, parentName string,
	kmsService *cloudkms.Service) (resultText []byte, err error) {
	req := &cloudkms.DecryptRequest{
		Ciphertext: base64.StdEncoding.EncodeToString(payload),
	}
	var resp *cloudkms.DecryptResponse
	if resp, err = kmsService.Projects.Locations.KeyRings.CryptoKeys.
		Decrypt(parentName, req).Do(); err != nil {
		return
	}
	var errm error
	resultText, errm = base64.StdEncoding.DecodeString(resp.Plaintext)
	check(errm)
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
