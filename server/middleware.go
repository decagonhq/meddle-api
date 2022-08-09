package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"

	errs "github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

// Authorize authorizes a request
func (s *Server) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := s.Config.JWTSecret
		if isJWTSecretEmpty(secret) {
			respondAndAbort(c, "", http.StatusInternalServerError, nil, errs.New("internal server error", http.StatusInternalServerError))
			return
		}

		accessToken := getTokenFromHeader(c)
		if isAccessTokenEmpty(accessToken) {
			respondAndAbort(c, "invalid token", http.StatusUnauthorized, nil, errs.New("unauthorized", http.StatusUnauthorized))
			return
		}

		token, err := validateToken(accessToken, secret)
		if err != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errs.New("unauthorized", http.StatusUnauthorized))
			return
		}

		accessClaims, err := getClaims(token)
		if err != nil {
			respondAndAbort(c, "", http.StatusInternalServerError, nil, errs.New("internal server error", http.StatusInternalServerError))
			return
		}

		if s.AuthRepository.TokenInBlacklist(accessToken) || isTokenExpired(accessClaims) {
			respondAndAbort(c, "expired token", http.StatusUnauthorized, nil, errs.New("expired token", http.StatusUnauthorized))
			return
		}

		email, ok := accessClaims["email"].(string)
		if !ok {
			respondAndAbort(c, "", http.StatusInternalServerError, nil, errs.New("internal server error", http.StatusInternalServerError))
			return
		}

		var user *models.User
		if user, err = s.AuthRepository.FindUserByEmail(email); err != nil {
			switch {
			case errors.Is(err, errs.InActiveUserError):
				respondAndAbort(c, "inactive user", http.StatusUnauthorized, nil, errs.New(err.Error(), http.StatusUnauthorized))
				return
			case errors.Is(err, gorm.ErrRecordNotFound):
				respondAndAbort(c, "user not found", http.StatusUnauthorized, nil, errs.New(err.Error(), http.StatusUnauthorized))
				return
			default:
				respondAndAbort(c, "", http.StatusInternalServerError, nil, errs.New("internal server error", http.StatusInternalServerError))
				return

			}
		}

		c.Set("access_token", accessToken)
		c.Set("user", user)

		c.Next()
	}
}

// respondAndAbort calls response.JSON and aborts the Context
func respondAndAbort(c *gin.Context, message string, status int, data interface{}, e *errs.Error) {
	response.JSON(c, message, status, data, e)
	c.Abort()
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

// getTokenFromHeader returns the token string in the authorization header
func getTokenFromHeader(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	if len(authHeader) > 8 {
		return authHeader[7:]
	}
	return ""
}
