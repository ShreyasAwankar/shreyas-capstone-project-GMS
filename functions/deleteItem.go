package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/Capstone/models"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	router.DELETE("/", DeleteItem)

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("DeleteItem", router.Handler().ServeHTTP)
}

// @Summary Delete a grocery item
// @Description Delete a grocery item by providing the ProductId
// @Tags Grocery
// @Security bearerToken
// @Accept json
// @Produce json
// @Param ProductId query string true "ProductId of the grocery item to be deleted"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /deleteItem [delete]
func DeleteItem(c *gin.Context) {

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

	if documentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Please provide ProductId"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Document Id not provided by user: %v", err),
			Severity: logging.Error,
		})
		return
	}

	doc, _ := firestoreClient.Collection("grocery").Doc(documentId).Get(ctx)

	if doc.Exists() {

		_, err := doc.Ref.Delete(ctx)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Error while deleting the item: %s", err.Error()),
				Severity: logging.Error,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Success": fmt.Sprintf("Document with ID %v deleted successfully.", doc.Ref.ID),
		})

		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Document fetched successfully. DocumentId - %v", doc.Ref.ID),
			Severity: logging.Info,
		})

		auditRecord := models.AuditRecord{
			Action:     "Delete",
			DocumentID: doc.Ref.ID,
			Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
			UserID:     claims.Username,
		}

		// Serialize the AuditRecord instance to JSON
		auditRecordJSON, err := json.Marshal(auditRecord)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Internal Server Error. %v",
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

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Document dose not exist"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Document dose not exist: %v", err),
			Severity: logging.Error,
		})
	}
}
