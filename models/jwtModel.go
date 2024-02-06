package models

import "github.com/golang-jwt/jwt/v5"

// Create a struct to read the username and password from the request body
type Credentials struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type SuccessResponseJWT struct {
	Success string `json:"Success"`
	Token   string `json:"Token"`
}
