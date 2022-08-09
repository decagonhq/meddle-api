package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/errors"
	"github.com/golang-jwt/jwt"
)

// verifyAccessToken verifies a token
func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func isTokenExpired(claims jwt.MapClaims) bool {
	if exp, ok := claims["exp"].(float64); ok {
		return float64(time.Now().Unix()) > exp
	}
	return true
}

func isJWTSecretEmpty(secret string) bool {
	return secret == ""
}

func isAccessTokenEmpty(token string) bool {
	return token == ""
}

func validateToken(token string, secret string) (*jwt.Token, error) {
	tk, err := verifyToken(token, secret)
	if err != nil {
		log.Println(err)                                 // TODO: remove
		return nil, fmt.Errorf("invalid token: %v", err) // TODO: probably need to errors.NEw
	}
	if !tk.Valid {
		return nil, errors.New("invalid token", http.StatusUnauthorized)
	}
	return tk, nil
}

func getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not get claims")
	}
	return claims, claims.Valid()
}
