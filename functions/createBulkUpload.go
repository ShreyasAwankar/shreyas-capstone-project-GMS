package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	router.POST("/", CreateBulkUpload)

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("CreateBulkUpload", router.Handler().ServeHTTP)
}

// @Summary Create grocery items in bulk from a CSV file
// @Description Upload a CSV file containing grocery items and create them in bulk
// @Tags Grocery
// @Accept json
// @Produce json
// @Security bearerToken
// @Param csvfile formData file true "CSV file containing grocery items"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /createBulkUpload [post]
func CreateBulkUpload(c *gin.Context) {

	// Check JWT token in Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"Error": "Unauthorized. Missing token."})
		c.Abort()
		return
	}

	// Verify the JWT token
	_, err := VerifyToken(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"Error": "Unauthorized. Invalid token."})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Unauthorized. Invalid token.: %v", err),
			Severity: logging.Error,
		})
		c.Abort()
		return
	}

	file, err := c.FormFile("csvfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Could not get the file: %v", err),
			Severity: logging.Error,
		})
		return
	}

	// Check if the uploaded file is a CSV file
	if !strings.HasSuffix(file.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format. Please upload a CSV file"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Invalid file format : %v", err),
			Severity: logging.Error,
		})
		return
	}

	// Open and read the CSV file
	csvFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error opening file : %v", err),
			Severity: logging.Error,
		})
		return
	}
	// defer csvFile.Close()

	randomString := uuid.New().String()[:5]

	// Upload the file to Google Cloud Storage
	file_url, err := UploadFile(csvFile, file.Filename+randomString, "shreyas-csv-bucket")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal Server Error"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while uploading the image to storage bucket: %s", err.Error()),
			Severity: logging.Error,
		})
		return
	}

	pubSubMessage, err := json.Marshal(file_url)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Internal Server Error",
		})

		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Could not serialize the file to send to pubsub topic data - %v", err),
			Severity: logging.Info,
		})
		return
	}

	topic := pubSubClient.Topic("create-bulk-records")

	// Publish the audit record to Pub/Sub topic
	result := topic.Publish(ctx, &pubsub.Message{
		Data: pubSubMessage,
	})

	// Check if the publish operation succeeded
	_, err = result.Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Internal Server Error",
		})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Failed to publish pub/sub message to topic create-bulk-records: %v", err),
			Severity: logging.Error,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "File uploaded successfully. Grocery will be added automatically.",
	})
	logger.Log(logging.Entry{
		Payload:  "File uploaded successfully. Pub/sub message published successfully",
		Severity: logging.Info,
	})

}
