package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func CreateJWT(acc *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"username":  acc.Username,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func PermissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Permission denied."})
}

func AuthGuard(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := r.Cookie("token")
		if err != nil {
			PermissionDenied(w)
			return
		}
		token, err := ValidateJWT(t.Value)
		if err != nil {
			PermissionDenied(w)
			return
		}
		if !token.Valid {
			PermissionDenied(w)
			return
		}
		f(w, r)
	}
}
