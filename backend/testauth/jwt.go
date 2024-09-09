package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT generates a JWT for testing purposes
func GenerateJWT() (string, error) {
	// Define token expiration time (1 hour from now)
	expirationTime := time.Now().Add(time.Hour * 1).Unix()

	// Create the JWT claims, which includes the user ID and the expiration time
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
		Issuer:    "uct",
		Subject:   "User 1", // We can pretend 'User 1' is the user ID for the sake of this example
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func main() {
	// Generate a JWT
	token, err := GenerateJWT()
	if err != nil {
		fmt.Println("Error generating JWT:", err)
		return
	}

	// Print the generated JWT
	fmt.Println(token)
}
