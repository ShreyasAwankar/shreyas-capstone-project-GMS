package functions

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/logging"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
)

func init() {
	// Initialize a Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Register your actual handler function with a route
	router.POST("/", StoreAuditRecords)

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("StoreAuditRecords", router.Handler().ServeHTTP)
}

func StoreAuditRecords(c *gin.Context) {

	// var auditRecord models.AuditRecord
	var auditRecord map[string]interface{}

	if err := c.BindJSON(&auditRecord); err != nil {
		// Handle JSON unmarshal error
		c.JSON(400, gin.H{"error": err.Error()})

		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while saving item to firestore: %s", err.Error()),
			Severity: logging.Error,
		})
		return
	}
	// Store other data to firestore.
	_, _, err := firestoreClient.Collection("audit-records").Add(ctx, auditRecord)

	if err != nil {
		c.JSON(500, gin.H{"500": "Internal Servrer Error"})

		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while saving item to firestore: %s", err.Error()),
			Severity: logging.Error,
		})
		// Manually NACK the Pub/Sub message
		c.Status(http.StatusNoContent)
		return
	}

	// Successful processing, manually ACK the Pub/Sub message with a different status code.
	c.Status(http.StatusAccepted)
}
