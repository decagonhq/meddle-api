package services

import (
	"fmt"
	"github.com/decagonhq/meddle-api/config"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var testConfig *config.Config

func TestMain(m *testing.M) {
	fmt.Println("Starting server tests")
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("couldn't load env vars: %v", err)
	}
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	testConfig = c
	fmt.Println(testConfig)
	exitCode := m.Run()
	os.Exit(exitCode)
}
