package functions

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"cloud.google.com/go/logging"
	"github.com/Capstone/models"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/google/uuid"
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
	router.POST("/", CreateItem)

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("CreateItem", router.Handler().ServeHTTP)
}

var newProduct models.Product

// @Summary Create a new grocery item
// @Description Create a new grocery item with product information and an optional image
// @Tags Grocery
// @Security bearerToken
// @Accept json
// @Produce json
// @Param Image formData file false "Image file (.jpg)"
// @Param ProductBody formData models.ProductRequestBody true "Product Information"
// @Success 201 {object} models.CreationSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /createItem [post]
func CreateItem(c *gin.Context) {
	var Caser = cases.Title(language.English)

	randomString := uuid.New().String()[:8]

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

	var (
		g          = galidator.New()
		customizer = g.Validator(models.Product{})
	)

	// Bind the form data to the Product struct and perform validation
	if err := c.ShouldBind(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": customizer.DecryptErrors(err)})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Failed to unmarshal JSON: %s", err.Error()),
			Severity: logging.Error,
		})
		return

	}

	newProduct.ProductName = Caser.String(newProduct.ProductName)
	newProduct.Category = Caser.String(newProduct.Category)

	// Retrieve the uploaded file from the request
	uploadedFile, err := c.FormFile("Image")

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
	} else {
		// Check if the uploaded file is a .jpg file
		if !strings.HasSuffix(uploadedFile.Filename, ".jpg") {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid file format. Image file should be '.jpg' file"})
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Invalid file format : %v", err),
				Severity: logging.Error,
			})
			return
		}

		// Open the file
		blobFile, err := uploadedFile.Open()
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
			img_url := GenerateThumbnail(c, randomString)

			newProduct.ThumbnailURL = img_url
		}()
		// wg.Wait() to wait for the thumbnail generation to complete.
		wg.Wait()

		// Upload the file to Google Cloud Storage
		img_url, err := UploadFile(blobFile, randomString+uploadedFile.Filename, "shreyas-grocery-image-bucket")
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

		newProduct.ImageURL = img_url
	}
	newProduct.ProductID = randomString

	//Store other data to firestore.
	_, err = firestoreClient.Collection("grocery").Doc(randomString).Set(ctx, newProduct)
	if err != nil {
		fmt.Println("Error creating product:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal Server Error"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while saving item to firestore: %s", err.Error()),
			Severity: logging.Error,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Document created with id : ": randomString,
		"Product Information":         newProduct})
	logger.Log(logging.Entry{
		Payload:  fmt.Sprintf("Grocery item created with id %v", randomString),
		Severity: logging.Info,
	})
}
