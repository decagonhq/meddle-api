package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type Config struct {
	Debug                    bool   `envconfig:"debug"`
	Port                     int    `envconfig:"port"`
	Env                      string `envconfig:"env"`
	PostgresHost             string `envconfig:"postgres_host"`
	PostgresPort             int    `envconfig:"postgres_port"`
	PostgresUser             string `envconfig:"postgres_user"`
	PostgresPassword         string `envconfig:"postgres_password"`
	PostgresDB               string `envconfig:"postgres_db"`
	JWTSecret                string `envconfig:"jwt_secret"`
	MEDDLE_MG_PUBLIC_API_KEY string `envconfig:"meddle_mg_public_api_key"`
	MEDDLE_MG_DOMAIN         string `envconfig:"meddle_mg_domain"`
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
