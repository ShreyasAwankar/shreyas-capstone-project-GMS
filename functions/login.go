package functions

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/logging"
	"github.com/Capstone/models"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golodash/galidator"
)

func init() {
	gin.SetMode(gin.ReleaseMode)

	// Initialize a Gin router
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
	router.POST("/", Login)
	// router.Use(corsMiddleware())

	// Register the Gin router as the handler for the Cloud Function
	functions.HTTP("Login", router.Handler().ServeHTTP)
}

// @Summary User login
// @Description Login to access APIs
// @Tags Authenticatuion
// @Accept       json
// @Produce      json
// @Param RequestBody body models.Credentials true "Credentials"
// @Success 200 {object} models.SuccessResponseJWT
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func Login(c *gin.Context) {

	var creds models.Credentials

	var (
		g          = galidator.New()
		customizer = g.Validator(models.Credentials{})
	)

	if err := c.BindJSON(&creds); err != nil {

		c.JSON(http.StatusNotAcceptable, gin.H{"Message": customizer.DecryptErrors(err)})

		logger.Log(logging.Entry{
			Payload:  "Invalid JSON input for user credential during login.",
			Severity: logging.Error,
		})
		return
	}

	// Get the expected password from our collection

	userSnap, err := firestoreClient.Collection("users").Doc(creds.Username).Get(ctx)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Invalid credentials"})

		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error retrieving user data: %v", err),
			Severity: logging.Error,
		})
		return
	}

	user := userSnap.Data()

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status

	password, ok := user["password"].(string)

	if !ok || password != creds.Password {
		c.JSON(401, gin.H{"Error": "Incorrect password."})
		logger.Log(logging.Entry{
			Payload:  "Incorrect password provided at login.",
			Severity: logging.Error,
		})
		return
	}

	expirationTime := time.Now().Add(20 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("jwtKey")))
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.JSON(500, gin.H{"Error": "Internal Server Error"})
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error while creating jwt token. %v", err),
			Severity: logging.Error,
		})
		return
	}

	// Set the token in the response header
	c.Header("Authorization", tokenString)
	c.JSON(200, gin.H{"Success": fmt.Sprintf("Login successful. Wellcome %v.", claims.Username),
		"Token": tokenString,
	})

	logger.Log(logging.Entry{
		Payload:  fmt.Sprintf("User logged in with username %v", claims.Username),
		Severity: logging.Info,
	})

}
