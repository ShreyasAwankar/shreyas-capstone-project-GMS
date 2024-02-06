package functions

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

func UploadFile(file multipart.File, object string, storageBucket string) (string, error) {
	// Create a new storage.Writer for the specified object in the bucket
	wc := storageClient.Bucket(storageBucket).Object(object).NewWriter(ctx)

	defer file.Close()

	// Copy the content of the file to the storage.Writer
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}

	// Close the storage.Writer to complete the upload
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	u, _ := url.Parse("https://storage.googleapis.com/shreyas-grocery-image-bucket/" + wc.Attrs().Name)

	return u.String(), nil
}
