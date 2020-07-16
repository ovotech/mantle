# Mantle
This is a Go program to simplify the encryption & decryption of byte arrays,
using 256 bit [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard)
keys in [Galois/Counter Mode](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
(GCM), with cloud-based KMS services (currently only AWS and GCP) and multiple
key layers (specifically
[Envelope Encryption](https://cloud.google.com/kms/docs/envelope-encryption)).

This avoids the need to send secret data to Cloud KMS services when encrypting
(only your own data encryption key is sent), and allows for encryption of data
bigger than 4096 bytes (a restriction often imposed by Cloud KMS).

## Install

### From Binary Releases

Darwin, Linux and Windows Binaries can be downloaded from [the Releases page](https://github.com/ovotech/mantle/releases).

Try it out:

```
$ mantle -h
```

### From Source

```
$ git clone git@github.com:ovotech/mantle.git

$ cd mantle

$ go build

$ ./mantle -h
```

### As a Go Dependency

```Go
$ go get -u github.com/ovotech/mantle
```

## Getting Started

### AWS

Create a symmetric Customer Managed Key (CMK). See AWS doco [here](https://docs.aws.amazon.com/kms/latest/developerguide/create-keys.html#create-symmetric-cmk).

It's recommended to give the key an alias, allowing you to change the underlying
CMK if required. Your user will [permission](https://docs.aws.amazon.com/kms/latest/developerguide/key-policies.html) to use the CMK.

Authenticate your AWS CLI, see [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) for help.

Now give it a go:

```bash
$ echo "helloworld" > plain.txt
$ ./mantle encrypt -n alias/my-kms-key -m aws
$ cat cipher.txt

$ ./mantle decrypt -m aws
$ cat plain.txt
```
#### Notes

* `mantle` doesn't need the keyId (passed in using the `-n,--keyName` flag) when
decrypting as it's stored in the ciphertext.

* You can use the Key ID, Key ARN, Alias name or Alias ARN in the `-n,--keyName`
flag when encrypyting, as described [here](https://docs.aws.amazon.com/cli/latest/reference/kms/encrypt.html#options).

* `plain.txt` and `cipher.txt` are used as default filepaths for encrypt
and decrypt respectively. You can control this in each operation using the 
`-f,--filepath` flag.

### GCP

You'll need to [create a Key Ring and Key](https://cloud.google.com/kms/docs/creating-keys#kms-create-keyring-cli) in Google KMS to allow the tool to
encrypt/decrypt.

#### Authorisation

If you have `gcloud` set up locally, and your user has the `Cloud KMS CryptoKey
Encrypter` and/or `Cloud KMS CryptoKey Decrypter` Role(s), the tool will
already be able to encrypt/decrypt.

If you're running the binary in an automated way (i.e. with a Service Account):

* The Service Account will need the Cloud KMS CryptoKey Encrypter and/or
Decrypter Role(s)
* Drop the Service Account's key.json onto the host you're running the tool on
* Set the `GOOGLE_APPLICATION_CREDENTIALS` env var as the path of the key.json

#### Obtain The Key's Resource ID

Using `gcloud`, you can do this by issuing:

```
# get the name of the Keyring that 'holds' the required Key
$ gcloud kms keyrings list --location <location>

# get the name of the Key
$ gcloud kms keys list --location <location> --keyring <keyring_name>
```

The `NAME` value returned by the last command is Google's `Resource ID` for the
Key. 

It's this value that you can give to the `mantle` binary in the
`-n,--keyName` flag, to get it to work.

Alternatively to using `gcloud` you can get the Resource ID from the [Google Cloud Console](https://console.cloud.google.com/security/kms); click on the KeyRing
you want to use, and select "Copy Resource ID" from the menu to the right of the
correct Key.

To test this out, you should be able to:

```bash
$ export KEY_NAME="projects/<project_name>/locations/<location>/keyRings/<keyring_name>/cryptoKeys/<key_name>"

$ echo "helloworld" > plain.txt
$ mantle encrypt -n $KEY_NAME -m gcp
$ cat cipher.txt

$ mantle decrypt -n $KEY_NAME -m gcp
$ cat plain.txt
```

#### Notes

* The `-m,--kmsProvider` flag defaults to `gcp`. This is to make `mantle`
backwards compatible for GCP users who don't use that flag (when `mantle` was
first released, GCP was the only provider integrated).

## How It Works

1. A new 256-bit AES key and a 96-bit nonce are generated every time you issue
 the `encrypt` command. The key and nonce are used to encrypt your plaintext.

2. The AES key, also known as the Data Encryption Key (DEK), is then encrypted
using the Cloud KMS Service.

3. The concatenated encrypted data, nonce and encrypted DEK are then given back
to the user as the ciphertext.

Decrypting is the same process but in reverse.


## Ciphertext Structure

For AWS, the structure of the ciphertext is:

```
encryptedData[n]nonce[12]encryptedDEK[185]
```

For GCP, the structure of the ciphertext is:

```
encryptedData[n]nonce[12]encryptedDEK[114]
```

The length of the `encryptedDEK` is determined by the length of the response
from the KMS providers. For AWS this is 185 chars, for GCP it's 114 chars. 
The length of the encrypted data will depend on the length of your plaintext.

## Notes

### Newlines

When encrypting, newline chars are by default inserted into the ciphertext,
every 40 chars. This is to play nicer with any max line lengths when storing in
source code. This functionality can be disabled using the `-s, --singleLine`
flag.

So long as the decrypting process is only removing newline chars at the end of
lines, it shouldn't need to differentiate the two 'modes'


### IV

Mantle uses the [crypto/rand](https://golang.org/pkg/crypto/rand/) Reader to
generate a new IV every time the `encrypt` command is called.


*"Reader is a global, shared instance of a cryptographically secure random number generator.*

*On Linux, Reader uses getrandom(2) if available, /dev/urandom otherwise. On
OpenBSD, Reader uses getentropy(2). On other Unix-like systems, Reader reads
from /dev/urandom. On Windows systems, Reader uses the CryptGenRandom API. On
Wasm, Reader uses the Web Crypto API."*

The IV is created by reading 96 bits (12 bytes) from this Reader.


### Zero-fill and Delete

By default, when performing either `encrypt` or `decrypt` commands, the tool
will zero-fill and delete the source file, so plain.txt or cipher.txt (or an
overriding filepath you've set) respectively.

When decrypting, you can use the `-r,--retainCipherText` flag in order to
retain the ciphertext file. There's no option to retain the source file when
encrypting.


## Example

```
$ mantle encrypt -n <key_name>

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


## Contributing

Contributions are very welcome, please fork or branch and raise a PR.
