package dto

import (
	"github.com/dgrijalva/jwt-go"
)

// JWTSigningKey represents Signing key to validate token.
// it is important to change below key for the first time
var JWTSigningKey = "S3cr3t"

// JwtCustomClaims represents jwt struct to be passed to client as authentication token
type JwtCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
