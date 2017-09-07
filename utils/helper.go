package utils

import (
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// CreateToken creates token on successful login or signup
func CreateToken(userID int64, username string) string {
	claims := MyCustomClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: 15000,
		},
	}

	godotenv.Load()

	key := os.Getenv("SECRET")
	mySigningKey := []byte(key)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

// MyCustomClaims customize jwt
type MyCustomClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
