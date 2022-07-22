package server

import (
	"fmt"
	"github.com/decagonhq/meddle-api/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

//var (
//	server     *Server
//	testRouter *gin.Engine
//)

var testServer struct {
	router  *gin.Engine
	handler *Server
}

func TestMain(m *testing.M) {
	fmt.Println("Starting server tests")
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("couldn't load env vars: %v", err)
	}
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
