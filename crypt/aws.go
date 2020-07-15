// Copyright 2020 OVO Technology
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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type awsKms struct{}

func (a awsKms) encryptedDekLength() int {
	return 185
}

// uses aws kms to either encrypt or decrypt a byte slice
func (a awsKms) crypto(payload []byte, projectid, locationid, keyringid,
	cryptokeyid, keyname string, encrypt bool) (resultText []byte, err error) {

	svc := kms.New(session.New(&aws.Config{
		Region: aws.String("eu-west-1")}))
	if encrypt {
		resultText, err = awsKMSEncrypt(payload, keyname, svc)
	} else {
		resultText, err = awsKMSDecrypt(payload, svc)
	}
	return
}

//awsKMSEncrypt uses aws kms to encypt a bite slice
func awsKMSEncrypt(payload []byte, keyname string, svc *kms.KMS) (resultText []byte, err error) {
	input := &kms.EncryptInput{
		KeyId:     aws.String(keyname),
		Plaintext: payload,
	}
	result, err := svc.Encrypt(input)
	if err == nil {
		resultText = result.CiphertextBlob
	}
	return
}

//awsKMSDecrypt uses aws kms to decypt a bite slice
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
