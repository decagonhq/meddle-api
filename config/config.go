package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"log"
	"os"
)

type Config struct {
	Debug                bool   `envconfig:"debug"`
	Port                 int    `envconfig:"port"`
	Env                  string `envconfig:"env"`
	PostgresHost         string `envconfig:"postgres_host"`
	PostgresPort         int    `envconfig:"postgres_port"`
	PostgresUser         string `envconfig:"postgres_user"`
	PostgresPassword     string `envconfig:"postgres_password"`
	PostgresDB           string `envconfig:"postgres_db"`
	JWTSecret            string `envconfig:"jwt_secret"`
	FacebookClientID     string `envconfig:"facebook_client_id"`
	FacebookClientSecret string `envconfig:"facebook_client_secret"`
	FacebookRedirectURL  string `envconfig:"facebook_redirect_url"`
	MailgunApiKey        string `envconfig:"mg_public_api_key"`
	MgDomain             string `envconfig:"mg_domain"`
	EmailFrom            string `envconfig:"email_from"`
	Host                 string `envconfig:"host"`
	GoogleClientID       string `envconfig:"google_client_id"`
	GoogleClientSecret   string `envconfig:"google_client_secret"`
	GoogleRedirectURL    string `envconfig:"google_redirect_url"`
	AppleClientID       string `envconfig:"apple_client_id"`
	AppleClientSecret   string `envconfig:"apple_client_secret"`
	AppleRedirectURL    string `envconfig:"apple_redirect_url"`
}

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

func GetFacebookOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     facebook.Endpoint,
		Scopes:       []string{"email"},
	}
}

func GetGoogleOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"email"},
	}
}

func GetAppleOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"email"},
	}
}