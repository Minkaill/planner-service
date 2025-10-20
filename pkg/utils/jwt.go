package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

func GenerateJWT() (string, error) {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return "", nil
	}

	hash := sha256.Sum256([]byte(password))
	hashStr := hex.EncodeToString(hash[:])

	expiration := time.Now().Add(8 * time.Hour)

	claims := &Claims{
		Hash: hashStr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) bool {
	if tokenString == "" {
		return false
	}

	password := os.Getenv("TODO_PASSWORD")
	hash := sha256.Sum256([]byte(password))
	expectedHash := hex.EncodeToString(hash[:])

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.Hash != expectedHash {
		return false
	}

	return true
}
