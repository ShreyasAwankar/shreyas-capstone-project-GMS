package functions

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"

	"cloud.google.com/go/logging"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

// ByteReaderCloser embeds *bytes.Reader and implements the Close method
type ByteReaderCloser struct {
	*bytes.Reader
}

// Close method for ByteReaderCloser
func (b *ByteReaderCloser) Close() error {
	// Close method for *bytes.Reader is a no-op
	return nil
}

func GenerateThumbnail(c *gin.Context, randomString string) string {
	f, err := c.FormFile("Image")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while fetching the image: %s", err.Error()),
			Severity: logging.Error,
		})
		return ""
	}

	// Open the file
	blobFile, err := f.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while opening the image: %s", err.Error()),
			Severity: logging.Error,
		})
		return ""
	}

	defer blobFile.Close()

	// Decode the image
	img, _, err := image.Decode(blobFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while decoding the image: %s", err.Error()),
			Severity: logging.Error,
		})
		return ""
	}

	// Resize the image using the imaging package
	thumbnail := imaging.Resize(img, 400, 200, imaging.Lanczos)

	// Create a new file buffer for the thumbnail
	thumbnailBuffer := new(bytes.Buffer)
	if err := jpeg.Encode(thumbnailBuffer, thumbnail, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while encoding the thumbnail: %s", err.Error()),
			Severity: logging.Error,
		})
		return ""
	}

	// Use ByteReaderCloser instead of bytes.NewReader directly
	thumbnailReader := &ByteReaderCloser{Reader: bytes.NewReader(thumbnailBuffer.Bytes())}

	// Upload the file to Google Cloud Storage
	img_url, err := UploadFile(thumbnailReader, "thumbnail"+randomString+f.Filename, "shreyas-grocery-image-bucket")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while uploading the image to storage bucket: %s", err.Error()),
			Severity: logging.Error,
		})
		return ""
	}

	return img_url
}
