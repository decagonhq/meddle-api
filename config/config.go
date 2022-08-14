package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"os"
)

type Config struct {
	Debug              bool   `envconfig:"debug"`
	Port               int    `envconfig:"port"`
	Env                string `envconfig:"env"`
	PostgresHost       string `envconfig:"postgres_host"`
	PostgresPort       int    `envconfig:"postgres_port"`
	PostgresUser       string `envconfig:"postgres_user"`
	PostgresPassword   string `envconfig:"postgres_password"`
	PostgresDB         string `envconfig:"postgres_db"`
	JWTSecret          string `envconfig:"jwt_secret"`
	MailgunApiKey      string `envconfig:"mg_public_api_key"`
	MgDomain           string `envconfig:"mg_domain"`
	EmailFrom          string `envconfig:"email_from"`
	Host               string `envconfig:"host"`
	GoogleClientID     string `envconfig:"google_client_id"`
	GoogleClientSecret string `envconfig:"google_client_secret"`
	GoogleRedirectURL  string `envconfig:"google_redirect_url"`
	OauthStateString   string `envconfig:"oauth_state_string"`
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

func GetGoogleOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"email"},
	}
}
