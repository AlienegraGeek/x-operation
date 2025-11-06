package test

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestGetClientID(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	t.Logf("env: %s", os.Getenv("CLIENT_ID"))
}
