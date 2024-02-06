package functions

import (
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Logout(c *gin.Context) {
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

	// Check if the token has already expired
	if claims.ExpiresAt.Time().Before(time.Now()) {
		c.JSON(401, gin.H{"Error": "Unauthorized. Token has expired."})
		c.Abort()
		return
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now())
	c.JSON(200, gin.H{"Message": "Logged out"})

}
