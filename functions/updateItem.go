package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/Capstone/models"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	// Initialize a Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"*"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for
	corsConfig.AddAllowMethods("OPTIONS")

	// Add the bellow header as we are sending this in request header so as to avoid cors (non-simple request)
	corsConfig.AllowHeaders = []string{"Authorization"}

	// Register the middleware
	router.Use(cors.New(corsConfig))

	// Register your actual handler function with a route
	router.PUT("/", UpdateItem)

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("UpdateItem", router.Handler().ServeHTTP)
}

// @Summary Update a grocery item
// @Description Update a grocery item by providing the ProductId and new product information
// @Tags Grocery
// @Security bearerToken
// @Accept json
// @Produce json
// @Param ProductId query string true "ProductId of the grocery item to be updated"
// @Param Image formData file false "Optional: Updated image file (.jpg)"
// @Param ProductBody formData models.ProductRequestBody true "Updated product information"
// @Success 200 {object} models.UpdateSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /updateItem [put]
func UpdateItem(c *gin.Context) {

	// Check JWT token in Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"Error": "Unauthorized. Missing token."})
		c.Abort()
		return
	}

	// Verify the JWT token
	claims, err := VerifyToken(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"Error": "Unauthorized. Invalid token."})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Unauthorized. Invalid token.: %v", err),
			Severity: logging.Error,
		})
		c.Abort()
		return
	}

	documentId := c.Query("ProductId")

	doc, _ := firestoreClient.Collection("grocery").Doc(documentId).Get(ctx)

	if documentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Please provide ProductId"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Document Id not provided by user: %v", err),
			Severity: logging.Error,
		})
		return
	}

	if doc.Exists() {

		var updateProductData models.Product

		var (
			g          = galidator.New()
			customizer = g.Validator(models.Product{})
		)

		// Bind the form data to the Product struct and perform validation
		if err := c.ShouldBind(&updateProductData); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": customizer.DecryptErrors(err)})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Failed to unmarshal JSON: %s", err.Error()),
				Severity: logging.Error,
			})
			return

		}

		var Caser = cases.Title(language.English)

		updateProductData.ProductName = Caser.String(updateProductData.ProductName)
		updateProductData.Category = Caser.String(updateProductData.Category)

		// Retrieve the uploaded file from the request
		f, err := c.FormFile("Image")
		if err != nil {
			if err != nil {
				if err != http.ErrMissingFile {
					c.JSON(http.StatusInternalServerError, gin.H{
						"Error": "Internal Server Error",
					})
					logger.Log(logging.Entry{
						Payload:  fmt.Sprintf("Error while fetching the image: %s", err.Error()),
						Severity: logging.Error,
					})
					return
				}
			}

		} else {
			// Check if the uploaded file is a .jpg file
			if !strings.HasSuffix(f.Filename, ".jpg") {
				c.JSON(http.StatusBadRequest, gin.H{
					"Error": "Invalid file format. Image file should be '.jpg' file"})
				logger.Log(logging.Entry{
					Payload:  fmt.Sprintf("Invalid file format : %v", err),
					Severity: logging.Error,
				})
				return
			}

			// Open the file
			blobFile, err := f.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Error": "Internal Server Error",
				})
				logger.Log(logging.Entry{
					Payload:  fmt.Sprintf("Error while opening the image: %s", err.Error()),
					Severity: logging.Error,
				})
				return
			}

			defer blobFile.Close()

			// Asynchronously generate and upload thumbnail
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				img_url := GenerateThumbnail(c, updateProductData.ProductID)

				updateProductData.ThumbnailURL = img_url
			}()
			// wg.Wait() to wait for the thumbnail generation to complete.
			wg.Wait()

			// Upload the file to Google Cloud Storage
			img_url, err := UploadFile(blobFile, updateProductData.ProductID+f.Filename, "shreyas-grocery-image-bucket")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Error": "Internal Server Error",
				})
				logger.Log(logging.Entry{
					Payload:  fmt.Sprintf("Error while uploading the image to storage bucket: %s", err.Error()),
					Severity: logging.Error,
				})
				return
			}

			updateProductData.ImageURL = img_url

		}

		updateProductData.ProductID = doc.Ref.ID

		_, err = doc.Ref.Set(ctx, updateProductData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Internal Server Error",
			})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Error while updating the item: %s", err.Error()),
				Severity: logging.Error,
			})
			return
		}

		auditRecord := models.AuditRecord{
			Action:     "Update",
			DocumentID: doc.Ref.ID,
			Timestamp:  time.Now().Format("2006-01-02 15:04:05 PM"),
			UserID:     claims.Username,
		}

		// Serialize the AuditRecord instance to JSON
		auditRecordJSON, err := json.Marshal(auditRecord)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Internal Server Error",
			})

			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Could not serialize the auditRecord data - %v", err),
				Severity: logging.Info,
			})
			return
		}

		auditTopic := pubSubClient.Topic("audit-records")

		// Publish the audit record to Pub/Sub topic
		result := auditTopic.Publish(ctx, &pubsub.Message{
			Data: auditRecordJSON,
		})

		// Check if the publish operation succeeded
		_, err = result.Get(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": fmt.Sprintf("Internal Server Error. %v", err),
			})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Failed to publish pub/sub message: %v", err),
				Severity: logging.Error,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Success":      "Product Updated successfully.",
			"Updated Data": updateProductData,
		})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Document Updated successfully. DocumentId - %v", doc.Ref.ID),
			Severity: logging.Info,
		})

	} else {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Product dose not exist"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Document dose not exist: %v", err),
			Severity: logging.Error,
		})
	}
}
