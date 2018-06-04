# aes-256-gcm-kms

## intro

A Go program to simplify the encryption & decryption of strings, using 256 bit AES keys in Galois/Counter Mode (GCM), with cloud-based KMS services and multiple key layers (specifically [Envelope Encryption](https://cloud.google.com/kms/docs/envelope-encryption))

## how-to decrypt

pre-requisites:
* `google credentials` on the host the aes-256-gcm-kms binary is run on, see [here](https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials), that has kms decrypt permission
* kms `cryptoKeyId`, `keyRingId`, `locationId` and `projectId` values
* the `ciphertext` you want to decrypt in a file on the host you're invoking the binary from (defaults to ./cipher.txt)
* a `nonce` string to be used in the [AEAD](https://golang.org/pkg/crypto/cipher/#AEAD) cipher `open` call, must match the nonce used for encryption

`aes-256-gcm-kms decrypt` will:
* split ciphertext into encrypted data + encrypted DEK
* decrypts the encrypted DEK using KMS
* decrypts the encrypted data
* outputs to ./plain.txt

```
$ aes-256-gcm-kms decrypt -h

Usage:
  aes-256-gcm-kms [OPTIONS] decrypt [decrypt-OPTIONS]

Decrypts the encrypted DEK via KMS, decrypts the data with the DEK, outputs to file

Application Options:
  -c, --cryptokeyId=  Google kms crytoKeyId
  -k, --keyringId=    Google kms keyRingId
  -l, --locationId=   Google kms locationId
  -p, --projectId=    Google projectId

Help Options:
  -h, --help          Show this help message

[decrypt command options]
      -f, --filepath= Path of file to get encrypted string from (default: ./cipher.txt)
      -n, --nonce=    Nonce for decryption
```

## how-to encrypt

pre-requisites:
* `google credentials` on the host the aes-256-gcm-kms binary is run on, see [here](https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials), that has kms encrypt permission
* kms `cryptoKeyId`, `keyRingId`, `locationId` and `projectId` values
* the `plaintext` you want to encrypt in a file on the host you're invoking the binary from (defaults to ./plain.txt)
* a `nonce` string to be used in the [AEAD](https://golang.org/pkg/crypto/cipher/#AEAD) cipher `seal` call

`aes-256-gcm-kms encrypt` will:
* create a new DEK
* encrypt data with DEK
* encrypt DEK using KMS
* outputs concatenated string of encrypted data + encrypted DEK. to command-line and ./cipher.txt

```
$ aes-256-gcm-kms encrypt -h

Usage:
  aes-256-gcm-kms [OPTIONS] encrypt [encrypt-OPTIONS]

Creates a new DEK, encrypts data with DEK, encrypts the DEK using KMS, spits out encrypted data + encrypted DEK.

Application Options:
  -c, --cryptokeyId=  Google kms crytoKeyId
  -k, --keyringId=    Google kms keyRingId
  -l, --locationId=   Google kms locationId
  -p, --projectId=    Google projectId

Help Options:
  -h, --help          Show this help message

[encrypt command options]
      -f, --filepath= Path of file to encrypt (default: ./plain.txt)
      -n, --nonce=    Nonce for encryption
```

