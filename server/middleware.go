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
		accToken := services.GetTokenFromHeader(c)
		accessToken, accessClaims, err := services.AuthorizeToken(&accToken, &secret)
		if err != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errs.New("unauthorized", http.StatusUnauthorized))
			return
		}

		//TODO find a way to make sure accesstoken wont be nil, because we allow
		//a token is epired error to reach here accessToken will be nill
		//when that happens
		if s.AuthRepository.TokenInBlacklist(&accessToken.Raw) || isTokenExpired(accessClaims) {
			rt := &struct {
				RefreshToken string `json:"refresh_token,omitempty" binding:"required"`
			}{}

			if err := c.ShouldBindJSON(rt); err != nil {
				respondAndAbort(c, "", http.StatusBadRequest, nil, errs.New("unauthorized", http.StatusBadRequest))
				return
			}

			if s.AuthRepository.TokenInBlacklist(&rt.RefreshToken) {
				respondAndAbort(c, "", http.StatusUnauthorized, nil, errs.New("refresh token is invalid", http.StatusUnauthorized))
				return
			}

			_, rtClaims, err := services.AuthorizeToken(&rt.RefreshToken, &secret)
			if err != nil {
				respondAndAbort(c, "", http.StatusUnauthorized, nil, errs.New("refresh token is invalid", http.StatusUnauthorized))
				return
			}

			if isTokenExpired(rtClaims) {
				respondAndAbort(c, "", http.StatusUnauthorized, nil, errs.New("refresh token is invalid", http.StatusUnauthorized))
				return
			}

			if sub, ok := rtClaims["sub"].(float64); ok && sub != 1 {
				respondAndAbort(c, "", http.StatusUnauthorized, nil, errs.New("refresh token is invalid", http.StatusUnauthorized))
				return
			}

			//generate a new access token, and rest its exp time
			accessClaims["exp"] = time.Now().Add(services.AccessTokenValidity).Unix()
			newAccessToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
			if err != nil {
				respondAndAbort(c, "", http.StatusUnauthorized, nil, errs.New("can't generate new access token", http.StatusUnauthorized))
				return
			}
			respondAndAbort(c, "new access token generated", http.StatusOK, gin.H{"access_token": *newAccessToken}, nil)
			return
		}

		email, ok := accessClaims["user_email"].(string)
		if !ok {
			respondAndAbort(c, "", http.StatusInternalServerError, nil, errs.New("internal server error", http.StatusInternalServerError))
			return
		}

		var user *models.User
		if user, err = s.AuthRepository.FindUserByUsername(email); err != nil {
			if errors.Is(errs.InActiveUserError, err) {
				respondAndAbort(c, "", http.StatusBadRequest, nil, errs.New(err.Error(), http.StatusUnauthorized))
				return
			}

			respondAndAbort(c, "", http.StatusNotFound, nil, errs.New("user not found", http.StatusUnauthorized))
			return
		}

		c.Set("user", user)
		c.Set("access_token", accessToken.Raw)

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
