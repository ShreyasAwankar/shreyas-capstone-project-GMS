package functions

import (
	"fmt"
	"os"

	"github.com/Capstone/models"
	"github.com/golang-jwt/jwt/v5"
)

// VerifyToken is the function to verify the JWT token
func VerifyToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("jwtKey")), nil // Replace with your actual jwtKey retrieval logic
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
