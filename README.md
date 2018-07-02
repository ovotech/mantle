# aes-256-gcm-kms

## intro

A Go program to simplify the encryption & decryption of strings, using 256 bit [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) keys in [Galois/Counter Mode](https://en.wikipedia.org/wiki/Galois/Counter_Mode) (GCM), with cloud-based KMS services (currently only Google KMS) and multiple key layers (specifically [Envelope Encryption](https://cloud.google.com/kms/docs/envelope-encryption))

In short, plaintext is encrypted using a generated 256 bit AES "Data Encryption Key/DEK" in GCM mode. This DEK is then encrypted using KMS. A 96-bit nonce is also generated. The concatenated encrypted data, nonce and encrypted DEK are given back to the user via CLI and a file.

Ciphertext is decrypted by following the same process in reverse. The DEK (obtained from ciphertext) is decrypted using KMS, then used, along with the nonce, to decrypt the 'data' section of the ciphertext. The resulting plaintext is given back to the user in a file.

## how-to encrypt

### pre-requisites:
* `google credentials` on the host the aes-256-gcm-kms binary is run on, see [here](https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials), that has kms encrypt permission
* kms `cryptoKeyId`, `keyRingId`, `locationId` and `projectId` values
* the `plaintext` you want to encrypt in a file on the host you're invoking the binary from (defaults to ./plain.txt)

### `aes-256-gcm-kms encrypt` will:
* create a new DEK and a new nonce
* encrypt data with DEK and nonce
* encrypt DEK using KMS
* concatenates string of encrypted data, nonce and encrypted DEK. to command-line and ./cipher.txt
* base64 encodes
* outputs string to CLI and cipher.txt
* zerofills and deletes plaintext file

### example

```
$ aes-256-gcm-kms encrypt -c <cryptoKeyId> -k <keyringId> -l europe-west2 -p <projectId>
Encrypting...
-----BEGIN (ENCRYPTED DATA + DEK) STRING-----
y+PvJrf0QJKKSp85C0MN6q2v7EhMeorNJG+5FLiN
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
  -n, --keyName=      Google kms keyName
  -l, --locationId=   Google kms locationId
  -p, --projectId=    Google projectId

Help Options:
  -h, --help          Show this help message

[encrypt command options]
      -f, --filepath= Path of file to encrypt (default: ./plain.txt)
      -s, --singleLine  Disable use of newline chars in ciphertext
```

## how-to decrypt

### pre-requisites:
* `google credentials` on the host the aes-256-gcm-kms binary is run on, see [here](https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials), that has kms decrypt permission
* kms `cryptoKeyId`, `keyRingId`, `locationId` and `projectId` values
* the `ciphertext` you want to decrypt in a file on the host you're invoking the binary from (defaults to ./cipher.txt)

### `aes-256-gcm-kms decrypt` will:
* base64 decodes
* split ciphertext into encrypted data, nonce and encrypted DEK
* decrypts the encrypted DEK using KMS
* decrypts the encrypted data using DEK and nonce
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
  -n, --keyName=      Google kms keyName
  -l, --locationId=   Google kms locationId
  -p, --projectId=    Google projectId

Help Options:
  -h, --help          Show this help message

[decrypt command options]
      -f, --filepath= Path of file to get encrypted string from (default: ./cipher.txt)
      -v, --validate  Validate decryption works; don't produce a plain.txt
```

## ciphertext contents

To clarify the contents of a ciphertext, after decoding (string length in brackets):

```
[encryptedData (n)][nonce (12)][encrypted DEK (113)]
```

The encrypted DEK is 113 chars. This is the length of the string returned by Google KMS after it's been base64 decoded.

## notes

### newlines

When encrypting, newline chars are by default inserted into the ciphertext, every 40 chars. This is to play nicer with any max line lengths when storing in source code. This functionality can be disabled using the `-s, --singleLine` flag.

So long as the decrypting process is only removing newline chars at the end of lines, it shouldn't need to differentiate the two 'modes'

### nonce

Currently, there's no purpose to the nonce. An auto-generated nonce is included in the ciphertext, which can then be used by the decrypting process. Originally, the nonce was required as an input parameter by the user in both `encrypt` and `decrypt`, and wasn't stored in the ciphertext at all.

There could be arguments for removing the nonce entirely (it doesn't really suit this kind of static file encryption), or resinstate it as a required user input parameter (at least for `decrypt`; it could still be auto-generated in `encrypt`).
