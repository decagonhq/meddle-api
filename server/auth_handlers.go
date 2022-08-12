package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/decagonhq/meddle-api/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/jwt"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

var state string

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
		tempState, _ := services.GenerateRandomString()
		state = tempState

		url := oAuthGoogleConfig().AuthCodeURL(state)
		response.JSON(c, "", http.StatusTemporaryRedirect, url, nil)
	}
}
func (s *Server) HandleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Param("state") != state {
			response.JSON(c, "", http.StatusTemporaryRedirect, nil, errors.New("invalid state", http.StatusTemporaryRedirect))
			return
		}
		token, err := oAuthGoogleConfig().Exchange(context.Background(), c.Param("code"))
		if err != nil {
			fmt.Print(err)
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("invalid code", http.StatusInternalServerError))
			return
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			fmt.Print(err)
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("error", http.StatusInternalServerError))
			return
		}
		defer resp.Body.Close()

		googleResponse := models.GoogleAuthResponse{}
		err = json.NewDecoder(resp.Body).Decode(&googleResponse)
		if err != nil {
			fmt.Println(err)
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("error", http.StatusInternalServerError))
			return
		}
		var userData = models.User{
			Name:        googleResponse.GivenName + " " + googleResponse.FamilyName,
			Email:       googleResponse.Email,
			AccessToken: token.AccessToken,
			Social:      "google",
		}
		//Create new user if not exists else update access token
		//err = s.AuthRepository.IsEmailExist(userData.Email)
		//if err != nil {
		//	if err.Error() == "email not found" {
		//		userData.Password, _ = services.GenerateRandomString()
		//		userResponse, err := s.AuthService.SignupUser(&userData)
		//		if err != nil {
		//			err.Respond(c)
		//			return
		//		}
		//		response.JSON(c, "user created successfully", http.StatusCreated, userResponse, nil)
		//	} else {
		//		response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("error", http.StatusInternalServerError))
		//	}
		//	return
		//}
		generateToken, err := services.GenerateToken(userData.Email, s.Config.JWTSecret)
		if err != nil {
			fmt.Println(err)
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("error", http.StatusInternalServerError))
			return
		}
		userData.AccessToken = generateToken
		response.JSON(c, "", http.StatusOK, userData, nil)
	}
}

func oAuthGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

//func (s *Server) HandleGoogleCallback() gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get oauth state from cookie for this user
//		oauthState, _ := c.Request.Cookie("oauthstate")
//		state := c.Request.FormValue("state")
//		code := c.Request.FormValue("code")
//		c.Writer.Header().Add("content-type", "application/json")
//
//		// ERROR : Invalid OAuth State
//		if state != oauthState.Value {
//			http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
//			fmt.Fprintf(c.Writer, "invalid oauth google state")
//			return
//		}
//
//		// Exchange Auth Code for Tokens
//		token, err := config.AppConfig.GoogleLoginConfig.Exchange(
//			context.Background(), code)
//
//		// ERROR : Auth Code Exchange Failed
//		if err != nil {
//			fmt.Fprintf(c.Writer, "falied code exchange: %s", err.Error())
//			return
//		}
//
//		// Fetch User Data from google server
//		resp, err := http.Get(config.OauthGoogleUrlAPI + token.AccessToken)
//
//		// ERROR : Unable to get user data from google
//		if err != nil {
//			fmt.Fprintf(c.Writer, "failed getting user info: %s", err.Error())
//			return
//		}
//
//		// Parse user data JSON Object
//		defer resp.Body.Close()
//		contents, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			fmt.Fprintf(c.Writer, "failed read response: %s", err.Error())
//			return
//		}
//
//		// send back response to browser
//		fmt.Fprintln(c.Writer, string(contents))
//
//		response.JSON(c, "login successful", http.StatusOK, resp, nil)
//	}
//}

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
