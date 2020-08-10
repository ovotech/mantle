#!/usr/bin/env bash
set -e
call_mantle_cmds () {
    ../mantle decrypt -n "$KEY_NAME" -m "$PROVIDER" -r
    plaintext=$(cat plain.txt)
    if [[ $plaintext != "helloworld" ]]
    then
        echo "Unexpected plaintext: $plaintext"
        exit 1
    fi
    echo "Successfully decrypted"
    echo "-----------------------------------------------------------"
    ../mantle encrypt -n "$KEY_NAME" -m "$PROVIDER"
    echo "Successfully encrypted"
    echo "-----------------------------------------------------------"
    ../mantle reencrypt -n "$KEY_NAME" -m "$PROVIDER"
    echo "Successfully re-encrypted"
    echo "-----------------------------------------------------------"
}

if [ -n "$MANTLE_GCP_KMS_KEY_NAME" ]; then
    cp gcp-cipher.txt cipher.txt
    KEY_NAME="$MANTLE_GCP_KMS_KEY_NAME"
    PROVIDER="gcp"
    echo "-----------------------------------------------------------"
    echo "GCP TESTS"
    echo "-----------------------------------------------------------"
    call_mantle_cmds
fi

if [ -n "$MANTLE_AWS_KMS_KEY_NAME" ]; then
    cp aws-cipher.txt cipher.txt
    KEY_NAME="$MANTLE_AWS_KMS_KEY_NAME"
    PROVIDER="aws"
    echo "-----------------------------------------------------------"
    echo "AWS TESTS"
    echo "-----------------------------------------------------------"
    call_mantle_cmds
fi