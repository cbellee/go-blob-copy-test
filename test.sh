srcAccountName="<add your source account name>"
srcAccountKey="<add your source account key>"
srcContainerName="<add your soruce container name>"

destAccountName="<add your destination account name>"
destAccountKey="<add your destination account key>"
destContainerName="<add your destination container name>"
blobName="<add your blob name>"

./go-blob-copy-test \
    -source-account-name $srcAccountName \
    -destination-account-name $destAccountName \
    -source-account-key $srcAccountKey \
    -destination-account-key $destAccountKey \
    -source-container-name $srcContainerName \
    -destination-container-name $destContainerName \
    -blob-name $blobName
