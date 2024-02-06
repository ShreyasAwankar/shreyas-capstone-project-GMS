package functions

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"github.com/Capstone/models"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	// Initialize a Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Register your actual handler function with a route
	router.POST("/", CreateBulk)

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("CreateBulk", router.Handler().ServeHTTP)
}

func CreateBulk(c *gin.Context) {

	// Read the request body as bytes
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read request body"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error reading request body: %s", err.Error()),
			Severity: logging.Error,
		})
		return
	}

	// Convert bytes to string
	file_url := string(bodyBytes)

	// Remove surrounding quotes from the URL
	file_url = strings.Trim(file_url, `"`)

	// Get the data from the URL
	response, err := http.Get(file_url)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file_url"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error reading file_url: %s", err.Error()),
			Severity: logging.Error,
		})
		return
	}
	defer response.Body.Close()

	// Read the content into memory
	data, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file_url"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error While storing data: %s", err.Error()),
			Severity: logging.Error,
		})
		return
	}

	// Create a CSV reader
	reader := csv.NewReader(strings.NewReader(string(data)))

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(500, gin.H{"Error": "Internal Server Error"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error While reading the records in csv file: %s", err.Error()),
			Severity: logging.Error,
		})
		return
	}

	// Process and store records in Firestore
	if err := StoreRecords(ctx, firestoreClient, records); err != nil {
		fmt.Printf("Error bulk creating records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal Server Error"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error storing data in Firestore : %v", err),
			Severity: logging.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Grocery added successfully"})
	logger.Log(logging.Entry{
		Payload:  "Grocery added successfully",
		Severity: logging.Info,
	})

}

// Function to store all records.
func StoreRecords(ctx context.Context, client *firestore.Client, records [][]string) error {

	var Caser = cases.Title(language.English)

	for _, record := range records {
		product := models.Product{
			ProductID:           "",
			ProductName:         record[0],
			Category:            record[1],
			Price:               0,
			Weight:              record[3],
			Brand:               record[4],
			Vegetarian:          true,
			PackageInformation:  record[6],
			ItemPackageQuantity: 0,
			CountryofOrigin:     record[8],
			Manufacturer:        record[9],
		}

		product.ProductName = Caser.String(product.ProductName)
		product.Category = Caser.String(product.Category)

		price, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Error adding product to Firestore: %v", err),
				Severity: logging.Error,
			})
			return err
		}
		product.Price = price

		vegetarian, _ := strconv.ParseBool(record[5])
		if err != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Error adding product to Firestore: %v", err),
				Severity: logging.Error,
			})
			return err
		}
		product.Vegetarian = vegetarian

		itemPackageQuantity, _ := strconv.Atoi(record[7])
		if err != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Error adding product to Firestore: %v", err),
				Severity: logging.Error,
			})
			return err
		}
		product.ItemPackageQuantity = itemPackageQuantity

		randomString := uuid.New().String()[:8]

		product.ProductID = randomString

		// Store in Firestore
		_, err = client.Collection("grocery").Doc(randomString).Set(ctx, product)

		if err != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Error adding product to Firestore: %v", err),
				Severity: logging.Error,
			})
			return err
		}

	}
	return nil
}
