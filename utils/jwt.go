package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// struct for payload token
type Claims struct {
	UserID uint `json:"user_id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// function to generate JWT Token
func GenerateToken(userID uint, email string, role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key_if_env_empty"
	}

	// make expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// input data user to Claims
	claims := &Claims{
		UserID: userID,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// making token with HS256 algorithm signing
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token with secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
