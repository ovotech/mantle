# aes-256-gcm-kms

## intro

A Go program to simplify the encryption & decryption of strings, using 256 bit AES keys in Galois/Counter Mode (GCM), with cloud-based KMS services and multiple key layers (specifically [Envelope Encryption](https://cloud.google.com/kms/docs/envelope-encryption))

In short, plaintext is encrypted using a 256 bit AES "Data Encryption Key/DEK" in GCM mode. This DEK is then encrypted using KMS. The concatenated encrypted data and encrypted DEK are given back to the user via CLI and a file.

Ciphertext is decrypted by following the same process in reverse. The DEK is decrypted using KMS, then used to decrypt the ciphertext. The resulting plaintext is given back to the user

## how-to decrypt

### pre-requisites:
* `google credentials` on the host the aes-256-gcm-kms binary is run on, see [here](https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials), that has kms decrypt permission
* kms `cryptoKeyId`, `keyRingId`, `locationId` and `projectId` values
* the `ciphertext` you want to decrypt in a file on the host you're invoking the binary from (defaults to ./cipher.txt)

### `aes-256-gcm-kms decrypt` will:
* split ciphertext into encrypted data + encrypted DEK
* decrypts the encrypted DEK using KMS
* decrypts the encrypted data
* outputs to ./plain.txt
* zerofills and deletes ciphertext file

### example

```
$ aes-256-gcm-kms decrypt -c <cryptoKeyId> -k <keyringId> -l europe-west2 -p <projectId>
Decrypting...
Decryption successful, plaintext available at ./plain.txt
Wiped 643 bytes from ./cipher.txt.
```

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
```

## how-to encrypt

### pre-requisites:
* `google credentials` on the host the aes-256-gcm-kms binary is run on, see [here](https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials), that has kms encrypt permission
* kms `cryptoKeyId`, `keyRingId`, `locationId` and `projectId` values
* the `plaintext` you want to encrypt in a file on the host you're invoking the binary from (defaults to ./plain.txt)

### `aes-256-gcm-kms encrypt` will:
* create a new DEK
* encrypt data with DEK
* encrypt DEK using KMS
* outputs concatenated string of encrypted data + encrypted DEK. to command-line and ./cipher.txt
* zerofills and deletes plaintext file

### example

```
$ aes-256-gcm-kms encrypt -c <cryptoKeyId> -k <keyringId> -l europe-west2 -p <projectId>
Encrypting...
-----BEGIN (ENCRYPTED DATA + DEK) STRING-----
***REMOVED***
wV/Wow6eWHFL80x3xl7vIgDVN5CdRAOVpZL2kJV3
coDbctszL5LJHaLL22YVYaJwojETz5Aff4Kss98p
MIRahCJ1D8EFNoBbTAQTUGNJAJGc11YcX3sWpsYB
h3BookBa6KEvnmNFfw8F6M71zpdmByS1p/k8/1Z/
TAX/Dj0wxcm2g/ez7gA0e/vFQXQjJYqSkb0xJuQX
SVaDoXap3HF7NbikcklBPBkDvy408Hogapvh4OF2
vL9tlhGoERUkrWcwXQfcZjk1B3Sjh45UDTHySTs+
m4Eco7MOur6LvfrGKJuX6qJhvppxUDv2ZTCeMCrK
d9AjmCqleD/iSthZN1FKjQ3zLowlnsvWIMnaeEC+
h5W8NIjKm4YQCY2yGj3V6AhdBMvujXLX1aYbIHSf
GfIzLhHSKI7vUm0RFN5irblcoC+sBkRf8NAKJAB1
PbjZJT8wZ94zMUnqrUNNCJqzoky5PFiAY0x077co
SHATyRJJAOR2fnkCjptlffrP0/y8Jhs7ogtttzwt
mkJtdbf9ltQw2ak1OJI3h7NC9vLqfDzGQFeO396C
RRt3E3ly9MifB+cFe4Fnowcq0g==
-----END (ENCRYPTED DATA + DEK) STRING-----
Encryption successful, ciphertext available at ./cipher.txt
Wiped 340 bytes from ./plain.txt.
```

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
```
