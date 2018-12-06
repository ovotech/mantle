# Example of Mantle with Kubernetes init-container

Follow the walkthrough below to get a feel for how Mantle can be used in a
Kubernetes init-container.

Don't forget to [install Mantle](https://github.com/ovotech/mantle#Install) and grab the Resource ID of your Google KMS key, see [here](https://github.com/ovotech/mantle#obtain-the-keys-resource-id) for help.

```
# copy our 'config' file to plain.txt so we don't lose the original upon
# encryption (Mantle will delete the source file)
$ cp banksy.txt plain.txt

# set your KMS Key Resource ID in an env var
$ export KMS_KEY=<kms_key_resource_id>

# create a k8s configmap containing the KMS Key Resource ID
$ kubectl create configmap mantle-kms-key --from-literal=resource.id=$KMS_KEY

# encrypt using Mantle (encrypted string is outputted to command-line, but will
# also be present in ./cipher.txt)
$ mantle encrypt -n $KMS_KEY

# create a k8s configmap from the encrypted cipher.txt
$ kubectl create configmap mantle-config --from-file=./cipher.txt

# create the k8s deployment
$ kubectl apply -f deployment.yaml

# get the k8s pod name, and get a shell into it
$ kubectl exec -it $(kubectl get pods --selector=app=mantle \
   --output=jsonpath={.items..metadata.name}) sh

# cat the file
/ # cat /etc/decrypted/banksy.txt

# delete the resources when you're done
$ kubectl delete configmap mantle-config \
    && kubectl delete configmap mantle-kms-key
    && kubectl delete deployment mantle
```

To summarise what's going on here:

1. Create a **k8s configmap** containing the KMS Key Resource ID.
1. **Encrypt config** to produce a cipher.txt file.
2. Create a **k8s configmap from the cipher.txt file**.
3. Create a **k8s deployment**, starting a pod that runs an init-container that
decrypts the config, and passes it to the 'app' container that runs afterwards.
4. Get a shell into the app container, to **check the decrypted file**.
