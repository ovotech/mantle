# AES-256-GCM-KMS

This is a Go program to simplify the encryption & decryption of byte arrays,
using 256 bit [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard)
keys in [Galois/Counter Mode](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
(GCM), with cloud-based KMS services (currently only Google KMS) and multiple
key layers (specifically
[Envelope Encryption](https://cloud.google.com/kms/docs/envelope-encryption)).


## Install

### From Binary Releases

Darwin, Linux and Windows Binaries can be downloaded from [the Releases page](https://github.com/ovotech/aes-256-gcm-kms/releases).

Try it out:

```
$ aes-256-gcm-kms -h
```

### From Source

```
$ git clone git@github.com:ovotech/aes-256-gcm-kms.git

$ cd aes-256-gcm-kms

$ go build

$ ./aes-256-gcm-kms -h
```

## Getting Started

### GCP

You'll need to [create a Key Ring and Key](https://cloud.google.com/kms/docs/creating-keys#kms-create-keyring-cli)
(if not already present) in Google KMS to allow the tool to encrypt/decrypt.

### Authorisation

If you have `gcloud` set up locally, and your user has the `Cloud KMS CryptoKey
Encrypter` and/or `Cloud KMS CryptoKey Decrypter` Role(s), the tool will
already be able to encrypt/decrypt.

If you're running the binary in an automated way (i.e. with a Service Account):

* The Service Account will need the Cloud KMS CryptoKey Encrypter and/or
Decrypter Role(s
* Drop the Service Account's key.json onto the host you're running the tool on
* Set the `GOOGLE_APPLICATION_CREDENTIALS` env var as the path of the key.json

### Obtain The Key's Resource Id

The final piece of the puzzle is obtaining the name of the KMS Key the tool is
going to use.

Using `gcloud`, you can do this by issuing:

```
# get the name of the Keyring that 'holds' the required Key
$ gcloud kms keyrings list --location <location>

# get the name of the Key
$ gcloud kms keys list --location <location> --keyring <keyring_name>
```

The `NAME` value returned by the last command is Google's `ResourceId` for the
Key. It's this value that you can give to the `aes-256-gcm-kms` binary in the
`-n,--keyName` flag, to get it to work.

To test this out, you should be able to:

```
# create plain.text
$ echo "helloworld" > plain.txt

# issue the encrypt command. The binary should output to command line the
# encrypted string, remove the plain.txt file, and create a cipher.txt file
$ aes-256-gcm-kms encrypt -n <key_name>

# now we can decrypt back again, you should be left with a new plain.txt
$ aes-256-gcm-kms decrypt -n <key_name>

$ cat plain.txt
```

## How It Works

1. A new 256-bit AES key and a 96-bit nonce are generated every time you issue
 the `encrypt` command. The key and nonce are used to encrypt your plaintext.

2. The AES key, also known as the Data Encryption Key (DEK), is then encrypted
using the Cloud KMS Service.

3. The concatenated encrypted data, nonce and encrypted DEK are then given back
to the user as the ciphertext.

Decrypting is the same process but in reverse.


## Ciphertext Structure

The structure of a ciphertext will be:

```
encryptedData[n]nonce[12]encryptedDEK[113]
```

The encrypted DEK is 113 chars. This is the length of the string returned by
Google KMS after it's been base64 decoded. The length of the encrypted data will
depend on the length of your plaintext.

## Notes

### Newlines

When encrypting, newline chars are by default inserted into the ciphertext,
every 40 chars. This is to play nicer with any max line lengths when storing in
source code. This functionality can be disabled using the `-s, --singleLine`
flag.

So long as the decrypting process is only removing newline chars at the end of
lines, it shouldn't need to differentiate the two 'modes'

### Nonce

The 96-bit nonce serves the purpose of the IV.


## Example

```
$ aes-256-gcm-kms encrypt -n <key_name>

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
m4Eco7MOur6LvfrGKJuX6qJubppxUDv2ZTCeMCrK
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

Your ciphertext is the string between (not including) the `BEGIN` and `END`
markers. The resulting `cipher.txt` file will only contain the ciphertext
string.
