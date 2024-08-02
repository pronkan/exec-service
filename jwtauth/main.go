package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {

	// Arg -key to get secret key on execution
	tokenSecret := flag.String("key", "", "Secret key for JWT token")
	flag.Parse()

	// Create claims (you can add more claims like user ID, roles, etc.)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	signedToken, err := token.SignedString([]byte(*tokenSecret))
	if err != nil {
		// Handle error
	}

	fmt.Println("Secret key: ", *tokenSecret)
	fmt.Println("JWT token: ", signedToken)
}
