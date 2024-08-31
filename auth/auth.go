package jwtutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	fmt.Println(username)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(username)
	fmt.Println("1")
	return token.SignedString(jwtKey)
}

func ValidateToken(tknStr string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func CheckToken(w http.ResponseWriter, r *http.Request) (string, int) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No token", http.StatusUnauthorized)
			return " ", 0
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return " ", 0
	}

	tokenStr := c.Value
	fmt.Println(tokenStr)
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return " ", 0
	}

	return claims.Username, 0
}
