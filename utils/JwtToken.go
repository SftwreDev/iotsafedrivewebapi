package utils

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"iotsafedriveapi/structs"
	"log"
	"os"
	"time"
)

func NewAccessToken(claims structs.UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) (*structs.UserClaims, error) {
	// Parse the access token with the custom claims struct
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &structs.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Retrieve the token secret from environment variables
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	// Check if there's an error parsing the access token
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !parsedAccessToken.Valid {
		return nil, errors.New("access token is invalid")
	}

	// Extract and return the custom claims from the parsed access token
	claims, ok := parsedAccessToken.Claims.(*structs.UserClaims)
	if !ok {
		return nil, errors.New("failed to extract custom claims from access token")
	}

	return claims, nil
}

func ParseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedRefreshToken.Claims.(*jwt.StandardClaims)
}

func GenerateAccessAndRefreshTokens(userClaims structs.UserClaims) (string, string, error) {
	signedAccessToken, err := NewAccessToken(userClaims)
	if err != nil {
		log.Fatal("error creating access token")
		return "", "", err
	}

	refreshClaims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	signedRefreshToken, err := NewRefreshToken(refreshClaims)
	if err != nil {
		log.Fatal("error creating refresh token")
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}
