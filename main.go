package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

var srcAccountName string
var destAccountName string
var srcAccountKey string
var destAccountKey string
var srcContainerName string
var destContainerName string
var blobName string

func init() {
	flag.StringVar(&srcAccountName, "source-account-name", "", "source storage account name")
	flag.StringVar(&destAccountName, "destination-account-name", "", "destination storage account name")
	flag.StringVar(&srcAccountKey, "source-account-key", "", "source storage account key")
	flag.StringVar(&destAccountKey, "destination-account-key", "", "destination storage account key")
	flag.StringVar(&srcContainerName, "source-container-name", "", "source storage account container")
	flag.StringVar(&destContainerName, "destination-container-name", "", "destination storage account container")
	flag.StringVar(&blobName, "blob-name", "", "blob name")
}

func main() {
	flag.Parse()

	srcCredential, err := azblob.NewSharedKeyCredential(srcAccountName, srcAccountKey)
	if err != nil {
		log.Fatal(err)
	}
	srcSasTokenUrl := getBlobSasToken(srcAccountName, srcContainerName, blobName, srcCredential, azblob.BlobSASPermissions{Read: true, Add: true, Create: true, Write: true})

	destCredential, err := azblob.NewSharedKeyCredential(destAccountName, destAccountKey)
	if err != nil {
		log.Fatal(err)
	}
	destSasTokenUrl := getBlobSasToken(destAccountName, destContainerName, blobName, destCredential, azblob.BlobSASPermissions{Read: true, Add: true, Create: true, Write: true})

	ctx := context.Background()

	s, _ := url.Parse(srcSasTokenUrl)
	d, _ := url.Parse(destSasTokenUrl)
	destBlobURL := azblob.NewBlobURL(*d, azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))

	startCopy, err := destBlobURL.StartCopyFromURL(ctx, *s, nil, azblob.ModifiedAccessConditions{}, azblob.BlobAccessConditions{}, azblob.DefaultAccessTier, nil)
	if err != nil {
		log.Fatal(err)
	}

	copyID := startCopy.CopyID()
	copyStatus := startCopy.CopyStatus()
	for copyStatus == azblob.CopyStatusPending {
		time.Sleep(time.Second * 2)
		getMetadata, err := destBlobURL.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
		if err != nil {
			log.Fatal(err)
		}
		copyStatus = getMetadata.CopyStatus()
	}
	fmt.Printf("Copy from %s to %s: ID=%s, Status=%s\n", s.String(), destBlobURL, copyID, copyStatus)
}

func getBlobSasToken(accountName string, containerName string, blobName string, credential azblob.StorageAccountCredential, blobPermissions azblob.BlobSASPermissions) (sasUrl string) {
	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		ContainerName: containerName,
		BlobName:      blobName,
		Permissions:   blobPermissions.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		log.Fatal(err)
	}

	qp := sasQueryParams.Encode()
	sasUrl = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", accountName, containerName, blobName, qp)
	fmt.Printf("SAS Blob Uri: %s", sasUrl)
	return sasUrl
}
