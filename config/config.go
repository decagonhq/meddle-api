package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"os"
)

const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type Config struct {
	Debug             bool          `envconfig:"debug"`
	Port              int           `envconfig:"port"`
	Env               string        `envconfig:"env"`
	PostgresHost      string        `envconfig:"postgres_host"`
	PostgresPort      int           `envconfig:"postgres_port"`
	PostgresUser      string        `envconfig:"postgres_user"`
	PostgresPassword  string        `envconfig:"postgres_password"`
	PostgresDB        string        `envconfig:"postgres_db"`
	JWTSecret         string        `envconfig:"jwt_secret"`
	GoogleLoginConfig oauth2.Config `envconfig:"google_login_config"`
}

var AppConfig Config

func Load() (*Config, error) {
	env := os.Getenv("GIN_MODE")
	if env != "release" {
		if err := godotenv.Load("../.env"); err != nil {
			log.Printf("couldn't load env vars: %v", err)
		}
	}

	c := &Config{}
	err := envconfig.Process("meddle", c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func LoadConfig() {
	// Oauth configuration for Google
	AppConfig.GoogleLoginConfig = oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/google_callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	//// Oauth configuration for Facebook
	//AppConfig.FacebookLoginConfig = oauth2.Config{
	//	ClientID:     os.Getenv("FB_CLIENT_ID"),
	//	ClientSecret: os.Getenv("FB_CLIENT_SECRET"),
	//	Endpoint:     facebook.Endpoint,
	//	RedirectURL:  "http://localhost:8080/fb_callback",
	//	Scopes: []string{
	//		"email",
	//		"public_profile",
	//	},
	//}
}
