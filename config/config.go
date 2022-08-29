package config

import (
	"github.com/BillSJC/appleLogin"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Debug                bool   `envconfig:"debug"`
	Port                 int    `envconfig:"port"`
	PostgresHost         string `envconfig:"postgres_host"`
	PostgresUser         string `envconfig:"postgres_user"`
	PostgresDB           string `envconfig:"postgres_db"`
	MailgunApiKey        string `envconfig:"mg_public_api_key"`
	EmailFrom            string `envconfig:"email_from"`
	BaseUrl              string `envconfig:"base_url"`
	Env                  string `envconfig:"env"`
	PostgresPort         int    `envconfig:"postgres_port"`
	PostgresPassword     string `envconfig:"postgres_password"`
	JWTSecret            string `envconfig:"jwt_secret"`
	FacebookClientID     string `envconfig:"facebook_client_id"`
	FacebookClientSecret string `envconfig:"facebook_client_secret"`
	FacebookRedirectURL  string `envconfig:"facebook_redirect_url"`
	MgDomain             string `envconfig:"mg_domain"`
	Host                 string `envconfig:"host"`
	GoogleClientID       string `envconfig:"google_client_id"`
	GoogleClientSecret   string `envconfig:"google_client_secret"`
	GoogleRedirectURL    string `envconfig:"google_redirect_url"`
	AppleClientID       string `envconfig:"apple_client_id"`
	AppleP8Cert   string `envconfig:"apple_p8_cert"`
	AppleRedirectURL    string `envconfig:"apple_redirect_url"`
	AppleKeyID   string `envconfig:"apple_key_id"`
	AppleTeamID   string `envconfig:"team_id"`
	AndroidClientID	   string `envconfig:"android_client_id"`
}

func Load() (*Config, error) {
	env := os.Getenv("GIN_MODE")
	if env != "release" {
		if err := godotenv.Load("./.env"); err != nil {
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

func GetAppleAuthConfig(appleTeamID, appleClientID, appleRedirectURL, appleKeyID string) *appleLogin.AppleConfig {
	return &appleLogin.AppleConfig{
		TeamID: appleTeamID,
		ClientID: appleClientID,
		RedirectURI: appleRedirectURL,
		KeyID: appleKeyID,
	}
}
