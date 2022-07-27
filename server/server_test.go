package server

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/decagonhq/meddle-api/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	server *Server
)

var testServer struct {
	router  *gin.Engine
	handler *Server
}

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("couldn't load env vars: %v", err)
	}
	fmt.Println("Starting server tests")
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	testServer.handler = &Server{
		Config: c,
	}
	testServer.handler.Config.JWTSecret = "testSecret"
	testServer.router = testServer.handler.setupRouter()
	exitCode := m.Run()
	os.Exit(exitCode)
}
