package functions

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/logging"
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

	// Register the middleware
	router.Use(cors.New(corsConfig))

	// Register your actual handler function with a route
	router.POST("/", GetItemById)

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("GetItemById", router.Handler().ServeHTTP)
}

// @Summary Get a grocery item by DocumentId
// @Description Retrieve a grocery item by providing the DocumentId
// @Tags Grocery
// @Security bearerToken
// @Accept json
// @Produce json
// @Param ProductId query string true "DocumentId of the grocery item to be retrieved"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /getItemById [get]
func GetItemById(c *gin.Context) {

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
		c.JSON(http.StatusOK, gin.H{
			"Success":      "Product fetched successfully.",
			"Product Data": doc.Data(),
		})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Document fetched successfully. DocumentId - %v", doc.Ref.ID),
			Severity: logging.Info,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Product dose not exist"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Document dose not exist: %v", err),
			Severity: logging.Error,
		})
	}
}
