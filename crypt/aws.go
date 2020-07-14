package crypt

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type awsKms struct{}

func (a awsKms) encryptedDekLength() int {
	return 185
}

func (a awsKms) crypto(payload []byte, projectid, locationid, keyringid,
	cryptokeyid, keyname string, encrypt bool) (resultText []byte, err error) {

	svc := kms.New(session.New())
	if encrypt {
		resultText, err = awsKMSEncrypt(payload, cryptokeyid, svc)
	} else {
		resultText, err = awsKMSDecrypt(payload, svc)
	}
	return
}

func awsKMSEncrypt(payload []byte, cryptokeyid string, svc *kms.KMS) (resultText []byte, err error) {
	input := &kms.EncryptInput{
		KeyId:     aws.String(cryptokeyid),
		Plaintext: payload,
	}
	result, err := svc.Encrypt(input)
	if err == nil {
		resultText, err = base64.StdEncoding.DecodeString(string(result.CiphertextBlob))
	}
	return
}

func awsKMSDecrypt(payload []byte, svc *kms.KMS) (resultText []byte, err error) {
	input := &kms.DecryptInput{
		CiphertextBlob: payload,
	}
	result, err := svc.Decrypt(input)
	if err == nil {
		resultText = result.Plaintext
	}
	return
}
