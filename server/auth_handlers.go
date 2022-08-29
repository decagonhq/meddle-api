package server

import (
	"github.com/Timothylock/go-signin-with-apple/apple"
	_"github.com/Timothylock/go-signin-with-apple/apple"
	_"github.com/Timothylock/go-signin-with-apple/apple"
	"context"
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/services/jwt"
	"golang.org/x/oauth2"

	"log"
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
		if err := decode(c, &user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.AuthService.SignupUser(&user)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "Signup successful, check your email for verification", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		if err := decode(c, &loginRequest); err != nil {
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
		conf := config.GetGoogleOAuthConfig(s.Config.GoogleClientID, s.Config.GoogleClientSecret, s.Config.GoogleRedirectURL)
		state, err := jwt.GenerateToken("", s.Config.JWTSecret)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, err)
			return
		}
		url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (s *Server) HandleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		var state = c.Query("state")
		var code = c.Query("code")

		_, err := jwt.ValidateToken(state, s.Config.JWTSecret)

		if err != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid login", http.StatusUnauthorized))
			return
		}

		var oauth2Config = config.GetGoogleOAuthConfig(s.Config.GoogleClientID, s.Config.GoogleClientID, s.Config.GoogleRedirectURL)

		token, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil || token == nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid token", http.StatusUnauthorized))
			return
		}

		authToken, errr := s.AuthService.GoogleSignInUser(token.AccessToken)
		if errr != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid authToken", http.StatusUnauthorized))
			return
		}

		response.JSON(c, "apple sign in successful", http.StatusOK, authToken, nil)
	}
}

func (s *Server) HandleAppleOauthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		appConf := config.GetAppleAuthConfig(s.Config.AppleTeamID, s.Config.AppleClientID, s.Config.AppleRedirectURL, s.Config.AppleP8Cert)

		state, err := jwt.GenerateToken("", s.Config.JWTSecret)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, err)
			return
		}

		url := appConf.CreateCallbackURL(state)
		//c.Redirect(http.StatusTemporaryRedirect, url)
		c.Redirect(http.StatusTemporaryRedirect, "https://appleid.apple.com/auth/authorize?&response_mode=form_post&client_id=com.herokuapp.meddle&redirect_uri=https%3A%2F%2Fmeddle-app.herokuapp.com%2Fapi%2Fv1%2Fapple%2Fauth&response_type=code&scope=name+email&state="+url)
	//https://appleid.apple.com/auth/authorize?&response_mode=form_post&client_id=com.herokuapp.meddle&redirect_uri=https%3A%2F%2Fmeddle-app.herokuapp.com%2Fapi%2Fv1%2Fapple%2Fauth&response_type=code&scope=name+email&state=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IiIsImV4cCI6MTY2MTg3OTEzNX0.A9HwfJ1nNRqUniSMESG6BpWrU3InTzTDBOqXAUj6JnE	//&response_mode=form_post
	}
}

func (s *Server) HandleAppleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret, err := apple.GenerateClientSecret(s.Config.AppleP8Cert, s.Config.AppleTeamID, s.Config.AppleClientID, s.Config.AppleKeyID)
		if err != nil {
			errors.New("invalid login", http.StatusUnauthorized)
		}
		_, err = s.AuthService.AppleSignInUser(secret)
		if err != nil {
			errors.New("invalid login", http.StatusInternalServerError)
		}
		response.JSON(c, "apple sign in successful", http.StatusOK, nil, nil)
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
		claims, errr := jwt.ValidateAndGetClaims(token, s.Config.JWTSecret)
		if errr != nil {
			response.JSON(c, "", http.StatusUnauthorized, nil, errr)
			return
		}
		convertClaims, _ := claims["exp"].(int64) //jwt pkg to validate
		if convertClaims < time.Now().Unix() {
			accBlacklist := &models.BlackList{
				Email: user.Email,
				Token: token,
			}
			if err := s.AuthRepository.AddToBlackList(accBlacklist); err != nil {
				log.Printf("can't add access token to blacklist: %v\n", err)
				response.JSON(c, "logout failed", http.StatusInternalServerError, nil, errors.New("can't add access token to blacklist", http.StatusInternalServerError))
				return
			}
		}
		response.JSON(c, "logout successful", http.StatusOK, nil, nil)

	}
}

func (s *Server) handleFBLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetFacebookOAuthConfig(s.Config.FacebookClientID, s.Config.FacebookClientSecret, s.Config.FacebookRedirectURL)
		state, err := jwt.GenerateToken("", s.Config.JWTSecret)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, err)
			return
		}
		url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (s *Server) fbCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var state = c.Query("state")
		var code = c.Query("code")

		_, err := jwt.ValidateToken(state, s.Config.JWTSecret)
		if err != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid login", http.StatusUnauthorized))
			return
		}

		var OAuth2Config = config.GetFacebookOAuthConfig(s.Config.FacebookClientID, s.Config.FacebookClientSecret, s.Config.FacebookRedirectURL)

		token, err := OAuth2Config.Exchange(context.Background(), code)
		if err != nil || token == nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid token", http.StatusUnauthorized))
			return
		}

		authToken, errr := s.AuthService.FacebookSignInUser(token.AccessToken)
		if errr != nil {
			log.Printf("Facebook Signin failed due to: %v", errr)
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid authToken", http.StatusUnauthorized))
			return
		}

		response.JSON(c, "facebook sign in successful", http.StatusOK, authToken, nil)
	}
}

func (s *Server) handleDeleteUserByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}

		if err := s.AuthService.DeleteUserByEmail(user.Email); err != nil {
			err.Respond(c)
			return
		}

		response.JSON(c, "user successfully deleted", http.StatusOK, nil, nil)
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

func (s *Server) HandleVerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		paramToken := c.Param("token")
		err := s.AuthService.VerifyEmail(paramToken)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	}
}
