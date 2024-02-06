package functions

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"github.com/Capstone/models"
	"github.com/gin-gonic/gin"
)

func bulkCreate(ctx context.Context, client *firestore.Client, records [][]string) error {

	for _, record := range records {
		product := models.Product{
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

		// Store in Firestore
		_, _, err = client.Collection("grocery").Add(ctx, product)

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

func HandleCSVUpload(c *gin.Context) {

	// Set CORS headers for the preflight request
	if c.Request.Method == http.MethodOptions {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Max-Age", "3600")
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	c.Header("Access-Control-Allow-Origin", "*")

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error."})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error opening file : %v", err),
			Severity: logging.Error,
		})
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading CSV file"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error reading CSV file : %v", err),
			Severity: logging.Error,
		})
		return
	}

	// Process and store records in Firestore
	if err := bulkCreate(ctx, firestoreClient, records); err != nil {
		fmt.Printf("Error bulk creating records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing data in Firestore"})
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
