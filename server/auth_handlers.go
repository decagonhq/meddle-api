package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/decagonhq/meddle-api/config"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Server) HandleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.AuthService.SignupUser(&user)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "user created successfully", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			response.JSON(c, "", errors.ErrBadRequest.Status, nil, err)
			return
		}
		userResponse, err := s.AuthService.LoginUser(&loginRequest)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}
		response.JSON(c, "login successful", http.StatusOK, userResponse, nil)
	}
}

func (s *Server) HandleGoogleOauthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create oauthState cookie
		oauthState := GenerateStateOauthCookie(c.Writer)
		/*
			AuthCodeURL receive state that is a token to protect the user
			from CSRF attacks. You must always provide a non-empty string
			and validate that it matches the the state query parameter
			on your redirect callback.
		*/
		u := config.AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
		http.Redirect(c.Writer, c.Request, u, http.StatusTemporaryRedirect)
	}
}

func (s *Server) HandleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get oauth state from cookie for this user
		oauthState, _ := c.Request.Cookie("oauthstate")
		state := c.Request.FormValue("state")
		code := c.Request.FormValue("code")
		c.Writer.Header().Add("content-type", "application/json")

		// ERROR : Invalid OAuth State
		if state != oauthState.Value {
			http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
			fmt.Fprintf(c.Writer, "invalid oauth google state")
			return
		}

		// Exchange Auth Code for Tokens
		token, err := config.AppConfig.GoogleLoginConfig.Exchange(
			context.Background(), code)

		// ERROR : Auth Code Exchange Failed
		if err != nil {
			fmt.Fprintf(c.Writer, "falied code exchange: %s", err.Error())
			return
		}

		// Fetch User Data from google server
		resp, err := http.Get(config.OauthGoogleUrlAPI + token.AccessToken)

		// ERROR : Unable to get user data from google
		if err != nil {
			fmt.Fprintf(c.Writer, "failed getting user info: %s", err.Error())
			return
		}

		// Parse user data JSON Object
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(c.Writer, "failed read response: %s", err.Error())
			return
		}

		// send back response to browser
		fmt.Fprintln(c.Writer, string(contents))

		response.JSON(c, "login successful", http.StatusOK, resp, nil)
	}
}

func GetValuesFromContext(c *gin.Context) (string, *models.User, *errors.Error) {
	var tokenI, userI interface{}
	var tokenExists, userExists bool

	if tokenI, tokenExists = c.Get("access_token"); !tokenExists {
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}
	if userI, userExists = c.Get("user"); !userExists {
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}

	token, ok := tokenI.(string)
	if !ok {
		return "", nil, errors.New("internal server error", http.StatusInternalServerError)
	}
	user, ok := userI.(*models.User)
	if !ok {
		return "", nil, errors.New("internal server error", http.StatusInternalServerError)
	}

	return token, user, nil
}

func (s *Server) handleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, user, err := GetValuesFromContext(c)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}

		claims, errr := getClaims(token, s.Config.JWTSecret)
		if errr != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errr)
			return
		}
		convertClaims, _ := claims["exp"].(int64) //jwt pkg to validate
		if convertClaims < time.Now().Unix() {
			response.JSON(c, "successfully logged out", http.StatusOK, nil, nil)
			return
		} else {
			accBlacklist := &models.BlackList{
				Email: user.Email,
				Token: token,
			}
			if err := s.AuthRepository.AddToBlackList(accBlacklist); err != nil {
				log.Printf("can't add access token to blacklist: %v\n", err)
				response.JSON(c, "logout failed", http.StatusInternalServerError, nil, errors.New("can't add access token to blacklist", http.StatusInternalServerError))
				return
			}
			response.JSON(c, "successfully added to blacklist", http.StatusOK, nil, nil)
		}
	}
}

func (s *Server) handleGetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleUpdateUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleShowProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(2 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	return state
}
