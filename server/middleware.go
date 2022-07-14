package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	errs "github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/decagonhq/meddle-api/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Authorize authorizes a request
func (s *Server) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := s.Config.JWTSecret
		accessToken := services.GetTokenFromHeader(c)
		verifiedToken, accessClaims, err := services.TokenChecker(&accessToken, &secret)
		if err != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errs.New("unauthorized", http.StatusUnauthorized))
			return
		}

		//TODO find a way to make sure accessToken wont be nil, because we allow
		//a token is expired error to reach here accessToken will be nil
		//when that happens
		if s.AuthRepository.TokenInBlacklist(&verifiedToken.Raw) || isTokenExpired(accessClaims) {
			respondAndAbort(c, "expired token", http.StatusUnauthorized, nil, errs.New("expired token", http.StatusUnauthorized))
			return
		}

		email, ok := accessClaims["user_email"].(string)
		if !ok {
			respondAndAbort(c, "", http.StatusInternalServerError, nil, errs.New("internal server error", http.StatusInternalServerError))
			return
		}

		var user *models.User
		if user, err = s.AuthRepository.FindUserByEmail(email); err != nil {
			if errors.Is(errs.InActiveUserError, err) {
				respondAndAbort(c, "", http.StatusBadRequest, nil, errs.New(err.Error(), http.StatusUnauthorized))
				return
			}

			respondAndAbort(c, "", http.StatusNotFound, nil, errs.New("user not found", http.StatusUnauthorized))
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

// respondAndAbort calls response.JSON
//and aborts the Context
func respondAndAbort(c *gin.Context, message string, status int, data interface{}, e *errs.Error) {
	response.JSON(c, message, status, data, e)
	c.Abort()
}

func isTokenExpired(claims jwt.MapClaims) bool {
	if exp, ok := claims["exp"].(float64); ok {
		return float64(time.Now().Unix()) > exp
	}
	return true
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
