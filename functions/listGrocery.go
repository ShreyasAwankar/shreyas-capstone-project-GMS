package functions

import (
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/logging"
	"github.com/Capstone/models"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/api/iterator"
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
	router.GET("/", ListGrocery)

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("ListGrocery", router.Handler().ServeHTTP)
}

// @Summary List grocery items with optional filters and pagination
// @Description Retrieve a list of grocery items with optional filters like productName, priceGreaterThan, priceLessThan, and category. Supports pagination.
// @Tags Grocery
// @Security bearerToken
// @Accept json
// @Produce json
// @Param productName query string false "Filter by product name"
// @Param priceGreaterThanOrEqualTo query string false "Filter by price greater than this value"
// @Param priceLessThanOrEqualTo query string false "Filter by price less than this value"
// @Param priceEqualTo query string false "Filter by price less than this value"
// @Param category query string false "Filter by category"
// @Param page query integer false "Page number for pagination"
// @Success 200 {object} models.ListGrocerySuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /listGrocery [get]
func ListGrocery(c *gin.Context) {

	// Get query parameters from the client
	productName := c.Query("productName")
	greaterThan := c.Query("priceGreaterThanOrEqualTo")
	lessThan := c.Query("priceLessThanOrEqualTo")
	equalTo := c.Query("priceEqualTo")
	category := c.Query("category")
	pageStr := c.Query("page")

	var page int

	if pageStr == "" {
		page = 1
	} else {
		// Get page number from query parameters
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Invalid page number provided by user: %v", err),
				Severity: logging.Error,
			})
			return
		}
	}

	// Calculate the start index for pagination
	pageSize := 10
	startIndex := (page - 1) * pageSize

	var products []models.Product

	// Construct the base query without pagination
	baseQuery := firestoreClient.Collection("grocery")

	filteredQuery := baseQuery.Query

	var Caser = cases.Title(language.English)

	// Apply filters based on query parameters
	if productName != "" {
		filterParam := Caser.String(productName)
		filteredQuery = filteredQuery.Where("ProductName", "==", filterParam)
	}

	if greaterThan != "" {
		greaterThanInt, err := strconv.ParseInt(greaterThan, 10, 64)
		if err == nil {
			filteredQuery = filteredQuery.Where("Price", ">=", greaterThanInt)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Error": " Price filter value should be a number only"})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Invalid query parameter provided by user: %v", err),
				Severity: logging.Error,
			})
			return
		}
	}

	if lessThan != "" {
		lessThanInt, err := strconv.ParseInt(lessThan, 10, 64)
		if err == nil {
			filteredQuery = filteredQuery.Where("Price", "<=", lessThanInt)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Price filter value should be a number only"})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Invalid query parameter provided by user: %v", err),
				Severity: logging.Error,
			})
			return
		}
	}

	if equalTo != "" {
		equalToInt, err := strconv.ParseInt(equalTo, 10, 64)
		if err == nil {
			filteredQuery = filteredQuery.Where("Price", "==", equalToInt)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Price filter value should be a number only"})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Invalid query parameter provided by user: %v", err),
				Severity: logging.Error,
			})
			return
		}
	}

	if category != "" {
		filterParam := Caser.String(category)
		filteredQuery = filteredQuery.Where("Category", "==", filterParam)
	}

	// Execute the filtered query with pagination
	iter := filteredQuery.Offset(startIndex).Limit(pageSize).Documents(ctx)
	defer iter.Stop()

	for {

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Failed to iterate through grocery items: %v", err),
				Severity: logging.Error,
			})

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to iterate through grocery items",
			})
			return
		}

		var product models.Product

		err = doc.DataTo(&product)
		if err != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Failed to parse employee data: %v", err),
				Severity: logging.Error,
			})

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to parse the product data",
			})
			return
		}

		products = append(products, product)
	}

	var nextPage interface{}
	if len(products) == pageSize {
		nextPage = page + 1
	} else {
		nextPage = "No more pages"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "Grocery items fetched",
		"page":               page,
		"number of products": len(products),
		"products":           products,
		"next page":          nextPage,
	})

	logger.Log(logging.Entry{
		Payload:  "Grocery fetched successfully",
		Severity: logging.Info,
	})

}
