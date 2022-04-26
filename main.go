package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func Upload(accountName, accountKey, containerName, fileName string, data io.ReadSeeker) error {
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	serviceClient, err := azblob.NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		return err
	}

	containerClient := serviceClient.NewContainerClient(containerName)

	blockBlobClient := containerClient.NewBlockBlobClient(fileName)

	_, err = blockBlobClient.Upload(context.TODO(), streaming.NopCloser(data), nil)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}

	containerName, ok := os.LookupEnv("AZURE_STORAGE_CONTAINER_NAME")
	if !ok {
		panic("AZURE_STORAGE_CONTAINER_NAME could not be found")
	}

	err := Upload(accountName, accountKey, containerName, "bar", strings.NewReader("Hello world!"))
	if err != nil {
		log.Fatal(err)
	}
}
