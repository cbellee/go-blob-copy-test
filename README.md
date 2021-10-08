# go-blob-copy-test

# pre-requisites
 - create two storage accounts in different regions
 - create a container in each storage account
 - upload a blob to the 'source' storage account container
 - get storage account keys from portal, 'az' cli, Azure PowwrShell, etc.

# usage
 - git clone this repo
 - `$ cd ./go-blob-copy-test`
 - `$ go build`
 - modify the ./test.sh script, adding your source & destination storage account names, keys, containers and blob name
 - `$ ./test.sh`
 - verify the blob has been copied to the destination storage account container
